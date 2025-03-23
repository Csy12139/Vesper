package main

import (
	"os"
	"time"

	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
)

// runMainLoop handles the main DN logic - sending heartbeats and processing commands
func runMainLoop() {
	mnClient, err := common.NewMNClient(GlobalConfig.MNAddr)
	if err != nil {
		log.Fatalf("Failed to create MN client: %v", err)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	cmdHandler := NewCommandHandler()

	for range ticker.C {
		// Get command results before sending heartbeat
		results := cmdHandler.GetResults()

		// Send heartbeat with results
		resp, err := mnClient.DoHeartbeat(GlobalConfig.UUID, results)
		if err != nil {
			log.Errorf("Heartbeat failed: %v", err)
			continue
		}

		// Process commands from heartbeat response
		for _, cmd := range resp.Commands {
			log.Infof("Received command: type=%v", cmd.Type)
			cmdHandler.HandleCommand(cmd)
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

	runMainLoop()
}
