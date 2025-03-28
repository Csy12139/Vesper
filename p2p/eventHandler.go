package p2p

import (
	"github.com/Csy12139/Vesper/log"
	"github.com/pion/webrtc/v4"
)

func (s *session) monitorConnectState() {
	s.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Info("connection state:", connectionState)
		switch connectionState {
		case webrtc.ICEConnectionStateFailed:
			s.ctxCancel()
		default:
		}
	})
}

func (s *session) monitorGatherState() {
	s.conn.OnICEGatheringStateChange(func(state webrtc.ICEGatheringState) {
		log.Info("gathering state:", state.String())
		if state == webrtc.ICEGatheringStateComplete {
			close(s.gatherDone)
		}
	})

}

func (s *session) candidatesHandler() {
	s.conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			log.Info("find ICECandidate is nil")
			return
		}
		log.Infof("find a new ICECandidate: %+v", candidate.String())
		s.candidatesLock.Lock()
		defer s.candidatesLock.Unlock()
		s.candidates = append(s.candidates, candidate.ToJSON())
	})
}

func (s *session) dataChHandler(ch *webrtc.DataChannel) {
	s.dataCh = ch
	ch.OnOpen(func() {
		s.dataChOpen <- struct{}{}
	})
	s.dataCh.OnClose(func() {
		s.ctxCancel()
	})

	s.dataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		s.msgCh <- &msg
	})
}
