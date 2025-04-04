package MN

import (
	"github.com/Csy12139/Vesper/common"
	"sync"
	"sync/atomic"
	"time"
)

type DataNodeState string

const (
	DEAD  DataNodeState = "DEAD"
	ALIVE DataNodeState = "ALIVE"
)

var cmdId atomic.Uint64

type DataNodeInfo struct {
	UUID                string
	state               DataNodeState
	cmdQueue            chan *common.Command
	waitResponseCommand map[uint64]*common.Command
	lastHeartbeatTime   time.Time
	mutex               sync.Mutex
}

func NewDataNodeInfo() *DataNodeInfo {
	return &DataNodeInfo{
		state:               ALIVE,
		cmdQueue:            make(chan *common.Command, 100),
		waitResponseCommand: make(map[uint64]*common.Command),
	}
}

func (dn *DataNodeInfo) SubmitAddChunkCmd(targetUUID string, timeout time.Duration, cb func(cmdId uint64, err error)) uint64{
	cmd := &common.Command{
		ID:             cmdId.Add(1),
		Type:           common.CommandType_ADD_CHUNK,
		ReadChunkCmd:   nil,
		AddChunkCmd:    &common.AddChunkCmd{TargetUUID: targetUUID},
		DeleteChunkCmd: nil,
		StartTimestamp: time.Now(),
		Timeout:        timeout,
		CallBack:       cb,
	}
	dn.cmdQueue <- cmd
	return cmd.ID
}

func (dn *DataNodeInfo) SetDead() {
	dn.state = DEAD
}
func (dn *DataNodeInfo) SetALIVE() {
	dn.state = ALIVE
}
