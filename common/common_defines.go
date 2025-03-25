// CommandType represents different types of chunk operations
package common

type CommandType int32

const (
	CommandType_NO_OP        CommandType = 0
	CommandType_READ_CHUNK   CommandType = 1
	CommandType_WRITE_CHUNK  CommandType = 2
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

// WriteChunkCmd represents a command to write a chunk
type WriteChunkCmd struct {
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
	WriteChunkCmd  *WriteChunkCmd
	DeleteChunkCmd *DeleteChunkCmd
}

// HeartbeatRequest represents a heartbeat request from DN to MN
type HeartbeatRequest struct {
	UUID           string
	CommandResults []CommandResult
}

// HeartbeatResponse represents MN's response to a heartbeat
type HeartbeatResponse struct {
	Commands []Command
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
	SourceUUID string
	TargetUUID string
	SDP        string
	Candidates []string
}
