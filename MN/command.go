package MN

import (
	"github.com/Csy12139/Vesper/common"
	"sync"
	"time"
)

type CommandHelper struct {
	mu            sync.Mutex
	DN2WaitSubmit map[string]chan common.Command
	Submitted     map[string]chan common.Command
	DN2ID2Result  map[string]map[string]*common.CommandResult
	Id            uint64
}

func NewCommandHelper() *CommandHelper {
	ch := &CommandHelper{
		DN2WaitSubmit: make(map[string]chan common.Command),
		DN2ID2Result:  make(map[string]map[string]*common.CommandResult),
		Id:            0,
	}
	go ch.routine()
	return ch
}

func (cm *CommandHelper) SubmitAddChunkCmd(DNuuid string, TargetUUID string, cb func(err error)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.DN2WaitSubmit[DNuuid] == nil {
		cm.DN2WaitSubmit[DNuuid] = make(chan common.Command, 1024)
	}
	cm.Id++
	cm.DN2WaitSubmit[DNuuid] <- common.Command{
		ID:           cm.Id,
		Type:         common.CommandType_WRITE_CHUNK,
		ReadChunkCmd: nil,
		WriteChunkCmd: &common.WriteChunkCmd{
			TargetUUID: TargetUUID,
		},
		DeleteChunkCmd: nil,
		StartTimestamp: time.Now(),
		Timeout:        60 * time.Second,
		CallBack:       cb,
	}
}
func (cm *CommandHelper) GetCommands(DNuuid string) []common.Command {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cmds := make([]common.Command, 0)
	if cm.DN2WaitSubmit[DNuuid] == nil {
		return cmds
	}
	for {
		select {
		case cmd := <-cm.DN2WaitSubmit[DNuuid]:
			cmds = append(cmds, cmd)
		default:
			return cmds
		}
	}
}
func (cm *CommandHelper) routine() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		cm.mu.Lock()

	}
}
