package main

import (
	"context"
	"github.com/Csy12139/Vesper/grpcutil"
	pb "github.com/Csy12139/Vesper/grpcutil/proto"
	"github.com/Csy12139/Vesper/log"
	"os"
	"time"
)

func heartDetect(c pb.MNClient, isRecover bool) {
	ticker := time.NewTicker(5 * time.Second)
	//defer ticker.Stop()

	go func() {
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			resp, err := c.HeartDetectRequest(ctx, &pb.HeartbeatDetection{IsRecover: isRecover})
			if err != nil {
				log.Infof("Heartbeat failed: %v", err)
				continue
			}
			log.Infof("Heartbeat Response: success=%v, isRecover=%v, isReport=%v, isReceiverData=%v",
				resp.Success, resp.IsRecover, resp.IsReport, resp.IsReceiverData)
		}
	}()
}

func main() {
	loadConfig(os.Args[1])
	err := log.InitLog("./logs", 10, 5, "info")
	if err != nil {
		log.Fatalf("log init failed: %v", err)
	}
	conn, c := grpcutil.SetupClient("localhost:50051")
	defer conn.Close()

	heartDetect(c, false)

	select {}
}
