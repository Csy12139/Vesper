package sender

import (
	"github.com/pion/webrtc/v4"
)

func (s *Session) onConnectionStateChange() func(connectionState webrtc.ICEConnectionState) {
	return func(connectionState webrtc.ICEConnectionState) {
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			s.stopSending <- struct{}{}
		}
	}
}
func (s *Session) onCloseHandler() func() {
	return func() {
		s.close(true)
	}
}

func (s *Session) onOpenHandler() func() {
	return func() {
		s.sess.NetworkStats.Start()
		s.writeToNetwork()
	}
}
func (s *Session) close(calledFromCloseHandler bool) {
	if !calledFromCloseHandler {
		s.dataChannel.Close()
	}
}
