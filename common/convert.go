package common

import (
	"errors"
	pb "github.com/Csy12139/Vesper/proto"
	"fmt"
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

func ChunkMeta2Proto(meta *ChunkMeta) *pb.ChunkMeta {
	if meta == nil {
		return nil
	}
	protoMeta := &pb.ChunkMeta{
		Id:    meta.ID,
		State: pb.ChunkState(meta.State),
	}
	// Convert map[string]struct{} to repeated string
	if meta.DnUuids != nil {
		protoMeta.DnUuids = make([]string, 0, len(meta.DnUuids))
		for uuid := range meta.DnUuids {
			protoMeta.DnUuids = append(protoMeta.DnUuids, uuid)
		}
	}

	return protoMeta
}

func Proto2ChunkMeta(protoMeta *pb.ChunkMeta) *ChunkMeta {
	if protoMeta == nil {
		return nil
	}
	meta := &ChunkMeta{
		ID:      protoMeta.Id,
		State:   ChunkState(protoMeta.State),
		DnUuids: make(map[string]struct{}),
	}

	// Convert repeated string to map[string]struct{}
	for _, uuid := range protoMeta.DnUuids {
		meta.DnUuids[uuid] = struct{}{}
	}
	return meta
}



func Error2Proto(err error) pb.ErrorCode {
	if err == nil {
		return pb.ErrorCode_OK
	}
	switch {
	case errors.Is(err, ErrChunkNotFound):
		return pb.ErrorCode_CHUNK_NOT_FOUND
	case errors.Is(err, ErrNoAvailableDN):
		return pb.ErrorCode_NO_AVAILABLE_DN
	case errors.Is(err, ErrCommitDNTimeout):
		return pb.ErrorCode_COMMIT_DN_CMD_TIMEOUT
	case errors.Is(err, ErrDNNotFound):
		return pb.ErrorCode_DN_NOT_FOUND
	default:
		return pb.ErrorCode_UNKNOWN
	}
}

func Proto2Error(code pb.ErrorCode) error {
	switch code {
		case pb.ErrorCode_OK:
		return nil
	case pb.ErrorCode_CHUNK_NOT_FOUND:
		return ErrChunkNotFound
	case pb.ErrorCode_NO_AVAILABLE_DN:
		return ErrNoAvailableDN
	case pb.ErrorCode_COMMIT_DN_CMD_TIMEOUT:
		return ErrCommitDNTimeout
	case pb.ErrorCode_DN_NOT_FOUND:
		return ErrDNNotFound
	default:
		return fmt.Errorf("unknown error code: %v", code)
	}
}

// AddChunkMeta conversions
func Proto2AddChunkMetaRequest(pbReq *pb.AddChunkMetaRequest) *AddChunkMetaRequest {
    return &AddChunkMetaRequest{}
}

func AddChunkMetaRequest2Proto(req *AddChunkMetaRequest) *pb.AddChunkMetaRequest {
    return &pb.AddChunkMetaRequest{}
}

func Proto2AddChunkMetaResponse(pbResp *pb.AddChunkMetaResponse) *AddChunkMetaResponse {
    return &AddChunkMetaResponse{
        ChunkId: pbResp.ChunkId,
    }
}

func AddChunkMetaResponse2Proto(resp *AddChunkMetaResponse) *pb.AddChunkMetaResponse {
    return &pb.AddChunkMetaResponse{
        ChunkId: resp.ChunkId,
    }
}

// CompleteAddChunkMeta conversions
func Proto2CompleteAddChunkMetaRequest(pbReq *pb.CompleteAddChunkMetaRequest) *CompleteAddChunkMetaRequest {
    return &CompleteAddChunkMetaRequest{
        ChunkID: pbReq.ChunkId,
    }
}

func CompleteAddChunkMetaRequest2Proto(req *CompleteAddChunkMetaRequest) *pb.CompleteAddChunkMetaRequest {
    return &pb.CompleteAddChunkMetaRequest{
        ChunkId: req.ChunkID,
    }
}

func Proto2CompleteAddChunkMetaResponse(pbResp *pb.CompleteAddChunkMetaResponse) *CompleteAddChunkMetaResponse {
    return &CompleteAddChunkMetaResponse{}
}

func CompleteAddChunkMetaResponse2Proto(resp *CompleteAddChunkMetaResponse) *pb.CompleteAddChunkMetaResponse {
    return &pb.CompleteAddChunkMetaResponse{}
}

// GetChunkMeta conversions
func Proto2GetChunkMetaRequest(pbReq *pb.GetChunkMetaRequest) *GetChunkMetaRequest {
    return &GetChunkMetaRequest{
        ChunkID: pbReq.ChunkId,
    }
}

func GetChunkMetaRequest2Proto(req *GetChunkMetaRequest) *pb.GetChunkMetaRequest {
    return &pb.GetChunkMetaRequest{
        ChunkId: req.ChunkID,
    }
}

func Proto2GetChunkMetaResponse(pbResp *pb.GetChunkMetaResponse) *GetChunkMetaResponse {
    return &GetChunkMetaResponse{
        Meta: Proto2ChunkMeta(pbResp.Meta),
        Code: Proto2Error(pbResp.Code),
    }
}

func GetChunkMetaResponse2Proto(resp *GetChunkMetaResponse) *pb.GetChunkMetaResponse {
    return &pb.GetChunkMetaResponse{
        Meta: ChunkMeta2Proto(resp.Meta),
        Code: Error2Proto(resp.Code),
    }
}

// AllocateDnForChunk conversions
func Proto2AllocateDnForChunkRequest(pbReq *pb.AllocateDnForChunkRequest) *AllocateDnForChunkRequest {
    return &AllocateDnForChunkRequest{
        ChunkId: pbReq.ChunkId,
        Excludes: pbReq.Excludes,
    }
}

func AllocateDnForChunkRequest2Proto(req *AllocateDnForChunkRequest) *pb.AllocateDnForChunkRequest {
    return &pb.AllocateDnForChunkRequest{
        ChunkId: req.ChunkId,
        Excludes: req.Excludes,
    }
}

func Proto2AllocateDnForChunkResponse(pbResp *pb.AllocateDnForChunkResponse) *AllocateDnForChunkResponse {
    return &AllocateDnForChunkResponse{
        Uuid: pbResp.Uuid,
        Code: Proto2Error(pbResp.Code),
    }
}

func AllocateDnForChunkResponse2Proto(resp *AllocateDnForChunkResponse) *pb.AllocateDnForChunkResponse {
    return &pb.AllocateDnForChunkResponse{
        Uuid: resp.Uuid,
        Code: Error2Proto(resp.Code),
    }
}

// AddChunkOnDN conversions
func Proto2AddChunkOnDNRequest(pbReq *pb.AddChunkOnDNRequest) *AddChunkOnDNRequest {
    return &AddChunkOnDNRequest{
        SdkUuid: pbReq.SdkUuid,
        DnUuid:  pbReq.DnUuid,
    }
}

func AddChunkOnDNRequest2Proto(req *AddChunkOnDNRequest) *pb.AddChunkOnDNRequest {
    return &pb.AddChunkOnDNRequest{
        SdkUuid: req.SdkUuid,
        DnUuid:  req.DnUuid,
    }
}

func Proto2AddChunkOnDNResponse(pbResp *pb.AddChunkOnDNResponse) *AddChunkOnDNResponse {
    return &AddChunkOnDNResponse{
        Code: Proto2Error(pbResp.Code),
    }
}

func AddChunkOnDNResponse2Proto(resp *AddChunkOnDNResponse) *pb.AddChunkOnDNResponse {
    return &pb.AddChunkOnDNResponse{
        Code: Error2Proto(resp.Code),
    }
}

// CompleteAddChunkOnDN conversions
func Proto2CompleteAddChunkOnDNRequest(pbReq *pb.CompleteAddChunkOnDNRequest) *CompleteAddChunkOnDNRequest {
    return &CompleteAddChunkOnDNRequest{
        ChunkId: pbReq.ChunkId,
        Uuid:    pbReq.Uuid,
    }
}

func CompleteAddChunkOnDNRequest2Proto(req *CompleteAddChunkOnDNRequest) *pb.CompleteAddChunkOnDNRequest {
    return &pb.CompleteAddChunkOnDNRequest{
        ChunkId: req.ChunkId,
        Uuid:    req.Uuid,
    }
}

func Proto2CompleteAddChunkOnDNResponse(pbResp *pb.CompleteAddChunkOnDNResponse) *CompleteAddChunkOnDNResponse {
    return &CompleteAddChunkOnDNResponse{
        Code: Proto2Error(pbResp.Code),
    }
}

func CompleteAddChunkOnDNResponse2Proto(resp *CompleteAddChunkOnDNResponse) *pb.CompleteAddChunkOnDNResponse {
    return &pb.CompleteAddChunkOnDNResponse{
        Code: Error2Proto(resp.Code),
    }
}






