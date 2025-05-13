package DN

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"os"
	"path/filepath"
	"strconv"
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
	filePath := filepath.Join(dn.dataPath, strconv.FormatUint(cmd.ChunkID, 10))
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file error [%v]", err)
	}
	defer file.Close()
	data := make([]byte, common.ChunkSize)
	_, err = file.Read(data)
	if err != nil {
		return fmt.Errorf("read file error [%v]", err)
	}
	//TODO:crc 校验
	err = dn.dataTransfer.Send(cmd.TargetUUID, strconv.FormatUint(cmd.ChunkID, 10), data)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	return nil
}

func (dn *DataNode) handleWriteChunk(cmd *common.AddChunkCmd) error {
	if cmd == nil {
		return fmt.Errorf("write chunk command is nil")
	}
	log.Infof("Processing write chunk command for target UUID: %s", cmd.TargetUUID)
	// TODO: CRC校验
	err := dn.dataTransfer.StartToReceive(cmd.TargetUUID, func(targetUUID string, label string, data []byte) {
		log.Info("Receive data !")
		//for _, b := range data {
		//	if b != 0 {
		//		log.Errorf("[Receiver]The received data from %s is incorrect.", label)
		//	}
		//}
		if len(data) != common.ChunkSize {
			log.Error("[Receiver]The received data length is incorrect.")
		} else {
			filePath := filepath.Join(dn.dataPath, label)
			file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
			defer file.Close()
			if err != nil {
				log.Error("[Receiver]The file open failed !")
			}
			_, err = file.Write(data)
			if err != nil {
				log.Error("[Receiver]The chunk write failed !")
			}
			// Complete the add chunk meta operation
			chunkId, err := strconv.ParseUint(label, 10, 64)
			if err != nil {
				log.Errorf("read chunk id from file error [%v]", err)
			}
			log.Debug("Start to run CompleteAddChunkMeta cmd, chunkId:", chunkId)
			err = dn.mnClient.CompleteAddChunkOnDN(chunkId, dn.UUID)
			if err != nil {
				log.Errorf("Failed to complete add chunk meta: %v", err)
				return
			}
		}
	})
	if err != nil {
		return fmt.Errorf("write chunk command error: %v", err)
	}
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
