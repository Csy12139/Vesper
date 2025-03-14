package main

import (
	"context"
	pb "github.com/Csy12139/Vesper/grpcutil/proto"
	"github.com/Csy12139/Vesper/log"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	pb.UnimplementedMNServer
}

func (s *server) CreateBucketRequest(ctx context.Context, req *pb.CreateBucket) (*pb.BoolResponse, error) {
	return &pb.BoolResponse{Success: true}, nil
}

func (s *server) UploadRequest(ctx context.Context, req *pb.UploadData) (*pb.DNReplicaID, error) {
	return &pb.DNReplicaID{}, nil
}

func (s *server) DownloadRequest(ctx context.Context, req *pb.DownloadData) (*pb.DNReplicaID, error) {
	return &pb.DNReplicaID{}, nil
}
func (s *server) HeartDetectRequest(ctx context.Context, req *pb.HeartbeatDetection) (*pb.HeartbeatMonitor, error) {
	return &pb.HeartbeatMonitor{Success: true}, nil
}

func setupListener(addr string) net.Listener {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return lis
}
func createServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterMNServer(s, &server{})
	return s
}
func startServer(s *grpc.Server, lis net.Listener) {
	log.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func initServer(addr string) {
	lis := setupListener(addr)
	s := createServer()
	startServer(s, lis)
}
func main() {
	err := log.InitLog("./logs", 10, 5, "info")
	if err != nil {
		log.Fatalf("log init failed: %v", err)
	}
	initServer(":50051")
}
