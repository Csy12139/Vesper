package session

import (
	"fmt"
	"github.com/Csy12139/Vesper/p2p/pkg/stats"
	"github.com/Csy12139/Vesper/p2p/pkg/utils"
	"github.com/pion/webrtc/v4"
	"io"
	"os"
)

type CompletionHandler func()

// Session contains common elements to perform send/receive
type Session struct {
	Done           chan struct{}
	sdpInput       io.Reader
	sdpOutput      io.Writer
	peerConnection *webrtc.PeerConnection
	NetworkStats   *stats.Stats
	Initialized    bool
	onCompletion   CompletionHandler
}

// New creates a new Session
func New(sdpInput io.Reader, sdpOutput io.Writer) *Session {
	// 后续改为从文件读取，sdp从服务器转发
	if sdpInput == nil {
		sdpInput = os.Stdin
	}
	if sdpOutput == nil {
		sdpOutput = os.Stdout
	}
	return &Session{
		sdpInput:    sdpInput,
		sdpOutput:   sdpOutput,
		Done:        make(chan struct{}),
		Initialized: false,
	}
}

// CreateConnection prepares a WebRTC connection
func (s *Session) CreateConnection(onConnectionStateChange func(connectionState webrtc.ICEConnectionState)) error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:101.34.238.168:3478"},
			},
			{
				URLs:           []string{"turn:101.34.238.168:3478"},
				Username:       "ranber",
				Credential:     "12138",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}
	s.peerConnection = peerConnection
	s.peerConnection.OnICEConnectionStateChange(onConnectionStateChange)
	return nil
}

func (s *Session) CreateOffer() error {
	//if err := s.iceGatherer.Gather(); err != nil {
	//	return err
	//}
	initialSdpOffer, err := s.peerConnection.CreateOffer(nil)
	if err != nil {
		return err
	}
	if err = s.peerConnection.SetLocalDescription(initialSdpOffer); err != nil {
		return err
	}
	gatherComplete := webrtc.GatheringCompletePromise(s.peerConnection)
	<-gatherComplete

	sdpOffer := s.peerConnection.LocalDescription()
	encodedSdp, err := utils.Encode(sdpOffer)
	if err != nil {
		return err
	}
	fmt.Println("Send this SDP:")
	fmt.Fprintf(s.sdpOutput, "%s\n", encodedSdp)
	return nil
}

// CreateDataChannel that will be used to send data
func (s *Session) CreateDataChannel(c *webrtc.DataChannelInit) (*webrtc.DataChannel, error) {
	return s.peerConnection.CreateDataChannel("data", c)
}

func (s *Session) ReadSDP() error {
	var remoteSDP *webrtc.SessionDescription
	fmt.Println("Please, paste the remote SDP:")
	for {
		encoded, err := utils.MustReadStream(s.sdpInput)
		if err == nil {
			remoteSDP, err = utils.Decode(encoded)
			if err == nil {
				break
			}
		}
	}
	return s.peerConnection.SetRemoteDescription(*remoteSDP)
}

// OnCompletion is called when session ends
func (s *Session) OnCompletion() {
	if s.onCompletion != nil {
		s.onCompletion()
	}
}
