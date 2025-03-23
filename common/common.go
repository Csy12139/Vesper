package common

import (
	"os"
	"path/filepath"
)

func GetExecName() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execName := filepath.Base(execPath)
	// log.Infof("execName:[%v]", execName)
	return execName, nil
}

// ChunkSize defines the size of each chunk in bytes (32 MiB)
const ChunkSize = 32 * 1024 * 1024

// Chunk represents a chunk of data in memory
type Chunk struct {
	// ID uniquely identifies the chunk
	ID uint64
	// Data contains the actual chunk data
	Data []byte
}
