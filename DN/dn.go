package main

import (
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"os"
	"time"
)

// runMainLoop handles the main DN logic - sending heartbeats and processing commands
func runMainLoop(mnClient *common.MNClient) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send heartbeat
			resp, err := mnClient.DoHeartbeat(GlobalConfig.UUID)
			if err != nil {
				log.Errorf("Heartbeat failed: %v", err)
				continue
			}

			// Process commands from heartbeat response
			for _, cmd := range resp.Commands {
				log.Infof("Received command: type=%v, chunkID=%v, targetUUID=%v", 
					cmd.CommandType, cmd.ChunkID, cmd.TargetUUID)
				// TODO: Implement command execution
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config_file>", os.Args[0])
	}

	err := loadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := log.InitLog(GlobalConfig.Log.LogDir, GlobalConfig.Log.MaxFileSizeMb, GlobalConfig.Log.MaxFileNum, GlobalConfig.Log.LogLevel); err != nil {
		log.Fatalf("Failed to initialize log: %v", err)
	}

	// Create MN client
	mnClient := common.NewMNClient(GlobalConfig.MNAddr)
	defer mnClient.Close()

	// Start main loop
	runMainLoop(mnClient)
}
