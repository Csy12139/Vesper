package DN

import (
	"fmt"
	"github.com/Csy12139/Vesper/p2p"
	"sync/atomic"
	"time"

	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
)

type DataNode struct {
	UUID         string
	stop         atomic.Bool
	MNAddr       string
	mnClient     *common.MNClient
	resultChan   chan common.CommandResult
	dataTransfer *p2p.DataTransfer
	dataPath     string
}

func NewDataNode(UUID string, MNAddr string, DataPath string) (*DataNode, error) {
	mnClient, err := common.NewMNClient(MNAddr)
	if err != nil {
		return nil, fmt.Errorf("NewMNClient error:%v", err)
	}
	trans, err := p2p.NewDataTransfer(UUID, MNAddr)
	if err != nil {
		return nil, fmt.Errorf("NewDataTransfer error:%v", err)
	}
	dn := DataNode{
		UUID:         UUID,
		MNAddr:       MNAddr,
		mnClient:     mnClient,
		resultChan:   make(chan common.CommandResult, 100),
		dataPath:     DataPath,
		dataTransfer: trans,
	}
	dn.stop.Store(true)
	return &dn, nil
}

func (dn *DataNode) Start() {
	dn.stop.Store(false)
	go dn.doHeartbeatLoop()
}

func (dn *DataNode) doHeartbeatLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	results := make([]common.CommandResult, 0)
	for range ticker.C {
		if dn.stop.Load() {
			return
		}
		results = append(results, dn.GetCommandResults()...)
		resp, err := dn.mnClient.DoHeartbeat(dn.UUID, results)
		if err != nil {
			log.Errorf("Heartbeat failed: %v", err)
			continue
		}
		results = make([]common.CommandResult, 0)
		// Process commands from heartbeat response
		for _, cmd := range resp.Commands {
			log.Infof("Received command: type=%v", cmd.Type)
			dn.HandleCommand(cmd)

		}
	}
}

func (dn *DataNode) updateReceivingRateLoop() {
	timeInterval := 5 * time.Second
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()
	for range ticker.C {
		totalBytes := dn.dataTransfer.GetTotalBytes()
		if totalBytes != 0 {
			currentRate := float64(totalBytes) / timeInterval.Seconds()
			err := dn.mnClient.UpdateReceivingRate(currentRate)
			if err != nil {
				log.Errorf("UpdateReceivingRate error:%v", err)
			}
		}
	}
}

func (dn *DataNode) IsRunning() bool {
	return !dn.stop.Load()
}

func (dn *DataNode) Stop() {
	dn.stop.Store(true)
}
