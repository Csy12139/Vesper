package grpcutil

import (
	"context"
	pb "github.com/Csy12139/Vesper/grpcutil/proto"
	"github.com/Csy12139/Vesper/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

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

func SetupClient(addr string) (*grpc.ClientConn, pb.MNClient) {
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
