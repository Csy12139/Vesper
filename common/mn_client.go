package common

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const mnClientTimeout = 5 * time.Second

type MNClient struct {
	addr   string
	client pb.MNServiceClient
}

func NewMNClient(addr string) (*MNClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}
	return &MNClient{
		addr:   addr,
		client: pb.NewMNServiceClient(conn),
	}, nil
}

func (c *MNClient) DoHeartbeat(uuid string, results []CommandResult) (*HeartbeatResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &HeartbeatRequest{
		UUID:           uuid,
		CommandResults: results,
	}

	pbReq := HeartbeatRequest2Proto(req)
	pbResp, err := c.client.DoHeartbeat(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("heartbeat RPC failed: %w", err)
	}

	return Proto2HeartbeatResponse(pbResp), nil
}

func (c *MNClient) PutSDPCandidates(req *PutSDPCandidatesRequest) (*PutSDPCandidatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	pbReq := PutSDPCandidatesRequest2Proto(req)
	pbResp, err := c.client.PutSDPCandidates(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("put sdp and candidates RPC failed: %w", err)
	}
	return Proto2PutSDPCandidatesResponse(pbResp), nil
}

func (c *MNClient) GetSDPCandidates(req *GetSDPCandidatesRequest) (*GetSDPCandidatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	pbReq := GetSDPCandidatesRequest2Proto(req)
	pbResp, err := c.client.GetSDPCandidates(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("get sdp and candidates RPC failed: %w", err)
	}
	return Proto2GetSDPCandidatesResponse(pbResp), nil
}
