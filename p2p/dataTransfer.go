package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"sync"
	"time"
)

type DataTransfer struct {
	sourceUUID          string
	mnAddr              string
	mnClient            *common.MNClient
	p2pConnectionsMutex sync.Mutex
	p2pConnections      map[string]*P2PConnection
	uuid2Mutex          sync.Map
}

func NewDataTransfer(SourceUUID string, MNAddr string) (*DataTransfer, error) {
	MNClient, err := common.NewMNClient(MNAddr)
	if err != nil {
		return nil, err
	}
	return &DataTransfer{
		sourceUUID:     SourceUUID,
		mnAddr:         MNAddr,
		mnClient:       MNClient,
		p2pConnections: make(map[string]*P2PConnection),
	}, nil
}

func (t *DataTransfer) Send(targetUUID string, dataName string, data []byte, ctx context.Context) error {
	if ctx.Err() != nil {
		return fmt.Errorf("sendDate failed: %v", ctx.Err())
	}

	t.LockUuid(targetUUID)
	defer t.UnlockUuid(targetUUID)
	t.p2pConnectionsMutex.Lock()
	p2p, exist := t.p2pConnections[targetUUID]
	t.p2pConnectionsMutex.Unlock()
	if !exist {
		var err error
		p2p, err = NewP2PConnection()
		if err != nil {
			return err
		}
	} else if !p2p.IsConnection() {
		err := p2p.CloseConnection()
		if err != nil {
			return err
		}
		p2p, err = NewP2PConnection()
		if err != nil {
			return err
		}
	}
	if !p2p.IsConnection() {
		offerSdp, offerCandidate, err := p2p.GetSenderSDPCandidates()
		if err != nil {
			return err
		}
		log.Debug("get sender sdp and candidate")
		err = t.mnClient.PutSDPCandidates(t.sourceUUID, targetUUID, offerSdp, offerCandidate)
		if err != nil {
			return err
		}
		log.Debug("put sender sdp and candidate")
		var answerSdp string
		var answerCandidate []string
		resultChan := make(chan struct{})
		go func() {
			for {
				answerSdp, answerCandidate, err = t.mnClient.GetSDPCandidates(targetUUID, t.sourceUUID)
				if err != nil {
					if errors.Is(common.ErrSDPNotFound, err) {
						time.Sleep(100 * time.Millisecond)
						continue
					} else {
						log.Debugf("common.ErrSDPNotFound [%v]", common.ErrSDPNotFound)
						log.Errorf("get sender sdp and candidates error [%v]", err)
						return
					}
				} else {
					close(resultChan)
					return
				}
			}
		}()
		select {
		case <-resultChan:
			if err != nil {
				return fmt.Errorf("get sender SDPCandidates failed: %v", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("get sender SDPCandidates failed: %v", ctx.Err())
		}

		err = p2p.SetRemotedDescription(answerSdp, answerCandidate)
		if err != nil {
			return fmt.Errorf("set sender remoted description failed: %v", err)
		}

		err = p2p.WaitConnection(ctx)
		if err != nil {
			return err
		}
		t.p2pConnectionsMutex.Lock()
		t.p2pConnections[targetUUID] = p2p
		t.p2pConnectionsMutex.Unlock()
	}
	err := p2p.SendDate(dataName, data, ctx)
	if err != nil {
		return fmt.Errorf("sendDate failed: %v", err)
	}
	return nil
}
func (t *DataTransfer) StartToReceive(targetUUID string, ctx context.Context, callback func(targetUUID string, label string, data []byte)) error {
	if ctx.Err() != nil {
		return fmt.Errorf("sendDate failed: %v", ctx.Err())
	}

	t.LockUuid(targetUUID)
	defer t.UnlockUuid(targetUUID)
	t.p2pConnectionsMutex.Lock()
	p2p, exist := t.p2pConnections[targetUUID]
	t.p2pConnectionsMutex.Unlock()
	if !exist {
		var err error
		p2p, err = NewP2PConnection()
		if err != nil {
			return fmt.Errorf("new p2p connection failed: %v", err)
		}
	} else if !p2p.IsConnection() {
		err := p2p.CloseConnection()
		if err != nil {
			return fmt.Errorf("close p2p connection failed: %v", err)
		}
		p2p, err = NewP2PConnection()
		if err != nil {
			return fmt.Errorf("new p2p connection failed: %v", err)
		}
	}
	if !p2p.IsConnection() {
		var offerSdp string
		var offerCandidate []string
		var err error
		resultChan := make(chan struct{})
		go func() {
			for {
				offerSdp, offerCandidate, err = t.mnClient.GetSDPCandidates(targetUUID, t.sourceUUID)
				if err != nil {
					if errors.Is(common.ErrSDPNotFound, err) {
						continue
					} else {
						return
					}
				} else {
					close(resultChan)
					return
				}
			}
		}()
		select {
		case <-resultChan:
			if err != nil {
				return fmt.Errorf("get sender SDPCandidates failed: %v", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("get sender SDPCandidates failed: %v", ctx.Err())
		}
		answerSdp, answerCandidate, err := p2p.SwapReceiverSDPCandidates(offerSdp, offerCandidate)
		if err != nil {
			return fmt.Errorf("swap sender SDPCandidates failed: %v", err)
		}
		err = t.mnClient.PutSDPCandidates(t.sourceUUID, targetUUID, answerSdp, answerCandidate)
		if err != nil {
			return fmt.Errorf("put sender SDPCandidates failed: %v", err)
		}
		err = p2p.WaitConnection(ctx)
		if err != nil {
			return fmt.Errorf("wait sender SDPCandidates failed: %v", err)
		}
		t.p2pConnectionsMutex.Lock()
		t.p2pConnections[targetUUID] = p2p
		t.p2pConnectionsMutex.Unlock()

		handle := func(label string, data []byte) {
			callback(targetUUID, label, data)
		}
		err = p2p.RegisterReceiveDataCallback(handle)
		if err != nil {
			return fmt.Errorf("register receive data failed: %v", err)
		}
	}
	return nil
}

func (t *DataTransfer) LockUuid(targetUUID string) {
	mu, _ := t.uuid2Mutex.LoadOrStore(targetUUID, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
}
func (t *DataTransfer) UnlockUuid(targetUUID string) {
	mu, ok := t.uuid2Mutex.Load(targetUUID)
	if ok {
		mu.(*sync.Mutex).Unlock()
	}
}
