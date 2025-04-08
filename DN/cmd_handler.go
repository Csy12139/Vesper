package DN

import (
	"context"
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"github.com/Csy12139/Vesper/p2p"
	"time"
)

// GetCommandResults returns all available results from the channel
func (dn *DataNode) GetCommandResults() []common.CommandResult {
	var results []common.CommandResult

	// Non-blocking read from channel
	for {
		select {
		case result := <-dn.resultChan:
			results = append(results, result)
		default:
			return results
		}
	}
}

// HandleCommand processes a command asynchronously
func (dn *DataNode) HandleCommand(cmd *common.Command) {
	go func() {
		result := common.CommandResult{
			CommandID: cmd.ID,
			Success:   true,
		}

		switch cmd.Type {
		case common.CommandType_READ_CHUNK:
			err := dn.handleReadChunk(cmd.ReadChunkCmd)
			if err != nil {
				result.Success = false
				result.ErrorMessage = err.Error()
			}

		case common.CommandType_ADD_CHUNK:
			err := dn.handleWriteChunk(cmd.AddChunkCmd)
			if err != nil {
				result.Success = false
				result.ErrorMessage = err.Error()
			}

		case common.CommandType_DELETE_CHUNK:
			err := dn.handleDeleteChunk(cmd.DeleteChunkCmd)
			if err != nil {
				result.Success = false
				result.ErrorMessage = err.Error()
			}

		default:
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("unknown command type: %v", cmd.Type)
		}

		// Send result through channel
		dn.resultChan <- result
	}()
}

func (dn *DataNode) handleReadChunk(cmd *common.ReadChunkCmd) error {
	if cmd == nil {
		return fmt.Errorf("read chunk command is nil")
	}
	log.Infof("Processing read chunk command for target UUID: %s", cmd.TargetUUID)
	// TODO: Implement read chunk logic
	return nil
}

func (dn *DataNode) handleWriteChunk(cmd *common.AddChunkCmd) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if cmd == nil {
		return fmt.Errorf("write chunk command is nil")
	}
	log.Infof("Processing write chunk command for target UUID: %s", cmd.TargetUUID)
	// TODO: Implement write chunk logic
	trans, err := p2p.NewDataTransfer(dn.UUID, dn.MNAddr)
	if err != nil {
		return fmt.Errorf("new data transfer error: %v", err)
	}
	// TODO: CRC校验
	err = trans.StartToReceive(cmd.TargetUUID, ctx, func(targetUUID string, label string, data []byte) {
		log.Info("Receive data !")
		for _, b := range data {
			if b != 0 {
				log.Errorf("[Receiver]The received data from %s is incorrect.", label)
			}
		}
		DataSize := 32 * 1024 * 1024
		if len(data) != DataSize {
			log.Error("[Receiver]The received data length is incorrect.")
		}
	})

	return nil
}

func (dn *DataNode) handleDeleteChunk(cmd *common.DeleteChunkCmd) error {
	if cmd == nil {
		return fmt.Errorf("delete chunk command is nil")
	}
	log.Infof("Processing delete chunk command for chunk ID: %d", cmd.ChunkID)
	// TODO: Implement delete chunk logic
	return nil
}
