package main

import (
	"fmt"
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

func test_put() {
	mnClient, err := common.NewMNClient(GlobalConfig.MNAddr)
	if err != nil {
		log.Fatalf("Failed to create MN client: %v", err)
	}
	req := common.PutSDPCandidatesRequest{
		SourceUUID: "1",
		TargetUUID: "2",
		SDP:        "abc",
		Candidates: nil,
	}
	resp, err := mnClient.PutSDPCandidates(&req)
	if err != nil {
		log.Errorf("PutSDPCandidates failed: %v", err)
	}
	fmt.Println(resp)
}

func test_get() {
	mnClient, err := common.NewMNClient(GlobalConfig.MNAddr)
	if err != nil {
		log.Fatalf("Failed to create MN client: %v", err)
	}
	req := common.GetSDPCandidatesRequest{
		SourceUUID: "1",
		TargetUUID: "2",
	}
	resp, err := mnClient.GetSDPCandidates(&req)
	if err != nil {
		log.Errorf("GetSDPCandidates failed: %v", err)
	}
	fmt.Println(resp)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <config_file>", os.Args[0])
	}

	err := loadConfig(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
	}

	if err := log.InitLog(GlobalConfig.Log.LogDir, GlobalConfig.Log.MaxFileSizeMb, GlobalConfig.Log.MaxFileNum, GlobalConfig.Log.LogLevel); err != nil {
		log.Fatalf("Failed to initialize log: %v", err)
	}

	//runMainLoop()
	test_put()
	test_get()
}
