package main

import (
	"context"
	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type mnServiceServer struct {
	//server                          pb.MNServiceServer
	pb.UnimplementedMNServiceServer // Embed the unimplemented interface to ensure compatibility.
}

func StartMNServer(network, address string) {
	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMNServiceServer(s, &mnServiceServer{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (m *mnServiceServer) PutSDPCandidates(ctx context.Context, req *pb.PutSDPCandidatesRequest) (*pb.PutSDPCandidatesResponse, error) {
	return &pb.PutSDPCandidatesResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (m *mnServiceServer) GetSDPCandidates(ctx context.Context, req *pb.GetSDPCandidatesRequest) (*pb.GetSDPCandidatesResponse, error) {
	return &pb.GetSDPCandidatesResponse{
		SourceUuid: req.SourceUuid,
		TargetUuid: req.TargetUuid,
		Sdp:        "xxx",
		Candidates: nil,
	}, nil
}
