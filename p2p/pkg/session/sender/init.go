package sender

import (
	"github.com/pion/webrtc/v4"
)

const (
	bufferThreshold = 512 * 1024 // 512kB
)

func (s *Session) Initialize() error {
	if s.sess.Initialized {
		return nil
	}
	if err := s.sess.CreateConnection(s.onConnectionStateChange()); err != nil {
		return err
	}
	if err := s.createDataChannel(); err != nil {
		return err
	}
	if err := s.sess.CreateOffer(); err != nil {
		return err
	}
	s.sess.Initialized = true
	return nil
}

func (s *Session) Start() error {
	if err := s.Initialize(); err != nil {
		return err
	}
	go s.readFile()
	if err := s.sess.ReadSDP(); err != nil {
		return err
	}
	<-s.sess.Done
	s.sess.OnCompletion()
	return nil
}

// CreateDataChannel that will be used to send data
func (s *Session) createDataChannel() error {
	ordered := true
	maxPacketLifeTime := uint16(10000)
	dataChannel, err := s.sess.CreateDataChannel(&webrtc.DataChannelInit{
		Ordered:           &ordered,
		MaxPacketLifeTime: &maxPacketLifeTime,
	})
	if err != nil {
		return err
	}
	s.dataChannel = dataChannel
	s.dataChannel.OnBufferedAmountLow(s.onBufferedAmountLow())
	s.dataChannel.SetBufferedAmountLowThreshold(bufferThreshold)
	s.dataChannel.OnOpen(s.onOpenHandler())
	s.dataChannel.OnClose(s.onCloseHandler())
	return nil
}
