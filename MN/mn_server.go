package main

import (
	"context"
	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"sync"
)

type mnServiceServer struct {
	//server                          pb.MNServiceServer
	mu                              sync.RWMutex
	SDPCandidatesMap                map[string][]byte
	pb.UnimplementedMNServiceServer // Embed the unimplemented interface to ensure compatibility.
}

func StartMNServer(network, address string) {
	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMNServiceServer(s, &mnServiceServer{
		SDPCandidatesMap: make(map[string][]byte),
	})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (m *mnServiceServer) PutSDPCandidates(ctx context.Context, req *pb.PutSDPCandidatesRequest) (*pb.PutSDPCandidatesResponse, error) {
	key := req.SourceUuid + "|" + req.TargetUuid

	value, err := proto.Marshal(&pb.PutSDPCandidatesRequest{
		Sdp:        req.Sdp,
		Candidates: req.Candidates,
	})
	if err != nil {
		return &pb.PutSDPCandidatesResponse{
			Success:      false,
			ErrorMessage: "failed to marshal SDP and Candidates data: " + err.Error(),
		}, nil
	}
	m.mu.Lock()
	m.SDPCandidatesMap[key] = value
	m.mu.Unlock()
	return &pb.PutSDPCandidatesResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (m *mnServiceServer) GetSDPCandidates(ctx context.Context, req *pb.GetSDPCandidatesRequest) (*pb.GetSDPCandidatesResponse, error) {
	key := req.SourceUuid + "|" + req.TargetUuid

	m.mu.RLock()
	value, ok := m.SDPCandidatesMap[key]
	m.mu.RUnlock()

	if !ok {
		return &pb.GetSDPCandidatesResponse{
			Success:      false,
			ErrorMessage: "failed to get SDP and Candidates data: ",
		}, nil
	}
	SDPCandidates := &pb.PutSDPCandidatesRequest{}
	err := proto.Unmarshal(value, SDPCandidates)
	if err != nil {
		return &pb.GetSDPCandidatesResponse{
			Success:      false,
			ErrorMessage: "failed to unmarshal SDP and Candidates data: " + err.Error(),
		}, nil
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
