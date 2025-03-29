package common

import (
	pb "github.com/Csy12139/Vesper/proto"
)

// Proto2Command converts a protobuf Command to internal format
func Proto2Command(pbCmd *pb.Command) *Command {
	cmd := &Command{
		ID:   pbCmd.Id,
		Type: CommandType(pbCmd.Type),
	}

	switch c := pbCmd.Command.(type) {
	case *pb.Command_ReadChunkCmd:
		cmd.ReadChunkCmd = &ReadChunkCmd{
			TargetUUID: c.ReadChunkCmd.TargetUuid,
		}
	case *pb.Command_WriteChunkCmd:
		cmd.AddChunkCmd = &AddChunkCmd{
			TargetUUID: c.WriteChunkCmd.TargetUuid,
		}
	case *pb.Command_DeleteChunkCmd:
		cmd.DeleteChunkCmd = &DeleteChunkCmd{
			ChunkID: c.DeleteChunkCmd.ChunkId,
		}
	}

	return cmd
}

// Command2Proto converts a Command to protobuf format
func Command2Proto(cmd *Command) *pb.Command {
	pbCmd := &pb.Command{
		Id:   cmd.ID,
		Type: pb.CommandType(cmd.Type),
	}

	switch {
	case cmd.ReadChunkCmd != nil:
		pbCmd.Command = &pb.Command_ReadChunkCmd{
			ReadChunkCmd: &pb.ReadChunkCmd{
				TargetUuid: cmd.ReadChunkCmd.TargetUUID,
			},
		}
	case cmd.AddChunkCmd != nil:
		pbCmd.Command = &pb.Command_WriteChunkCmd{
			WriteChunkCmd: &pb.WriteChunkCmd{
				TargetUuid: cmd.AddChunkCmd.TargetUUID,
			},
		}
	case cmd.DeleteChunkCmd != nil:
		pbCmd.Command = &pb.Command_DeleteChunkCmd{
			DeleteChunkCmd: &pb.DeleteChunkCmd{
				ChunkId: cmd.DeleteChunkCmd.ChunkID,
			},
		}
	}

	return pbCmd
}

// Proto2HeartbeatRequest converts a protobuf HeartbeatRequest to internal format
func Proto2HeartbeatRequest(pbReq *pb.HeartbeatRequest) *HeartbeatRequest {
	results := make([]CommandResult, len(pbReq.CommandResults))
	for i, pbRes := range pbReq.CommandResults {
		results[i] = CommandResult{
			CommandID:    pbRes.CommandId,
			Success:      pbRes.Success,
			ErrorMessage: pbRes.ErrorMessage,
		}
	}
	return &HeartbeatRequest{
		UUID:           pbReq.Uuid,
		CommandResults: results,
	}
}

// HeartbeatRequest2Proto converts a HeartbeatRequest to protobuf format
func HeartbeatRequest2Proto(req *HeartbeatRequest) *pb.HeartbeatRequest {
	results := make([]*pb.CommandResult, len(req.CommandResults))
	for i, res := range req.CommandResults {
		results[i] = &pb.CommandResult{
			CommandId:    res.CommandID,
			Success:      res.Success,
			ErrorMessage: res.ErrorMessage,
		}
	}
	return &pb.HeartbeatRequest{
		Uuid:           req.UUID,
		CommandResults: results,
	}
}

// Proto2HeartbeatResponse converts a protobuf HeartbeatResponse to internal format
func Proto2HeartbeatResponse(pbResp *pb.HeartbeatResponse) *HeartbeatResponse {
	commands := make([]*Command, len(pbResp.Commands))
	for i, pbCmd := range pbResp.Commands {
		commands[i] = Proto2Command(pbCmd)
	}
	return &HeartbeatResponse{
		Commands: commands,
	}
}

// HeartbeatResponse2Proto converts a HeartbeatResponse to protobuf format
func HeartbeatResponse2Proto(resp *HeartbeatResponse) *pb.HeartbeatResponse {
	pbCommands := make([]*pb.Command, len(resp.Commands))
	for i, cmd := range resp.Commands {
		pbCommands[i] = Command2Proto(cmd)
	}
	return &pb.HeartbeatResponse{
		Commands: pbCommands,
	}
}

func Proto2PutSDPCandidatesRequest(pbReq *pb.PutSDPCandidatesRequest) *PutSDPCandidatesRequest {
	return &PutSDPCandidatesRequest{
		SourceUUID: pbReq.SourceUuid,
		TargetUUID: pbReq.TargetUuid,
		SDP:        pbReq.Sdp,
		Candidates: pbReq.Candidates,
	}
}

func PutSDPCandidatesRequest2Proto(req *PutSDPCandidatesRequest) *pb.PutSDPCandidatesRequest {
	return &pb.PutSDPCandidatesRequest{
		SourceUuid: req.SourceUUID,
		TargetUuid: req.TargetUUID,
		Sdp:        req.SDP,
		Candidates: req.Candidates,
	}
}

func Proto2PutSDPCandidatesResponse(pbResp *pb.PutSDPCandidatesResponse) *PutSDPCandidatesResponse {
	return &PutSDPCandidatesResponse{
		Success:      pbResp.Success,
		ErrorMessage: pbResp.ErrorMessage,
	}
}

func PutSDPCandidatesResponse2Proto(resp *PutSDPCandidatesResponse) *pb.PutSDPCandidatesResponse {
	return &pb.PutSDPCandidatesResponse{
		Success:      resp.Success,
		ErrorMessage: resp.ErrorMessage,
	}
}

func Proto2GetSDPCandidatesRequest(pbReq *pb.GetSDPCandidatesRequest) *GetSDPCandidatesRequest {
	return &GetSDPCandidatesRequest{
		SourceUUID: pbReq.SourceUuid,
		TargetUUID: pbReq.TargetUuid,
	}
}

func GetSDPCandidatesRequest2Proto(req *GetSDPCandidatesRequest) *pb.GetSDPCandidatesRequest {
	return &pb.GetSDPCandidatesRequest{
		SourceUuid: req.SourceUUID,
		TargetUuid: req.TargetUUID,
	}
}
func Proto2GetSDPCandidatesResponse(pbResp *pb.GetSDPCandidatesResponse) *GetSDPCandidatesResponse {
	return &GetSDPCandidatesResponse{
		Success:      pbResp.Success,
		ErrorMessage: pbResp.ErrorMessage,
		SourceUUID:   pbResp.SourceUuid,
		TargetUUID:   pbResp.TargetUuid,
		SDP:          pbResp.Sdp,
		Candidates:   pbResp.Candidates,
	}
}
func GetSDPCandidatesResponse2Proto(resp *GetSDPCandidatesResponse) *pb.GetSDPCandidatesResponse {
	return &pb.GetSDPCandidatesResponse{
		Success:      resp.Success,
		ErrorMessage: resp.ErrorMessage,
		SourceUuid:   resp.SourceUUID,
		TargetUuid:   resp.TargetUUID,
		Sdp:          resp.SDP,
		Candidates:   resp.Candidates,
	}
}
