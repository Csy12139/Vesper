package MN

import (
	"context"
	"fmt"
	"time"

	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	pb "github.com/Csy12139/Vesper/proto"
	"github.com/pkg/errors"
)

func (mn *MateNode) PutSDPCandidates(ctx context.Context, req *pb.PutSDPCandidatesRequest) (*pb.PutSDPCandidatesResponse, error) {
	key := req.SourceUuid + "|" + req.TargetUuid
	mn.mu.Lock()
	mn.SDPCandidatesMap[key] = req
	mn.mu.Unlock()
	return &pb.PutSDPCandidatesResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

func (mn *MateNode) GetSDPCandidates(ctx context.Context, req *pb.GetSDPCandidatesRequest) (*pb.GetSDPCandidatesResponse, error) {
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

func (mn *MateNode) DoHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	request := common.Proto2HeartbeatRequest(req)

	mn.dataNodeLock.RLock()
	dn, exist := mn.dataNodes[request.UUID]
	mn.dataNodeLock.RUnlock()
	if !exist {
		log.Infof("receive first heartbeat from data node [%s], now join it", request.UUID)
		dn = NewDataNodeInfo()
		dn.SetALIVE()
		mn.dataNodeLock.Lock()
		defer mn.dataNodeLock.Unlock()
		mn.dataNodes[request.UUID] = dn
		return common.HeartbeatResponse2Proto(&common.HeartbeatResponse{Commands: make([]*common.Command, 0)}), nil
	}
	dn.mutex.Lock()
	defer dn.mutex.Unlock()
	log.Debugf("receive heartbeat from data node [%s] interval [%v]s cmd response [%+v]",
		request.UUID, time.Since(dn.lastHeartbeatTime).Seconds(), request.CommandResults)
	dn.SetALIVE()
	dn.lastHeartbeatTime = time.Now()
	for _, cmdResult := range request.CommandResults {
		cmd, exist := dn.waitResponseCommand[cmdResult.CommandID]
		if !exist {
			log.Errorf("receive heartbeat from data node [%s] but not found command id [%d]",
				request.UUID, cmdResult.CommandID)
			continue
		}
		if cmdResult.Success {
			go cmd.CallBack(cmd.ID, nil)
		} else {
			go cmd.CallBack(cmd.ID, fmt.Errorf(cmdResult.ErrorMessage))
		}
	}
	response := &common.HeartbeatResponse{
		Commands: make([]*common.Command, 0),
	}
	// TODO:limit cmd number
	for {
		select {
		case cmd := <-dn.cmdQueue:
			response.Commands = append(response.Commands, cmd)
			dn.waitResponseCommand[cmd.ID] = cmd
		default:
			return common.HeartbeatResponse2Proto(response), nil
		}
	}
}

func (mn *MateNode) AddChunkMeta(ctx context.Context, req *pb.AddChunkMetaRequest) (*pb.AddChunkMetaResponse, error) {
	id, err := mn.kv.AllocateChunkId()
	mn.lockChunk(id)
	defer mn.unlockChunk(id)
	if err != nil {
		return nil, err
	}
	chunk := &common.ChunkMeta{
		ID:      id,
		State:   common.ChunkState_CREATING,
		DnUuids: make(map[string]struct{}),
	}
	err = mn.kv.PutChunkMeta(chunk)
	if err != nil {
		return nil, err
	}
	return &pb.AddChunkMetaResponse{
		ChunkId: chunk.ID,
	}, nil
}

func (mn *MateNode) GetChunkMeta(ctx context.Context, req *pb.GetChunkMetaRequest) (*pb.GetChunkMetaResponse, error) {
	meta, err := mn.kv.GetChunkMeta(req.ChunkId)
	if err != nil {
		if errors.Is(err, common.ErrChunkNotFound) {
			return &pb.GetChunkMetaResponse{
				Code: pb.ErrorCode_CHUNK_NOT_FOUND,
			}, nil
		}
		return nil, err
	}
	return &pb.GetChunkMetaResponse{
		Meta: common.ChunkMeta2Proto(meta),
	}, nil
}

func (mn *MateNode) CompleteAddChunkMeta(ctx context.Context, req *pb.CompleteAddChunkMetaRequest) (*pb.CompleteAddChunkMetaResponse, error) {
	chunkId := req.ChunkId
	mn.lockChunk(chunkId)
	defer mn.unlockChunk(chunkId)
	meta, err := mn.kv.GetChunkMeta(req.ChunkId)
	if err != nil {
		return nil, err
	}
	meta.State = common.ChunkState_CREATED
	err = mn.kv.PutChunkMeta(meta)
	if err != nil {
		return nil, err
	}
	return &pb.CompleteAddChunkMetaResponse{}, nil
}

func (mn *MateNode) AllocateDnForChunk(ctx context.Context, req *pb.AllocateDnForChunkRequest) (*pb.AllocateDnForChunkResponse, error) {
	excludes := req.Excludes
	dnUUID, err := mn.allocateDn(excludes)
	if err != nil {
		if errors.Is(err, common.ErrNoAvailableDN) {
			return &pb.AllocateDnForChunkResponse{
				Uuid: "",
				Code: pb.ErrorCode_NO_AVAILABLE_DN,
			}, nil
		}
		return nil, err
	}
	return &pb.AllocateDnForChunkResponse{
		Uuid: dnUUID,
		Code: pb.ErrorCode_OK,
	}, nil
}

func (mn *MateNode) AddChunkOnDN(ctx context.Context, req *pb.AddChunkOnDNRequest) (*pb.AddChunkOnDNResponse, error) {
	mn.dataNodeLock.RLock()
	dn, exist := mn.dataNodes[req.DnUuid]
	mn.dataNodeLock.RUnlock()
	if !exist {
		return &pb.AddChunkOnDNResponse{
			Code: pb.ErrorCode_DN_NOT_FOUND,
		}, nil
	}
	dn.mutex.Lock()
	defer dn.mutex.Unlock()

	done := make(chan error, 1)
	dn.SubmitAddChunkCmd(req.SdkUuid, 5*time.Second, func(cmdId uint64, err error) {
		done <- err
	})

	if err := <-done; err != nil {
		log.Errorf("Failed to add chunk on DN %s: %v", req.DnUuid, err)
		return &pb.AddChunkOnDNResponse{
			Code: pb.ErrorCode_COMMIT_DN_CMD_TIMEOUT,
		}, nil
	}
	return &pb.AddChunkOnDNResponse{
		Code: pb.ErrorCode_OK,
	}, nil
}

func (mn *MateNode) CompleteAddChunkOnDN(ctx context.Context, req *pb.CompleteAddChunkOnDNRequest) (*pb.CompleteAddChunkOnDNResponse, error) {
	chunkID := req.ChunkId
	dnUUID := req.Uuid

	mn.lockChunk(chunkID)
	defer mn.unlockChunk(chunkID)

	chunkMeta, err := mn.kv.GetChunkMeta(chunkID)
	if err != nil {
		if errors.Is(err, common.ErrChunkNotFound) {
			return &pb.CompleteAddChunkOnDNResponse{
				Code: pb.ErrorCode_CHUNK_NOT_FOUND,
			}, nil
		}
		return nil, err
	}
	chunkMeta.DnUuids[dnUUID] = struct{}{}
	err = mn.kv.PutChunkMeta(chunkMeta)
	if err != nil {
		return nil, err
	}
	return &pb.CompleteAddChunkOnDNResponse{}, nil
}
