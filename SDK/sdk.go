package main

import (
	"context"
	"github.com/Csy12139/Vesper/grpcutil"
	pb "github.com/Csy12139/Vesper/grpcutil/proto"
	"github.com/Csy12139/Vesper/log"
)

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
	conn, c := grpcutil.SetupClient("localhost:50051")
	defer conn.Close()

	createBucket(c, 1001)
}
