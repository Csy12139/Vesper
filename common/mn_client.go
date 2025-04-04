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


func (c *MNClient) AddChunkMeta() (*AddChunkMetaResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	pbReq := AddChunkMetaRequest2Proto(&AddChunkMetaRequest{})
	pbResp, err := c.client.AddChunkMeta(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("add chunk meta RPC failed: %w", err)
	}
	return Proto2AddChunkMetaResponse(pbResp), nil
}

func (c *MNClient) CompleteAddChunkMeta(chunkId uint64) (*CompleteAddChunkMetaResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &CompleteAddChunkMetaRequest{
		ChunkID: chunkId,
	}
	pbReq := CompleteAddChunkMetaRequest2Proto(req)
	pbResp, err := c.client.CompleteAddChunkMeta(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("complete add chunk meta RPC failed: %w", err)
	}
	return Proto2CompleteAddChunkMetaResponse(pbResp), nil
}

func (c *MNClient) GetChunkMeta(chunkId uint64) (*GetChunkMetaResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &GetChunkMetaRequest{
		ChunkID: chunkId,
	}
	pbReq := GetChunkMetaRequest2Proto(req)
	pbResp, err := c.client.GetChunkMeta(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("get chunk meta RPC failed: %w", err)
	}
	return Proto2GetChunkMetaResponse(pbResp), nil
}

func (c *MNClient) AllocateDnForChunk(chunkId uint64, excludes []string) (*AllocateDnForChunkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &AllocateDnForChunkRequest{
		ChunkId:  chunkId,
		Excludes: excludes,
	}
	pbReq := AllocateDnForChunkRequest2Proto(req)
	pbResp, err := c.client.AllocateDnForChunk(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("allocate dn for chunk RPC failed: %w", err)
	}
	return Proto2AllocateDnForChunkResponse(pbResp), nil
}

func (c *MNClient) AddChunkOnDN(sdkUuid string, dnUuid string) (*AddChunkOnDNResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &AddChunkOnDNRequest{
		SdkUuid: sdkUuid,
		DnUuid:  dnUuid,
	}
	pbReq := AddChunkOnDNRequest2Proto(req)
	pbResp, err := c.client.AddChunkOnDN(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("add chunk on dn RPC failed: %w", err)
	}
	return Proto2AddChunkOnDNResponse(pbResp), nil
}

func (c *MNClient) CompleteAddChunkOnDN(chunkId uint64, uuid string) (*CompleteAddChunkOnDNResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &CompleteAddChunkOnDNRequest{
		ChunkId: chunkId,
		Uuid:    uuid,
	}
	pbReq := CompleteAddChunkOnDNRequest2Proto(req)
	pbResp, err := c.client.CompleteAddChunkOnDN(ctx, pbReq)
	if err != nil {
		return nil, fmt.Errorf("complete add chunk on dn RPC failed: %w", err)
	}
	return Proto2CompleteAddChunkOnDNResponse(pbResp), nil
}
