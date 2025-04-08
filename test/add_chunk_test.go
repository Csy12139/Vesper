package test

import (
	"context"
	"fmt"
	"github.com/Csy12139/Vesper/p2p"
	"os"
	"path/filepath"
	"strconv"
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
	if err := log.InitLog("./logs", 10, 5, "debug"); err != nil {
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
	defer mn.Stop()

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
	defer func() {
		for _, dn := range dataNodes {
			dn.Stop()
		}
	}()

	t.Log("Successfully started MetaNode and 5 DataNodes")
	time.Sleep(5 * time.Second)

	sdkUUID := uuid.New().String()
	mnClient, err := common.NewMNClient(mnAddr)
	if err != nil {
		t.Fatalf("Failed to create MetaNode client: %v", err)
	}

	var excludes []string
	for i := 0; i < 5; i++ {
		chunkId, err := mnClient.AddChunkMeta()
		if err != nil {
			t.Fatalf("Failed to add chunk meta: %v", err)
		}
		dnId, err := mnClient.AllocateDnForChunk(chunkId, excludes)
		if err != nil {
			t.Fatalf("Failed to allocate dn for chunk: %v", err)
		}
		t.Logf("Allocated dn %v for chunk %v", dnId, chunkId)
		excludes = append(excludes, dnId)
		err = mnClient.AddChunkOnDN(sdkUUID, dnId)
		if err != nil {
			t.Fatalf("Failed to add chunk on dn: %v", err)
		}
		time.Sleep(5 * time.Second)
		trans, err := p2p.NewDataTransfer(sdkUUID, mnAddr)
		if err != nil {
			t.Fatalf("Failed to create DataTransfer: %v", err)
		}
		data := make([]byte, 32*1024*1024)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			err = trans.Send(dnId, strconv.FormatUint(chunkId, 10), data, ctx)
			if err != nil {
				t.Errorf("Failed to send data: %v", err)
				return
			}
			// Get the chunk meta to verify it's in CREATING state
			metaResp, err := mnClient.GetChunkMeta(chunkId)
			if err != nil {
				t.Errorf("Failed to get chunk meta: %v", err)
				return
			}
			if metaResp.State != common.ChunkState_CREATING {
				t.Errorf("Expected chunk state to be CREATING, got %v", metaResp.State)
				return
			}

			// Complete the add chunk meta operation
			err = mnClient.CompleteAddChunkMeta(chunkId)
			if err != nil {
				t.Errorf("Failed to complete add chunk meta: %v", err)
				return
			}
			t.Logf("Successfully completed add chunk meta for chunk ID: %v", chunkId)
			// Verify the chunk state changed to CREATED
			metaResp, err = mnClient.GetChunkMeta(chunkId)
			if err != nil {
				t.Errorf("Failed to get chunk meta after completion: %v", err)
				return
			}
			if metaResp.State != common.ChunkState_CREATED {
				t.Errorf("Expected chunk state to be CREATED, got %v", metaResp.State)
				return
			}
		}()
	}
}
