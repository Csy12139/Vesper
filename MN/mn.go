package main

import (
	"github.com/Csy12139/Vesper/log"
)

func main() {
	err := log.InitLog("./logs", 10, 5, "info")
	if err != nil {
		log.Fatalf("log init failed: %v", err)
	}

}
