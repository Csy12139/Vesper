package common

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const mnClientTimeout = 10 * time.Second

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

func (c *MNClient) PutSDPCandidates(SourceUUID string, TargetUUID string, SDP string, Candidates []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()
	req := &pb.PutSDPCandidatesRequest{
		SourceUuid: SourceUUID,
		TargetUuid: TargetUUID,
		Sdp:        SDP,
		Candidates: Candidates,
	}
	_, err := c.client.PutSDPCandidates(ctx, req)
	return err
}

func (c *MNClient) GetSDPCandidates(SourceUUID string, TargetUUID string) (SDP string, Candidates []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()
	req := &pb.GetSDPCandidatesRequest{
		SourceUuid: SourceUUID,
		TargetUuid: TargetUUID,
	}
	pbResp, err := c.client.GetSDPCandidates(ctx, req)
	if err != nil {
		return "", nil, fmt.Errorf("[Technical Error]get sdp and candidates failed [%w]", err)
	}

	err = Proto2Error(pbResp.Code)
	if err != nil {
		return "", nil, err
	}
	return pbResp.Sdp, pbResp.Candidates, nil
}

func (c *MNClient) AddChunkMeta() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	pbResp, err := c.client.AddChunkMeta(ctx, &pb.AddChunkMetaRequest{})
	if err != nil {
		return 0, fmt.Errorf("add chunk meta failed: %w", err)
	}
	return pbResp.ChunkId, nil
}

func (c *MNClient) CompleteAddChunkMeta(chunkId uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &pb.CompleteAddChunkMetaRequest{
		ChunkId: chunkId,
	}
	_, err := c.client.CompleteAddChunkMeta(ctx, req)
	if err != nil {
		return fmt.Errorf("complete add chunk meta failed: %w", err)
	}
	return nil
}

func (c *MNClient) GetChunkMeta(chunkId uint64) (*ChunkMeta, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &pb.GetChunkMetaRequest{
		ChunkId: chunkId,
	}
	pbResp, err := c.client.GetChunkMeta(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get chunk meta failed: %w", err)
	}
	err = Proto2Error(pbResp.Code)
	if err != nil {
		return nil, err
	}

	return Proto2ChunkMeta(pbResp.Meta), nil
}

func (c *MNClient) AllocateDnForChunk(chunkId uint64, excludes []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &pb.AllocateDnForChunkRequest{
		ChunkId:  chunkId,
		Excludes: excludes,
	}
	pbResp, err := c.client.AllocateDnForChunk(ctx, req)
	if err != nil {
		return "", fmt.Errorf("allocate dn for chunk failed: %w", err)
	}
	err = Proto2Error(pbResp.Code)
	if err != nil {
		return "", err
	}
	return pbResp.Uuid, nil
}

func (c *MNClient) AddChunkOnDN(sdkUuid string, dnUuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &pb.AddChunkOnDNRequest{
		SdkUuid: sdkUuid,
		DnUuid:  dnUuid,
	}
	pbResp, err := c.client.AddChunkOnDN(ctx, req)
	if err != nil {
		return fmt.Errorf("[Technical Error] add chunk on dn %v failed: %w", req.DnUuid, err)
	}
	err = Proto2Error(pbResp.Code)
	if err != nil {
		return fmt.Errorf("[Business Error] add chunk on dn failed: %w", err)
	}
	return nil
}

func (c *MNClient) CompleteAddChunkOnDN(chunkId uint64, uuid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), mnClientTimeout)
	defer cancel()

	req := &pb.CompleteAddChunkOnDNRequest{
		ChunkId: chunkId,
		Uuid:    uuid,
	}
	pbResp, err := c.client.CompleteAddChunkOnDN(ctx, req)
	if err != nil {
		return fmt.Errorf("complete add chunk on dn failed: %w", err)
	}
	err = Proto2Error(pbResp.Code)
	if err != nil {
		return err
	}
	return nil
}
