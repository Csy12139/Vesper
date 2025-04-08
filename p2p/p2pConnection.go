package p2p

import (
	"context"
	"fmt"
	"github.com/Csy12139/Vesper/log"
	"github.com/pion/webrtc/v4"
	"sync"
	"time"
)

type P2PConnection struct {
	mnAddr string

	conn *webrtc.PeerConnection

	//dataCh *webrtc.DataChannel
	//msgCh  chan webrtc.DataChannelMessage

	candidates     []webrtc.ICECandidateInit
	candidatesLock sync.Mutex
	gatherDone     chan struct{}
}

func NewP2PConnection() (*P2PConnection, error) {
	p := &P2PConnection{
		gatherDone: make(chan struct{}, 1),
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
	p.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Info("connection state:", connectionState)
		switch connectionState {
		case webrtc.ICEConnectionStateFailed:
			err := p.conn.GracefulClose()
			if err != nil {
				log.Error("Error closing peer connection:", err)
			}
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

func (p *P2PConnection) CloseConnection() error {
	return p.conn.GracefulClose()
}

func (p *P2PConnection) WaitConnection(ctx context.Context) error {
	for !p.IsConnection() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

func (p *P2PConnection) SendDate(label string, data []byte, ctx context.Context) error {
	if ctx.Err() != nil {
		return fmt.Errorf("sendDate failed: %v", ctx.Err())
	}

	ch, err := p.conn.CreateDataChannel(label, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := ch.GracefulClose(); err != nil {
			log.Errorf("Error closing DataChannel %s: %v", ch.Label(), err)
		}
	}()

	ch.OnClose(func() {
		log.Info("DataChannel closed:", ch.Label())
	})
	// Wait DataChannel Open
	for !(ch.ReadyState() == webrtc.DataChannelStateOpen) {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait dataChannel open failed: %v", ctx.Err())
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	// Send data
	const bufferedAmountLowThreshold = 1 * 1024 * 1024
	offset := 0
	var end int
	dataLen := len(data)
	sendNextChunk := func() {
		if offset >= dataLen {
			log.Info("All data sent to buffer")
			return
		}
		end = offset + bufferedAmountLowThreshold
		if end > dataLen {
			end = dataLen
		}
		chunk := data[offset:end]

		if err := ch.Send(chunk); err != nil {
			log.Errorf("Send chunk %d-%d failed: %v", offset, end, err)
			return
		}
		log.Infof("Sent chunk %d-%d bytes", offset, end)
		offset = end
	}
	ch.SetBufferedAmountLowThreshold(bufferedAmountLowThreshold)
	ch.OnBufferedAmountLow(func() {
		log.Infof("BufferedAmountLow triggered, current: %d", ch.BufferedAmount())
		sendNextChunk()
	})
	// Initial sending
	sendNextChunk()
	// Wait data send over and wait buffer decrease to zero
	for ch.BufferedAmount() > 0 || end < dataLen {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait buffer clear failed: %v", ctx.Err())
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	log.Infof("All data sent.")
	return nil
}

func (p *P2PConnection) RegisterReceiveDataCallback(callback func(label string, data []byte)) error {
	p.conn.OnDataChannel(func(ch *webrtc.DataChannel) {
		msgCh := make(chan webrtc.DataChannelMessage, 200)
		if ch.Label() == "init" {
			return
		}
		ch.OnOpen(func() {
			log.Infof("%s dataChannel opened", ch.Label())
		})
		ch.OnClose(func() {
			log.Infof("%s dataChannel closed", ch.Label())
			close(msgCh)
		})
		ch.OnMessage(func(msg webrtc.DataChannelMessage) {
			select {
			case msgCh <- msg:
			default:
				log.Warnf("msgCh full or closed for %s", ch.Label())
			}
		})
		go func() {
			var receivedData []byte
			for msg := range msgCh {
				log.Infof("Received data of %d bytes on %s", len(msg.Data), ch.Label())
				receivedData = append(receivedData, msg.Data...)
			}
			callback(ch.Label(), receivedData)
		}()
	})
	return nil
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

func (p *P2PConnection) waitGatherComplete() {
	<-p.gatherDone
}
