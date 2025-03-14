package grpcutil

import (
	"context"
	pb "github.com/Csy12139/Vesper/grpcutil/proto"
	"github.com/Csy12139/Vesper/log"
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
