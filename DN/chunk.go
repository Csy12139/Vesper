package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Csy12139/Vesper/common"
)

// AddChunk creates a new chunk with the given chunk
func AddChunk(chunk *common.Chunk) error {
	// Verify data size
	if len(chunk.Data) != common.ChunkSize {
		return fmt.Errorf("invalid chunk size: got %d bytes, want %d bytes", len(chunk.Data), common.ChunkSize)
	}

	// Create chunk file path using configured data directory
	chunkPath := filepath.Join(GlobalConfig.DataPath, fmt.Sprintf("%d", chunk.ID))

	// Create file with O_EXCL flag to ensure atomic creation
	file, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("chunk %d already exists", chunk.ID)
		}
		return fmt.Errorf("failed to create chunk file: %w", err)
	}
	defer file.Close()

	// Write data to file
	_, err = file.Write(chunk.Data)
	if err != nil {
		// Try to cleanup the file if write fails
		os.Remove(chunkPath)
		return fmt.Errorf("failed to write chunk data: %w", err)
	}

	return nil
}

// GetChunk reads and returns the chunk with the given ID
func GetChunk(id uint64) (*common.Chunk, error) {
	chunkPath := filepath.Join(GlobalConfig.DataPath, fmt.Sprintf("%d", id))
	
	data, err := os.ReadFile(chunkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("chunk %d does not exist", id)
		}
		return nil, fmt.Errorf("failed to read chunk: %w", err)
	}

	if len(data) != common.ChunkSize {
		return nil, fmt.Errorf("corrupted chunk: got %d bytes, want %d bytes", len(data), common.ChunkSize)
	}

	return &common.Chunk{
		ID:   id,
		Data: data,
	}, nil
}

// RemoveChunk deletes the chunk with the given ID
func RemoveChunk(id uint64) error {
	chunkPath := filepath.Join(GlobalConfig.DataPath, fmt.Sprintf("%d", id))
	
	err := os.Remove(chunkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("chunk %d does not exist", id)
		}
		return fmt.Errorf("failed to remove chunk: %w", err)
	}
	
	return nil
}


