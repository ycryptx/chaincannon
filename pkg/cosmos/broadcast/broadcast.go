package broadcast

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/types/tx"
)

func sendTx(ctx context.Context, txData []byte) error {
	grpcConn, ok := ctx.Value("rpc").(*grpc.ClientConn)
	if !ok {
		return fmt.Errorf("failed to retrieve rpc connection")
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
		return err
	}

	fmt.Println(grpcRes.TxResponse.Code) // Should be `0` if the tx is successful

	return nil
}
