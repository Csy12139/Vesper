package sender

import (
	internalSess "github.com/Csy12139/Vesper/p2p/internal/session"
	"github.com/Csy12139/Vesper/p2p/pkg/stats"
	"github.com/pion/webrtc/v4"
	"io"
)

type outputMsg struct {
	n    int
	buff []byte
}
type Session struct {
	sess         *internalSess.Session
	stream       io.Reader
	stopSending  chan struct{}
	dataChannel  *webrtc.DataChannel
	output       chan outputMsg
	msgToBeSent  []outputMsg
	readingStats *stats.Stats
	dataBuff     []byte
}

func newSession(s *internalSess.Session, f io.Reader) *Session {
	return &Session{
		sess:        s,
		stream:      f,
		stopSending: make(chan struct{}, 1),
	}
}

func New(f io.Reader) *Session {
	return newSession(internalSess.New(nil, nil), f)
}
