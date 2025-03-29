package MN

import (
	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type MNServer struct {
	MNNetwork string
	MNAddr    string
	mu        sync.RWMutex
	// TODO change to map[string][string]
	// TODO add timestamp
	SDPCandidatesMap map[string]*pb.PutSDPCandidatesRequest
	grpcServer       *grpc.Server
}

func NewMNServer(MNAddr string) (*MNServer, error) {
	return &MNServer{
		MNNetwork:        "tcp",
		MNAddr:           MNAddr,
		SDPCandidatesMap: make(map[string][]*pb.PutSDPCandidatesRequest),
	}, nil
}

func (mn *MNServer) StartMNServer() {
	go func() {
		lis, err := net.Listen(mn.MNNetwork, mn.MNAddr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		mn.grpcServer = grpc.NewServer()
		pb.RegisterMNServiceServer(mn.grpcServer, mn)
		log.Printf("Server listening at %v", lis.Addr())
		if err := mn.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (mn *MNServer) StopMNServer() {
	mn.grpcServer.Stop()
}
