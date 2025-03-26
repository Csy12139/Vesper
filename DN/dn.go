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
	var results []common.CommandResult

	for range ticker.C {
		// Get command results before sending heartbeat
		// TODO(cissy) here is a bug, where heartbeat failed, results loss
		results = append(results, cmdHandler.GetResults()...)
		// Send heartbeat with results
		resp, err := mnClient.DoHeartbeat(GlobalConfig.UUID, results)
		if err != nil {
			log.Errorf("Heartbeat failed: %v", err)
			continue
		}
		results = nil
		// Process commands from heartbeat response
		for _, cmd := range resp.Commands {
			log.Infof("Received command: type=%v", cmd.Type)
			cmdHandler.HandleCommand(cmd)
		}
	}
}

func put(req *common.PutSDPCandidatesRequest) {
	mnClient, err := common.NewMNClient(GlobalConfig.MNAddr)
	if err != nil {
		log.Fatalf("Failed to create MN client: %v", err)
	}
	resp, err := mnClient.PutSDPCandidates(req)
	if err != nil {
		log.Errorf("PutSDPCandidates failed: %v", err)
	}
	log.Infof("PutSDPCandidates: %v", resp)
}

func get(req *common.GetSDPCandidatesRequest) (string, []string) {
	mnClient, err := common.NewMNClient(GlobalConfig.MNAddr)
	if err != nil {
		log.Fatalf("Failed to create MN client: %v", err)
	}
	resp, err := mnClient.GetSDPCandidates(req)
	if err != nil {
		log.Errorf("GetSDPCandidates failed: %v", err)
	}
	log.Infof("GetSDPCandidates: %v", resp)
	return resp.SDP, resp.Candidates
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

	runMainLoop()
}
