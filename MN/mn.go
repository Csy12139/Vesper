package MN

import (
	"fmt"
	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type MateNode struct {
	MNNetwork string
	MNAddr    string
	mu        sync.RWMutex
	// TODO change to map[string][string]
	// TODO add timestamp
	SDPCandidatesMap map[string]*pb.PutSDPCandidatesRequest
	grpcServer       *grpc.Server
	pb.UnimplementedMNServiceServer
	stop         atomic.Bool
	dataNodes    map[string]*DataNodeInfo
	dataNodeLock sync.RWMutex
}

func NewMNServer(MNAddr string) (*MateNode, error) {
	mn := &MateNode{
		MNNetwork:        "tcp",
		MNAddr:           MNAddr,
		SDPCandidatesMap: make(map[string]*pb.PutSDPCandidatesRequest),
	}
	mn.stop.Store(true)
	return mn, nil
}

func (mn *MateNode) StartMetaNode() {
	mn.stop.Store(false)
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
		mn.stop.Store(true)
	}()
	// 启动后台任务
	go mn.monitorQueueTimeoutCmd()
	go mn.monitorWaitResponseTimeoutCmd()

}

func (mn *MateNode) IsRunning() bool {
	return !mn.stop.Load()
}

func (mn *MateNode) StopMetaNode() {
	mn.stop.Store(true)
	mn.grpcServer.Stop()
}

func (mn *MateNode) monitorQueueTimeoutCmd() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if mn.stop.Load() {
			return
		}
		for _, dn := range mn.dataNodes {
			select {
			case cmd := <-dn.cmdQueue:
				if time.Since(cmd.StartTimestamp) > cmd.Timeout {
					go cmd.CallBack(cmd.ID, fmt.Errorf("cmd queue timeout!start time:%v,timeout:%v", cmd.StartTimestamp, cmd.Timeout))
				} else {
					dn.cmdQueue <- cmd
				}
			default:
				break
			}
		}

	}
}
func (mn *MateNode) monitorWaitResponseTimeoutCmd() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if mn.stop.Load() {
			return
		}
		for _, dn := range mn.dataNodes {
			dn.mutex.Lock()
			for taskId, cmd := range dn.waitResponseCommand {
				if time.Since(cmd.StartTimestamp) > cmd.Timeout {
					go cmd.CallBack(cmd.ID, fmt.Errorf("cmd wait response timeout!start time:%v,timeout:%v", cmd.StartTimestamp, cmd.Timeout))
					delete(dn.waitResponseCommand, taskId) // go 可以一边遍历一边删除
				}
			}
			dn.mutex.Unlock()
		}
	}
}
