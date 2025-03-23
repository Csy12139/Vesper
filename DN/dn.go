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
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config_file>", os.Args[0])
	}
	config, err := loadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := log.InitLog(config.LogPath, config.MaxSize, config.MaxBackups, config.LogLevel); err != nil {
		log.Fatalf("Failed to initialize log: %v", err)
	}
}
