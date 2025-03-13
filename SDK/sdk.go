package main

import (
	"context"
	pb "github.com/Csy12139/Vesper/common/communicate"
	"github.com/Csy12139/Vesper/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
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

type RequestFunc func(context.Context, pb.MNClient) (interface{}, error)

func Request(c pb.MNClient, requestFunc RequestFunc) interface{} {
	const maxAttempts = 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		resp, err := requestFunc(ctx, c)
		if err == nil {
			return resp
		}
		log.Infof("Attempt %d failed: %v", attempt, err)
		if attempt < maxAttempts {
			time.Sleep(time.Second * time.Duration(attempt))
		}
	}
	log.Fatalf("Failed after %d attempts", maxAttempts)
	return nil
}
func createBucket(c pb.MNClient, bucketID int64) {
	createBucketFunc := func(ctx context.Context, client pb.MNClient) (interface{}, error) {
		return client.CreateBucketRequest(ctx, &pb.CreateBucket{BucketId: bucketID})
	}
	resp := Request(c, createBucketFunc)
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
