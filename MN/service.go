package MN

import (
	"context"
	"fmt"
	pb "github.com/Csy12139/Vesper/proto"
)

func (mn *MNServer) PutSDPCandidates(ctx context.Context, req *pb.PutSDPCandidatesRequest) (*pb.PutSDPCandidatesResponse, error) {
	key := req.SourceUuid + "|" + req.TargetUuid
	mn.mu.Lock()
	mn.SDPCandidatesMap[key] = req
	mn.mu.Unlock()
	return &pb.PutSDPCandidatesResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (mn *MNServer) GetSDPCandidates(ctx context.Context, req *pb.GetSDPCandidatesRequest) (*pb.GetSDPCandidatesResponse, error) {
	key := req.SourceUuid + "|" + req.TargetUuid

	mn.mu.RLock()
	SDPCandidates, ok := mn.SDPCandidatesMap[key]
	mn.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("SDPCandidates not found")
	}
	return &pb.GetSDPCandidatesResponse{
		Success:      true,
		ErrorMessage: "",
		SourceUuid:   req.SourceUuid,
		TargetUuid:   req.TargetUuid,
		Sdp:          SDPCandidates.Sdp,
		Candidates:   SDPCandidates.Candidates,
	}, nil
}

func (mn *MNServer) DoHeartbeat(context.Context, *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {

}
