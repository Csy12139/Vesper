package main

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"os"
	"path/filepath"

	"github.com/Csy12139/Vesper/SDK"
	"github.com/Csy12139/Vesper/log"
)

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
	sdk, err := SDK.NewSdk(GlobalConfig.UUID, GlobalConfig.MNAddr, GlobalConfig.DataPath)
	if err != nil {
		log.Fatalf("Failed to initialize sdk: %v", err)
	}
	// generate file
	filePath := filepath.Join(GlobalConfig.DataPath, "testFile")
	file, err := os.Create(filePath)
	if err != nil {
		log.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	err = file.Truncate(33 * 1024 * 1024)
	if err != nil {
		log.Errorf("failed to set file size: %v", err)
	}
	SendMd5, err := common.CalculateMD5(filePath)
	if err != nil {
		log.Errorf("failed to calculate MD5: %v", err)
	}
	// store
	err = sdk.StoreFile("testFile", 1, false, 1)
	if err != nil {
		log.Errorf("failed to store file: %v", err)
	}
	log.Debug("Successfully stored file testFile")
	// remove file
	err = os.Remove(filePath)
	if err != nil {
		log.Errorf("failed to remove file: %v", err)
	}
	// read
	err = sdk.ReadFile("testFile", 1)
	if err != nil {
		log.Errorf("failed to read file: %v", err)
	}
	log.Debug("Successfully read file testFile")
	ReceiveMd5, err := common.CalculateMD5(filePath)
	if err != nil {
		log.Errorf("failed to calculate MD5: %v", err)
	}
	log.Debug("Successfully calculated SendMd5: ", SendMd5)
	log.Debug("Successfully calculated ReceiveMd5: ", ReceiveMd5)
	if SendMd5 == ReceiveMd5 {
		log.Infof("The MD5 is the same!")
	} else {
		log.Infof("The MD5 is not the same.")
	}
	return
}
