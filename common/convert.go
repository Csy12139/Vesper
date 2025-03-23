package common

import (
	pb "github.com/Csy12139/Vesper/proto"
)

// Proto2Command converts a protobuf ChunkCommand to internal format
func Proto2Command(pbCmd *pb.ChunkCommand) ChunkCommand {
	return ChunkCommand{
		CommandType: CommandType(pbCmd.CommandType),
		ChunkID:    pbCmd.ChunkId,
		TargetUUID: pbCmd.TargetUuid,
	}
}

// Command2Proto converts a ChunkCommand to protobuf format
func Command2Proto(cmd ChunkCommand) *pb.ChunkCommand {
	return &pb.ChunkCommand{
		CommandType: pb.CommandType(cmd.CommandType),
		ChunkId:    cmd.ChunkID,
		TargetUuid: cmd.TargetUUID,
	}
}

// Proto2HeartbeatRequest converts a protobuf HeartbeatRequest to internal format
func Proto2HeartbeatRequest(pbReq *pb.HeartbeatRequest) HeartbeatRequest {
	return HeartbeatRequest{
		UUID: pbReq.Uuid,
	}
}

// HeartbeatRequest2Proto converts a HeartbeatRequest to protobuf format
func HeartbeatRequest2Proto(req HeartbeatRequest) *pb.HeartbeatRequest {
	return &pb.HeartbeatRequest{
		Uuid: req.UUID,
	}
}

// Proto2HeartbeatResponse converts a protobuf HeartbeatResponse to internal format
func Proto2HeartbeatResponse(pbResp *pb.HeartbeatResponse) HeartbeatResponse {
	commands := make([]ChunkCommand, len(pbResp.Commands))
	for i, pbCmd := range pbResp.Commands {
		commands[i] = Proto2Command(pbCmd)
	}
	return HeartbeatResponse{
		Commands: commands,
	}
}

// HeartbeatResponse2Proto converts a HeartbeatResponse to protobuf format
func HeartbeatResponse2Proto(resp HeartbeatResponse) *pb.HeartbeatResponse {
	pbCommands := make([]*pb.ChunkCommand, len(resp.Commands))
	for i, cmd := range resp.Commands {
		pbCommands[i] = Command2Proto(cmd)
	}
	return &pb.HeartbeatResponse{
		Commands: pbCommands,
	}
}

