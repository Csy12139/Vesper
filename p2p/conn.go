package p2p

import (
	"context"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"github.com/pion/webrtc/v4"
	"sync"
)

type session struct {
	mnAddr string

	conn       *webrtc.PeerConnection
	dataCh     *webrtc.DataChannel
	dataChOpen chan struct{}

	ctx       context.Context
	ctxCancel context.CancelFunc

	msgCh          chan *webrtc.DataChannelMessage
	candidates     []webrtc.ICECandidateInit
	candidatesLock sync.Mutex
	gatherDone     chan struct{}
}

func newSession(MNAddr string) *session {
	ctx, cancel := context.WithCancel(context.Background())
	return &session{
		mnAddr:     MNAddr,
		ctx:        ctx,
		ctxCancel:  cancel,
		dataChOpen: make(chan struct{}, 10),
		gatherDone: make(chan struct{}, 1),
		msgCh:      make(chan *webrtc.DataChannelMessage, 200),
	}
}

func (s *session) setupP2PConnection() error {
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
		return err
	}
	s.conn = conn
	s.monitorConnectState()
	s.monitorGatherState()
	go s.connClose()
	s.candidatesHandler()
	return nil
}

func (s *session) createDataChannel() error {
	ch, err := s.conn.CreateDataChannel("data", nil)
	if err != nil {
		return err
	}
	s.dataChHandler(ch)
	return nil
}
func (s *session) createOffer() (string, []string, error) {
	initialSDPOffer, err := s.conn.CreateOffer(nil)
	if err != nil {
		return "", nil, err
	}
	err = s.conn.SetLocalDescription(initialSDPOffer)
	if err != nil {
		return "", nil, err
	}
	sdpOffer := s.conn.LocalDescription()
	s.waitGatherComplete()

	encodedSDP, err := encodeSDP(sdpOffer)
	if err != nil {
		return "", nil, err
	}
	encodedCandidates, err := encodeCandidates(s.candidates)
	if err != nil {
		return "", nil, err
	}
	return encodedSDP, encodedCandidates, nil
}
func (s *session) waitGatherComplete() {
	<-s.gatherDone
}
func (s *session) closeP2PConnection() error {
	s.ctxCancel()
	return nil
}
func (s *session) initP2PConnection() (string, []string, error) {
	err := s.setupP2PConnection()
	if err != nil {
		return "", nil, err
	}
	err = s.createDataChannel()
	if err != nil {
		return "", nil, err
	}

	sdpOffer, candidates, err := s.createOffer()
	if err != nil {
		return "", nil, err
	}

	return sdpOffer, candidates, nil
}
func (s *session) putRequest(req *common.PutSDPCandidatesRequest) (bool, error) {
	mnClient, err := common.NewMNClient(s.mnAddr)
	if err != nil {
		log.Fatalf("Failed to create MN client: %v", err)
	}
	resp, err := mnClient.PutSDPCandidates(req)
	if err != nil {
		log.Errorf("PutSDPCandidates failed: %v", err)
	}
	log.Infof("PutSDPCandidates: %v", resp)
	return resp.Success, err
}
func (s *session) getRequest(req *common.GetSDPCandidatesRequest) (string, []string) {
	mnClient, err := common.NewMNClient(s.mnAddr)
	if err != nil {
		log.Fatalf("Failed to create MN client: %v", err)
	}
	resp, err := mnClient.GetSDPCandidates(req)
	if err != nil {
		log.Errorf("GetSDPCandidates failed: %v", err)
	}
	log.Infof("GetSDPCandidates: %v", resp)
	return resp.SDP, resp.Candidates
}
func (s *session) exchangeSDPCandidates(sourceUUID string, targetUUID string, sdpOffer string, candidates []string) {
	putRequest := common.PutSDPCandidatesRequest{
		SourceUUID: sourceUUID,
		TargetUUID: targetUUID,
		SDP:        sdpOffer,
		Candidates: candidates,
	}
	s.putRequest(&putRequest)

}
func (s *session) connClose() {
	<-s.ctx.Done()

	if s.dataCh != nil {
		err := s.dataCh.GracefulClose()
		if err != nil {
			log.Error("Error closing control channel:", err)
		}
	}
	if s.conn == nil {
		err := s.conn.GracefulClose()
		if err != nil {
			log.Error("Error closing peer connection:", err)
		}
	}
}
