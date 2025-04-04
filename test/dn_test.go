package main

// import (
// 	"fmt"
// 	main2 "github.com/Csy12139/Vesper/DNmain"
// 	"github.com/Csy12139/Vesper/common"
// 	"github.com/Csy12139/Vesper/log"
// 	"os"
// 	"testing"
// 	"time"
// )

// func equalSlices(a, b []string) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i := range a {
// 		if a[i] != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }
// func TestPutAndGet(t *testing.T) {
// 	if len(os.Args) < 2 {
// 		fmt.Printf("Usage: %s <config_file>", os.Args[0])
// 	}

// 	err := main2.loadConfig(os.Args[1])
// 	if err != nil {
// 		fmt.Printf("Failed to load config: %v", err)
// 	}

// 	if err := log.InitLog(main2.GlobalConfig.Log.LogDir, main2.GlobalConfig.Log.MaxFileSizeMb, main2.GlobalConfig.Log.MaxFileNum, main2.GlobalConfig.Log.LogLevel); err != nil {
// 		log.Fatalf("Failed to initialize log: %v", err)
// 	}
// 	putRequest := common.PutSDPCandidatesRequest{
// 		SourceUUID: "1",
// 		TargetUUID: "2",
// 		SDP:        "abc",
// 		Candidates: []string{"a", "b", "c"},
// 	}
// 	put(&putRequest)

// 	time.Sleep(1 * time.Second)

// 	sdp, candidates := get(&common.GetSDPCandidatesRequest{
// 		SourceUUID: "1",
// 		TargetUUID: "2",
// 	})
// 	if putRequest.SDP != sdp {
// 		t.Errorf("get SDP = %s; want %s", sdp, putRequest.SDP)
// 	}
// 	if !equalSlices(putRequest.Candidates, candidates) {
// 		t.Errorf("get candidates = %s; want %s", candidates, putRequest.SDP)
// 	}
// }
