package p2p

import (
	"fmt"
	"github.com/Csy12139/Vesper/log"
	"github.com/pion/webrtc/v4"
	"sync"
	"time"
)

type P2PConnection struct {
	mnAddr string

	conn *webrtc.PeerConnection

	dataCh   *webrtc.DataChannel
	msgChMap map[string]chan *webrtc.DataChannelMessage
	mu       sync.Mutex

	candidates     []webrtc.ICECandidateInit
	candidatesLock sync.Mutex
	gatherDone     chan struct{}

	SDPCandidatesExchangeTimeout time.Duration
}

func NewP2PConnection() (*P2PConnection, error) {
	p := &P2PConnection{
		gatherDone: make(chan struct{}, 1),
		msgChMap:   make(map[string]chan *webrtc.DataChannelMessage),
	}

	conn, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:101.34.238.168:3478"},
			},
			{
				URLs:           []string{"turn:101.34.238.168:3478"},
				Username:       "ranber",
				Credential:     "12138",
				CredentialType: webrtc.ICECredentialTypePassword,
			}},
	})
	if err != nil {
		return nil, err
	}
	p.conn = conn
	p.registerConnectionCallback()

	_, err = p.conn.CreateDataChannel("init", nil)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *P2PConnection) GetSenderSDPCandidates() (string, []string, error) {
	initialSDPOffer, err := p.conn.CreateOffer(nil)
	if err != nil {
		return "", nil, err
	}
	OfferSdp, err := p.getLocalDescription(initialSDPOffer)
	if err != nil {
		return "", nil, err
	}
	OfferCandidates, err := p.getCandidates()
	if err != nil {
		return "", nil, err
	}
	return OfferSdp, OfferCandidates, nil
}

func (p *P2PConnection) SwapReceiverSDPCandidates(sdp string, candidates []string) (string, []string, error) {
	err := p.SetRemotedDescription(sdp, candidates)
	if err != nil {
		return "", nil, err
	}
	initAnswer, err := p.conn.CreateAnswer(nil)
	if err != nil {
		return "", nil, err
	}
	answerSdp, err := p.getLocalDescription(initAnswer)
	if err != nil {
		return "", nil, err
	}
	answerCandidates, err := p.getCandidates()
	if err != nil {
		return "", nil, err
	}
	return answerSdp, answerCandidates, nil
}

func (p *P2PConnection) IsConnection() bool {
	return p.conn.ICEConnectionState() == webrtc.ICEConnectionStateConnected
	//return (p.conn.ICEConnectionState() == webrtc.ICEConnectionStateConnected) && (p.dataCh.ReadyState() == webrtc.DataChannelStateOpen)
}

func (p *P2PConnection) WaitConnection(timeout time.Duration) error {
	timeoutCh := time.After(timeout)
	for !p.IsConnection() {
		select {
		case <-timeoutCh:
			return fmt.Errorf("wait connection timed out after %v", timeout)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

func (p *P2PConnection) SendDate(label string, data []byte, timeout time.Duration) error {
	ch, err := p.conn.CreateDataChannel(label, nil)
	if err != nil {
		return err
	}
	p.registerChannelCallback(label, ch)

	for !(ch.ReadyState() == webrtc.DataChannelStateOpen) {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("wait connection timed out after %v", timeout)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- ch.Send(data)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Info("send data error: ", err)
			return fmt.Errorf("send data failed: %w", err)
		}
		return nil
	case <-time.After(timeout):
		log.Info("send data timeout")
		return fmt.Errorf("send data timed out after %v", timeout)
	}

}

func (p *P2PConnection) RegisterReceiveDataCallback(label string, callback func(label string, data []byte)) error {
	p.mu.Lock()
	p.msgChMap[label] = make(chan *webrtc.DataChannelMessage, 200)
	msgCh := p.msgChMap[label]
	p.mu.Unlock()

	p.conn.OnDataChannel(func(ch *webrtc.DataChannel) {
		if ch.Label() == label {
			p.registerChannelCallback(label, ch)

			go func() {
				for msg := range msgCh {
					log.Infof("Received data of %d bytes on %s", len(msg.Data), label)
					callback(label, msg.Data)
				}
				log.Info("Message channel closed for:", label)
				ch.Close()
			}()
		}
	})
	return nil
}

func (p *P2PConnection) CloseConnection() {
	if p.conn == nil {
		err := p.conn.GracefulClose()
		if err != nil {
			log.Error("Error closing peer connection:", err)
		}
	}
}

func (p *P2PConnection) SetRemotedDescription(remotedSdp string, remotedCandidates []string) error {
	remoteSDP, err := decodeSDP(remotedSdp)
	if err != nil {
		return err
	}
	remoteCandidates, err := decodeCandidates(remotedCandidates)
	if err != nil {
		return err
	}
	err = p.conn.SetRemoteDescription(*remoteSDP)
	if err != nil {
		return err
	}
	for i := range remoteCandidates {
		err := p.conn.AddICECandidate(remoteCandidates[i])
		if err != nil {
			log.Errorf("add ICECandidate error: %v candidate:%v", err, remoteCandidates[i])
			continue
		} else {
			log.Debugf("add ICECandidate success candidate:%v", remoteCandidates[i])
		}
	}
	return nil
}
func (p *P2PConnection) getCandidates() ([]string, error) {
	p.waitGatherComplete()
	encodedCandidates, err := encodeCandidates(p.candidates)
	if err != nil {
		return nil, err
	}

	return encodedCandidates, nil
}

func (p *P2PConnection) getLocalDescription(Description webrtc.SessionDescription) (string, error) {
	err := p.conn.SetLocalDescription(Description)
	if err != nil {
		return "", err
	}
	sdp := p.conn.LocalDescription()
	encodedSDP, err := encodeSDP(sdp)
	if err != nil {
		return "", err
	}
	return encodedSDP, nil
}

func (p *P2PConnection) registerConnectionCallback() {
	p.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Info("connection state:", connectionState)
		switch connectionState {
		case webrtc.ICEConnectionStateFailed:
			p.CloseConnection()
		default:
		}
	})
	p.conn.OnICEGatheringStateChange(func(state webrtc.ICEGatheringState) {
		log.Info("gathering state:", state.String())
		if state == webrtc.ICEGatheringStateComplete {
			close(p.gatherDone)
		}
	})
	p.conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			log.Info("find ICECandidate is nil")
			return
		}
		log.Infof("find a new ICECandidate: %+v", candidate.String())
		p.candidatesLock.Lock()
		defer p.candidatesLock.Unlock()
		p.candidates = append(p.candidates, candidate.ToJSON())
	})

}

func (p *P2PConnection) registerChannelCallback(label string, ch *webrtc.DataChannel) {
	// p.dataCh = ch
	ch.OnClose(func() {
		// p.CloseConnection()
		log.Info("DataChannel closed:", ch.Label())
	})
	ch.OnMessage(func(msg webrtc.DataChannelMessage) {
		p.mu.Lock()
		p.msgChMap[label] <- &msg
		p.mu.Unlock()
	})
}

func (p *P2PConnection) waitGatherComplete() {
	<-p.gatherDone
}
