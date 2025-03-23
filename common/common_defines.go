// CommandType represents different types of chunk operations
type CommandType int32

const (
	CommandType_NO_OP        CommandType = 0
	CommandType_READ_CHUNK   CommandType = 1
	CommandType_WRITE_CHUNK  CommandType = 2
	CommandType_DELETE_CHUNK CommandType = 3
)

// ChunkCommand represents a command to operate on a chunk
type ChunkCommand struct {
	CommandType CommandType
	ChunkID    uint64
	TargetUUID string
}

// HeartbeatRequest represents a heartbeat request from DN to MN
type HeartbeatRequest struct {
	UUID string
}

// HeartbeatResponse represents MN's response to a heartbeat
type HeartbeatResponse struct {
	Commands []ChunkCommand
}
