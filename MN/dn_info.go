package MN

import (
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
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

func NewDataNodeInfo(uuid string) *DataNodeInfo {
	return &DataNodeInfo{
		UUID:                uuid,
		state:               ALIVE,
		cmdQueue:            make(chan *common.Command, 100),
		waitResponseCommand: make(map[uint64]*common.Command),
	}
}

func (dn *DataNodeInfo) SubmitAddChunkCmd(targetUUID string, timeout time.Duration, cb func(cmdId uint64, err error)) uint64 {
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
	log.Infof("Submit add chunk cmd id [%v] to cmdQueue, cmd info[%+v]", cmd.ID, cmd)
	dn.cmdQueue <- cmd
	log.Debugf("commandQueue addr [%v] of [%v] queue size [%d]", dn.cmdQueue, dn.UUID, len(dn.cmdQueue))
	return cmd.ID
}

func (dn *DataNodeInfo) SetDead() {
	dn.state = DEAD
}
func (dn *DataNodeInfo) SetALIVE() {
	dn.state = ALIVE
}
