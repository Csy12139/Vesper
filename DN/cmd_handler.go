package main

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
)

// CommandHandler handles execution of commands from MN
type CommandHandler struct {
	resultChan chan common.CommandResult
}

// NewCommandHandler creates a new command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		resultChan: make(chan common.CommandResult, 1024), // Buffered channel to avoid blocking
	}
}

// GetResults returns all available results from the channel
func (h *CommandHandler) GetResults() []common.CommandResult {
	var results []common.CommandResult

	// Non-blocking read from channel
	for {
		select {
		case result := <-h.resultChan:
			results = append(results, result)
		default:
			return results
		}
	}
}

// HandleCommand processes a command asynchronously
func (h *CommandHandler) HandleCommand(cmd common.Command) {
	go func() {
		result := common.CommandResult{
			CommandID: cmd.ID,
			Success:   true,
		}

		switch cmd.Type {
		case common.CommandType_READ_CHUNK:
			err := h.handleReadChunk(cmd.ReadChunkCmd)
			if err != nil {
				result.Success = false
				result.ErrorMessage = err.Error()
			}

		case common.CommandType_WRITE_CHUNK:
			err := h.handleWriteChunk(cmd.WriteChunkCmd)
			if err != nil {
				result.Success = false
				result.ErrorMessage = err.Error()
			}

		case common.CommandType_DELETE_CHUNK:
			err := h.handleDeleteChunk(cmd.DeleteChunkCmd)
			if err != nil {
				result.Success = false
				result.ErrorMessage = err.Error()
			}

		default:
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("unknown command type: %v", cmd.Type)
		}

		// Send result through channel
		h.resultChan <- result
	}()
}

func (h *CommandHandler) handleReadChunk(cmd *common.ReadChunkCmd) error {
	if cmd == nil {
		return fmt.Errorf("read chunk command is nil")
	}
	log.Infof("Processing read chunk command for target UUID: %s", cmd.TargetUUID)
	// TODO: Implement read chunk logic
	return nil
}

func (h *CommandHandler) handleWriteChunk(cmd *common.WriteChunkCmd) error {
	if cmd == nil {
		return fmt.Errorf("write chunk command is nil")
	}
	log.Infof("Processing write chunk command for target UUID: %s", cmd.TargetUUID)
	// TODO: Implement write chunk logic
	return nil
}

func (h *CommandHandler) handleDeleteChunk(cmd *common.DeleteChunkCmd) error {
	if cmd == nil {
		return fmt.Errorf("delete chunk command is nil")
	}
	log.Infof("Processing delete chunk command for chunk ID: %d", cmd.ChunkID)
	// TODO: Implement delete chunk logic
	return nil
}



