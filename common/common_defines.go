// CommandType represents different types of chunk operations
package common

import (
	"errors"
	"time"
)

type CommandType int32

const (
	CommandType_NO_OP        CommandType = 0
	CommandType_READ_CHUNK   CommandType = 1
	CommandType_ADD_CHUNK    CommandType = 2
	CommandType_DELETE_CHUNK CommandType = 3
)

// CommandResult represents the result of executing a command
type CommandResult struct {
	CommandID    uint64
	Success      bool
	ErrorMessage string
}

// ReadChunkCmd represents a command to read a chunk
type ReadChunkCmd struct {
	TargetUUID string
}

// AddChunkCmd represents a command to write a chunk
type AddChunkCmd struct {
	TargetUUID string
}

// DeleteChunkCmd represents a command to delete a chunk
type DeleteChunkCmd struct {
	ChunkID uint64
}

// Command represents a command from MN to DN
type Command struct {
	ID             uint64
	Type           CommandType
	ReadChunkCmd   *ReadChunkCmd
	AddChunkCmd    *AddChunkCmd
	DeleteChunkCmd *DeleteChunkCmd
	StartTimestamp time.Time
	Timeout        time.Duration
	CallBack       func(cmdId uint64, err error)
}

// HeartbeatRequest represents a heartbeat request from DN to MN
type HeartbeatRequest struct {
	UUID           string
	CommandResults []CommandResult
}

// HeartbeatResponse represents MN's response to a heartbeat
type HeartbeatResponse struct {
	Commands []*Command
}

type PutSDPCandidatesRequest struct {
	SourceUUID string
	TargetUUID string
	SDP        string
	Candidates []string
}

type PutSDPCandidatesResponse struct {
	Success      bool
	ErrorMessage string
}

type GetSDPCandidatesRequest struct {
	SourceUUID string
	TargetUUID string
}

type GetSDPCandidatesResponse struct {
	Success      bool
	ErrorMessage string
	SourceUUID   string
	TargetUUID   string
	SDP          string
	Candidates   []string
}


type AddChunkMetaRequest struct {
}

type AddChunkMetaResponse struct {
	ChunkId uint64
}

type AllocateDnForChunkRequest struct {
	ChunkId  uint64
	Excludes []string
}

type AllocateDnForChunkResponse struct {
	Uuid string
	Code error
}

type AddChunkOnDNRequest struct {
	SdkUuid string
	DnUuid  string
}

type AddChunkOnDNResponse struct {
	Code error
}

type CompleteAddChunkOnDNRequest struct {
	ChunkId uint64
	Uuid    string
}

type CompleteAddChunkOnDNResponse struct {
	Code error
}

type CompleteAddChunkMetaRequest struct {
	ChunkID uint64
}

type CompleteAddChunkMetaResponse struct {
}

type GetChunkMetaRequest struct {
	ChunkID uint64
}

type GetChunkMetaResponse struct {
	Meta *ChunkMeta
	Code error
}



type ChunkState int

const (
	ChunkState_CREATING ChunkState = 0
	ChunkState_CREATED  ChunkState = 1
)

// ChunkSize defines the size of each chunk in bytes (32 MiB)
const ChunkSize = 32 * 1024 * 1024

// Chunk represents a chunk of data in memory
type Chunk struct {
	// ID uniquely identifies the chunk
	ID uint64
	// Data contains the actual chunk data
	Data []byte
}

type ChunkMeta struct {
	ID      uint64
	State   ChunkState
	DnUuids map[string]struct{}
}

var (
	ErrChunkNotFound = errors.New("chunk not found")
	ErrNoAvailableDN = errors.New("no available data node")
	ErrCommitDNTimeout = errors.New("commit dn timeout")
	ErrDNNotFound = errors.New("data node not found")
)
