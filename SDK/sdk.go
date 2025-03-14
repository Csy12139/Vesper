package main

import (
	"context"
	"github.com/Csy12139/Vesper/grpcutil"
	pb "github.com/Csy12139/Vesper/grpcutil/proto"
	"github.com/Csy12139/Vesper/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupClient(addr string) (*grpc.ClientConn, pb.MNClient) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect to %s: %v", addr, err)
	}
	c := pb.NewMNClient(conn)
	return conn, c
}

func createBucket(c pb.MNClient, bucketID int64) {
	createBucketFunc := func(ctx context.Context, client pb.MNClient) (interface{}, error) {
		return client.CreateBucketRequest(ctx, &pb.CreateBucket{BucketId: bucketID})
	}
	resp := grpcutil.Request(c, createBucketFunc)

	if createResp, ok := resp.(*pb.BoolResponse); ok {
		log.Info("Response: %s", createResp.Success)
	} else {
		log.Fatalf("unexpected response type")
	}
}
func main() {
	err := log.InitLog("./logs", 10, 5, "info")
	if err != nil {
		log.Fatalf("log init failed: %v", err)
	}
	conn, c := setupClient("localhost:50051")
	defer conn.Close()

	createBucket(c, 1001)
}
