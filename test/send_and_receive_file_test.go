package test

import (
	"fmt"
	"github.com/Csy12139/Vesper/DN"
	"github.com/Csy12139/Vesper/MN"
	"github.com/Csy12139/Vesper/SDK"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSendAndReceiveFile(t *testing.T) {
	// Initialize log
	if err := log.InitLog("./logs", 10, 5, "debug"); err != nil {
		t.Fatalf("Failed to initialize log: %v", err)
	}
	// Create temporary directories for MetaNode ,SDK and DataNodes
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
	log.Info("Successfully started MetaNode and 5 DataNodes")
	time.Sleep(5 * time.Second)
	// Start 1 SDK
	sdkDir := filepath.Join(baseDir, "sdk")
	if err := os.Mkdir(sdkDir, 0755); err != nil {
		t.Fatalf("Failed to create sdk directory: %v", err)
	}
	sdkUUID := uuid.New().String()
	log.Infof("sdk uuid: [%s]", sdkUUID)
	sdk, err := SDK.NewSdk(sdkUUID, mnAddr, sdkDir)
	if err != nil {
		t.Fatalf("Failed to create Sdk: %v", err)
	}
	// generate file
	filePath := filepath.Join(sdkDir, "testFile")
	file, err := os.Create(filePath)
	if err != nil {
		t.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	err = file.Truncate(1 * 1024 * 1024)
	if err != nil {
		t.Errorf("failed to set file size: %v", err)
	}
	SendMd5, err := common.CalculateMD5(filePath)
	if err != nil {
		t.Errorf("failed to calculate MD5: %v", err)
	}
	// store
	err = sdk.StoreFile("testFile", 1, false, 1)
	if err != nil {
		t.Errorf("failed to store file: %v", err)
	}
	log.Debug("Successfully stored file testFile")
	// remove file
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("failed to remove file: %v", err)
	}
	// read
	err = sdk.ReadFile("testFile", 1)
	if err != nil {
		t.Errorf("failed to read file: %v", err)
	}
	log.Debug("Successfully read file testFile")
	ReceiveMd5, err := common.CalculateMD5(filePath)
	if err != nil {
		t.Errorf("failed to calculate MD5: %v", err)
	}
	log.Debug("Successfully calculated SendMd5: ", SendMd5)
	log.Debug("Successfully calculated ReceiveMd5: ", ReceiveMd5)
	if SendMd5 == ReceiveMd5 {
		log.Infof("The MD5 is the same!")
		return
	}
}
