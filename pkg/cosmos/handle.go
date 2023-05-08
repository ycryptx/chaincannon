package cosmos

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/schollz/progressbar/v3"
	"github.com/ycryptx/chaincannon/pkg/benchmark"
	"github.com/ycryptx/chaincannon/pkg/ierror"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Handle(ctx context.Context) {
	log, _ := ctx.Value("log").(*zap.Logger)
	recipe, _ := ctx.Value("recipe").(*benchmark.Recipe)
	bar, _ := ctx.Value("bar").(*progressbar.ProgressBar)

	ctx, cancel := context.WithTimeout(ctx, recipe.Duration)

	monitor := benchmark.InitMonitoring(ctx, cancel)
	ctx = context.WithValue(ctx, "monitoring", monitor)

	err := blockMonitor(ctx, cancel)
	if err != nil {
		log.Error(err.Error())
		cancel()
		return
	}

	// make grpc connection to blockchain node
	grpcConn, err := grpc.Dial(
		recipe.Endpoint,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Error(fmt.Sprintf("%s: %s", ierror.ERROR_CONNECTION, err.Error()))
		cancel()
		return
	}
	ctx = context.WithValue(ctx, "rpc", grpcConn)

	// open all tx file scanners
	transactionFileReadersPerRun := [][]*bufio.Scanner{}
	for i, run := range recipe.Runs {
		transactionFileReadersPerRun = append(transactionFileReadersPerRun, []*bufio.Scanner{})
		for _, path := range run.TxPaths {
			scanner, err := GetSignedTxFileScanner(path)
			if err != nil {
				log.Error(err.Error())
				cancel()
			}
			transactionFileReadersPerRun[i] = append(transactionFileReadersPerRun[i], scanner)
		}
	}

	for i, run := range transactionFileReadersPerRun {
		isLast := i == len(transactionFileReadersPerRun)-1
		go Run(ctx, run, cancel, isLast)
	}

	<-ctx.Done()
	grpcConn.Close()
	bar.Finish()
	cancel()
	monitor.Report.PrintReport(ctx)
}

func Run(ctx context.Context, txFiles []*bufio.Scanner, cancel context.CancelFunc, isLast bool) {
	log, _ := ctx.Value("log").(*zap.Logger)
	grpcConn, ok := ctx.Value("rpc").(*grpc.ClientConn)
	if !ok {
		log.Error(ierror.ERROR_CONNECTION_NOT_IN_CONTEXT)
		cancel()
		return
	}
	monitoring, _ := ctx.Value("monitoring").(*benchmark.Monitoring)

	txClient := tx.NewServiceClient(grpcConn)
	for _, fileReader := range txFiles {
		for fileReader.Scan() {
			if ctx.Err() != nil { // circuit-break if cancel has already been called
				return
			}
			txData := fileReader.Text()
			txData = strings.Trim(txData, "\n")
			decoded, err := base64.StdEncoding.DecodeString(txData)
			if err != nil {
				log.Error(err.Error())
			}
			start := time.Now()
			hash, err := SendTx(ctx, txClient, decoded)
			if err != nil {
				log.Error(err.Error())
				cancel()
				return
			}
			go func() {
				monitoring.Stream <- &benchmark.TxPending{
					Start: &start,
					Hash:  hash,
				}
			}()
		}
	}
	if isLast {
		monitoring.Done = true
	}
}

func GetSignedTxFileScanner(path string) (*bufio.Scanner, error) {
	readFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	return fileScanner, nil
}

func SendTx(ctx context.Context, txClient tx.ServiceClient, txData []byte) (string, error) {
	grpcRes, err := txClient.BroadcastTx(
		ctx,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txData,
		},
	)
	if err != nil {
		return "", fmt.Errorf("%s: %s", ierror.BROADCAST_TX_ERROR, err)
	}
	if grpcRes.TxResponse.Code != 0 {
		return "", fmt.Errorf("%s: %s", ierror.BROADCAST_TX_ERROR, grpcRes.TxResponse.RawLog)
	}
	return grpcRes.TxResponse.TxHash, nil
}

func blockMonitor(ctx context.Context, cancel context.CancelFunc) error {
	log, _ := ctx.Value("log").(*zap.Logger)
	monitoring, _ := ctx.Value("monitoring").(*benchmark.Monitoring)
	rpcClient, err := client.NewClientFromNode(fmt.Sprintf("tcp://%s", "0.0.0.0:26657"))
	if err != nil {
		return fmt.Errorf("%s: %s", ierror.ERROR_CONNECTION, err.Error())
	}

	err = rpcClient.Start()
	if err != nil {
		return fmt.Errorf("%s: %s", ierror.ERROR_CONNECTION, err.Error())
	}
	// Subscribe to new blocks
	query := "tm.event = 'NewBlock'"
	eventCh, err := rpcClient.Subscribe(context.Background(), "block-subscriber", query)
	if err != nil {
		return err
	}

	var lastBlockTime time.Time

	// Listen for new blocks
	go func() {
		for {
			event := <-eventCh
			blockEvent, ok := event.Data.(types.EventDataNewBlock)
			if !ok {
				log.Warn("Cosmos block subscriber: nexpected event data")
				continue
			}
			endTime := time.Now()
			block := blockEvent.Block
			hashes := []string{}

			blockTime := block.Time
			if !lastBlockTime.IsZero() {
				monitoring.Report.RecordBlockTime(ctx, lastBlockTime, blockTime)
			}
			lastBlockTime = blockTime

			resBlock, err := rpcClient.TxSearch(ctx, fmt.Sprintf("tx.height=%d", block.Height), false, nil, nil, "asc")
			if err != nil {
				log.Fatal(err.Error())
			}

			for _, tx := range resBlock.Txs {
				hash := tx.Hash.String()
				if monitoring.IsOurTx(ctx, hash) {
					hashes = append(hashes, hash)
				}
			}
			if len(hashes) > 0 {
				monitoring.RecordTxs(ctx, hashes, endTime)
			}

			if monitoring.Done && monitoring.NoMorePendingTxs() && len(hashes) == 0 {
				cancel()
			}
		}
	}()
	return nil
}
