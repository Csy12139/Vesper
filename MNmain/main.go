package MNmain

import (
	"fmt"
	"os"
	"time"

	"github.com/Csy12139/Vesper/MN"
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
	server, err := MN.NewMNServer(GlobalConfig.MNAddr, GlobalConfig.MNDataPath)
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
	for server.IsRunning() {
		time.Sleep(10 * time.Second)
	}
}
