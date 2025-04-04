package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Csy12139/Vesper/DN"
	"github.com/Csy12139/Vesper/MN"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"github.com/google/uuid"
)

func TestAddChunk(t *testing.T) {
	// Initialize log
	if err := log.InitLog("./logs", 10, 5, "info"); err != nil {
		t.Fatalf("Failed to initialize log: %v", err)
	}

	// Create temporary directories for MetaNode and DataNodes
	baseDir, err := os.MkdirTemp("", "vesper-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(baseDir)

	mnDir := filepath.Join(baseDir, "mn")
	if err := os.Mkdir(mnDir, 0755); err != nil {
		t.Fatalf("Failed to create MN directory: %v", err)
	}

	// Start MetaNode
	mnAddr := "127.0.0.1:8000"
	mn, err := MN.NewMNServer(mnAddr, mnDir)
	if err != nil {
		t.Fatalf("Failed to create MetaNode: %v", err)
	}
	mn.Start()

	// Start 5 DataNodes with different paths
	dataNodes := make([]*DN.DataNode, 5)
	for i := 0; i < 5; i++ {
		dnDir := filepath.Join(baseDir, fmt.Sprintf("dn%d", i))
		if err := os.Mkdir(dnDir, 0755); err != nil {
			t.Fatalf("Failed to create DN directory: %v", err)
		}

		// Generate a unique UUID for each DataNode
		dnUUID := uuid.New().String()

		dn, err := DN.NewDataNode(dnUUID, mnAddr, dnDir)
		if err != nil {
			t.Fatalf("Failed to create DataNode %d: %v", i, err)
		}
		dn.Start()
		dataNodes[i] = dn
	}

	t.Log("Successfully started MetaNode and 5 DataNodes")
	time.Sleep(10 * time.Second)
	mnClient, err := common.NewMNClient(mnAddr)
	if err != nil {
		t.Fatalf("Failed to create MetaNode client: %v", err)
	}
	response, err := mnClient.AddChunkMeta()
	if err != nil {
		t.Fatalf("Failed to add chunk meta: %v", err)
	}

	res, err := mnClient.AllocateDnForChunk(response.ChunkId, []string{})
	if err != nil {
		t.Fatalf("Failed to allocate dn for chunk: %v", err)
	}
	if res.Code != nil {
		t.Fatalf("Failed to allocate dn for chunk: %v", res.Code)
	}
	t.Logf("Allocated dn for chunk: %v", res.Uuid)
	res, err = mnClient.AllocateDnForChunk(response.ChunkId, []string{res.Uuid})
	if err != nil {
		t.Fatalf("Failed to allocate dn for chunk: %v", err)
	}
	if res.Code != nil {
		t.Fatalf("Failed to allocate dn for chunk: %v", res.Code)
	}
	t.Logf("Allocated dn for chunk: %v", res.Uuid)
	// Get the chunk meta to verify it's in CREATING state
	metaResp, err := mnClient.GetChunkMeta(response.ChunkId)
	if err != nil {
		t.Fatalf("Failed to get chunk meta: %v", err)
	}
	if metaResp.Meta.State != common.ChunkState_CREATING {
		t.Fatalf("Expected chunk state to be CREATING, got %v", metaResp.Meta.State)
	}
	
	// Complete the add chunk meta operation
	_, err = mnClient.CompleteAddChunkMeta(response.ChunkId)
	if err != nil {
		t.Fatalf("Failed to complete add chunk meta: %v", err)
	}
	t.Logf("Successfully completed add chunk meta for chunk ID: %v", response.ChunkId)
	// Verify the chunk state changed to CREATED
	metaResp, err = mnClient.GetChunkMeta(response.ChunkId)
	if err != nil {
		t.Fatalf("Failed to get chunk meta after completion: %v", err)
	}
	if metaResp.Meta.State != common.ChunkState_CREATED {
		t.Fatalf("Expected chunk state to be CREATED, got %v", metaResp.Meta.State)
	}
	
	mn.Stop()
	for _, dn := range dataNodes {
		dn.Stop()
	}
}
