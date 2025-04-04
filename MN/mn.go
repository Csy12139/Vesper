package MN

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Csy12139/Vesper/common"
	pb "github.com/Csy12139/Vesper/proto"
	"google.golang.org/grpc"
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
	stop           atomic.Bool
	dataNodes      map[string]*DataNodeInfo
	dataNodeLock   sync.RWMutex
	chunkMetaMutex sync.Map
	kv             *KVTable
}

func NewMNServer(MNAddr string, dataPath string) (*MateNode, error) {
	mn := &MateNode{
		MNNetwork:        "tcp",
		MNAddr:           MNAddr,
		SDPCandidatesMap: make(map[string]*pb.PutSDPCandidatesRequest),
		kv:               NewKVTable(),
		dataNodes:        make(map[string]*DataNodeInfo),
		dataNodeLock:     sync.RWMutex{},
		chunkMetaMutex:   sync.Map{},
	}
	mn.stop.Store(true)
	err := mn.kv.Open(dataPath)
	if err != nil {
		return nil, err
	}
	return mn, nil
}

func (mn *MateNode) Start() {
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

func (mn *MateNode) Stop() {
	mn.stop.Store(true)
	mn.grpcServer.Stop()
}

func (mn *MateNode) monitorQueueTimeoutCmd() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if mn.stop.Load() {
			return
		}
		mn.dataNodeLock.RLock()
		for _, dn := range mn.dataNodes {
			select {
			case cmd := <-dn.cmdQueue:
				if time.Since(cmd.StartTimestamp) > cmd.Timeout {
					go cmd.CallBack(cmd.ID, fmt.Errorf("cmd queue timeout!start time:%v,timeout:%v", cmd.StartTimestamp, cmd.Timeout))
				} else {
					dn.cmdQueue <- cmd
				}
			default:
				// No commands in queue, continue to next DN
			}
		}
		mn.dataNodeLock.RUnlock()
	}
}
func (mn *MateNode) monitorWaitResponseTimeoutCmd() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if mn.stop.Load() {
			return
		}
		mn.dataNodeLock.RLock()
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
		mn.dataNodeLock.RUnlock()
	}
}

func (mn *MateNode) lockChunk(ChunkId uint64) {
	mu, _ := mn.chunkMetaMutex.LoadOrStore(ChunkId, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
}

func (mn *MateNode) unlockChunk(ChunkId uint64) {
	mu, ok := mn.chunkMetaMutex.Load(ChunkId)
	if ok {
		mu.(*sync.Mutex).Unlock()
	}
}

func (mn *MateNode) allocateDn(excludes []string) (string, error) {
	mn.dataNodeLock.RLock()
	availableDNs := make([]string, 0, len(mn.dataNodes))
	for uuid, dn := range mn.dataNodes {
		// Skip excluded and dead DNs
		dn.mutex.Lock()
		defer dn.mutex.Unlock()
		if dn.state == DEAD {
			continue
		}
		excluded := false
		for _, exclude := range excludes {
			if uuid == exclude {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		availableDNs = append(availableDNs, uuid)
	}
	mn.dataNodeLock.RUnlock()

	if len(availableDNs) == 0 {
		return "", common.ErrNoAvailableDN
	}

	// Randomly select one DN
	selectedIndex := rand.Intn(len(availableDNs))
	return availableDNs[selectedIndex], nil
}
