package cosmos

import (
	"context"

	"google.golang.org/grpc"
)

func Handle(endpoint string, connections int, duration int, amount int, threads int) {
	ctx := context.Background()

	grpcConn, _ := grpc.Dial(
		endpoint,
		grpc.WithInsecure(),
	)
	defer grpcConn.Close()
	ctx = context.WithValue(ctx, "rpc", grpcConn)
}
