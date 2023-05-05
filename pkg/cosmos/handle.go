package cosmos

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ycryptx/chaincannon/pkg/benchmark"
	"github.com/ycryptx/chaincannon/pkg/ierror"
	"google.golang.org/grpc"
)

func Handle(recipe *benchmark.Recipe) error {
	ctx := context.Background()

	// make grpc connection to blockchain node
	grpcConn, err := grpc.Dial(
		recipe.Endpoint,
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("%s: %s", ierror.ERROR_CONNECTION, err)
	}
	defer grpcConn.Close()
	ctx = context.WithValue(ctx, "rpc", grpcConn)

	// open all tx file scanners
	transactionFileReadersPerRun := [][]*bufio.Scanner{}
	for i, run := range recipe.Runs {
		for _, path := range run.TxPaths {
			scanner, err := GetSignedTxFileScanner(path)
			if err != nil {
				return err
			}
			transactionFileReadersPerRun[i] = append(transactionFileReadersPerRun[i], scanner)
		}
	}

	for _, run := range transactionFileReadersPerRun {
		go Run(ctx, run)
	}

	grpcResp, err := SendTx(ctx, []byte{})
	if err != nil {
		return err
	}

	fmt.Println(grpcResp.Code)

	return nil
}

func Run(ctx context.Context, txFiles []*bufio.Scanner) {
	// TODO: implement
}

func GetSignedTxFileScanner(path string) (*bufio.Scanner, error) {
	readFile, err := os.Open(path)
	defer readFile.Close()

	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	return fileScanner, nil
}

func SendTx(ctx context.Context, txData []byte) (*types.TxResponse, error) {
	grpcConn, ok := ctx.Value("rpc").(*grpc.ClientConn)
	if !ok {
		return nil, errors.New(ierror.ERROR_CONNECTION_NOT_IN_CONTEXT)
	}

	txClient := tx.NewServiceClient(grpcConn)
	grpcRes, err := txClient.BroadcastTx(
		ctx,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txData,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", ierror.BROADCAST_TX_ERROR, err)
	}

	return grpcRes.TxResponse, nil
}
