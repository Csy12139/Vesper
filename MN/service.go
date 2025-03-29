package MN

import (
	"context"
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	pb "github.com/Csy12139/Vesper/proto"
	"time"
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
		request.UUID, time.Now().Sub(dn.lastHeartbeatTime).Seconds(), request.CommandResults)
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
