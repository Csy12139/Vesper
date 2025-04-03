package test

import (
	"github.com/Csy12139/Vesper/log"
	"github.com/Csy12139/Vesper/p2p"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

var offerReady = make(chan bool)
var answerReady = make(chan bool)
var testReady = make(chan bool, 1)
var DataSize = 10 * 1024 * 1024

func WriteFile(sdpPath string, CandidatesPath string, sdp string, candidate []string) error {
	if _, err := os.Stat(sdpPath); err == nil {
		if err := os.Remove(sdpPath); err != nil {
			return err
		}
	}
	if _, err := os.Stat(CandidatesPath); err == nil {
		if err := os.Remove(CandidatesPath); err != nil {
			return err
		}
	}

	err := os.WriteFile(sdpPath, []byte(sdp), 0644)
	if err != nil {
		panic("write sdpPath failed: " + err.Error())
	}
	candidateContent := strings.Join(candidate, "\n")
	err = os.WriteFile(CandidatesPath, []byte(candidateContent), 0644)
	if err != nil {
		panic("write CandidatesPath failed: " + err.Error())
	}
	return nil
}
func ReadFile(sdpPath string, CandidatesPath string) (string, []string, error) {
	sdpData, err := os.ReadFile(sdpPath)
	if err != nil {
		panic("read sdpPath failed: " + err.Error())
	}
	sdp := string(sdpData)

	candidateData, err := os.ReadFile(CandidatesPath)
	if err != nil {
		panic("read CandidatesPath failed: " + err.Error())
	}
	candidate := strings.Split(strings.TrimSpace(string(candidateData)), "\n")

	return sdp, candidate, nil
}

func Sender(t *testing.T) {
	pc, err := p2p.NewP2PConnection()
	if err != nil {
		t.Fatal("[Sender]NewP2PConnection failed: " + err.Error())
	}
	offerSdp, offerCandidate, err := pc.GetSenderSDPCandidates()
	if err != nil {
		t.Fatal("[Sender]GetSenderSDPCandidates failed: " + err.Error())
	}
	err = WriteFile("offerSdp.txt", "offerCandidate.txt", offerSdp, offerCandidate)
	if err != nil {
		t.Fatal("[Sender]WriteFile failed: " + err.Error())
	}
	offerReady <- true
	<-answerReady
	answerSdp, answerCandidate, err := ReadFile("answerSdp.txt", "answerCandidate.txt")
	if err != nil {
		t.Fatal("[Sender]ReadFile failed: " + err.Error())
	}
	err = pc.SetRemotedDescription(answerSdp, answerCandidate)
	if err != nil {
		t.Fatal("[Sender]SetRemotedDescription failed: " + err.Error())
	}
	err = pc.WaitConnection(time.Minute)
	if err != nil {
		t.Fatal("[Sender]WaitConnection failed: " + err.Error())
	}
	data := make([]byte, DataSize)
	err = pc.SendDate("test", data, time.Minute)
	if err != nil {
		t.Fatal("[Sender]SendDate failed: " + err.Error())
	}
	//data = bytes.Repeat([]byte{1}, DataSize)
	err = pc.SendDate("another", data, time.Minute)
	if err != nil {
		t.Fatal("[Sender]SendDate failed: " + err.Error())
	}
}

func Receiver(t *testing.T) {
	pc, err := p2p.NewP2PConnection()
	if err != nil {
		t.Fatal("[Receiver]NewP2PConnection failed: " + err.Error())
	}
	<-offerReady
	offerSdp, offerCandidate, err := ReadFile("offerSdp.txt", "offerCandidate.txt")
	if err != nil {
		t.Fatal("[Receiver]ReadFile failed: " + err.Error())
	}
	answerSdp, answerCandidate, err := pc.SwapReceiverSDPCandidates(offerSdp, offerCandidate)
	if err != nil {
		t.Fatal("[Receiver]SwapReceiverSDPCandidates failed: " + err.Error())
	}
	err = WriteFile("answerSdp.txt", "answerCandidate.txt", answerSdp, answerCandidate)
	if err != nil {
		t.Fatal("[Receiver]WriteFile failed: " + err.Error())
	}
	answerReady <- true
	err = pc.WaitConnection(time.Minute)
	if err != nil {
		t.Fatal("[Receiver]WaitConnection failed: " + err.Error())
	}
	testReady <- false
	err = pc.RegisterReceiveDataCallback(func(label string, data []byte) {
		for _, b := range data {
			if b != 0 {
				t.Errorf("[Receiver]The received data from %s is incorrect.", label)
			}
		}
		if len(data) != DataSize {
			t.Error("[Receiver]The received data length is incorrect.")
		}
		testReady <- true
	})
	if err != nil {
		t.Fatal("[Receiver]RegisterReceiveDataCallback failed: " + err.Error())
	}
	for !<-testReady {
		time.Sleep(100 * time.Millisecond)
	}
}

func TestP2P(t *testing.T) {
	err := log.InitLog("./p2p", 100, 10, "INFO")
	if err != nil {
		t.Fatal("InitLog failed: " + err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		Sender(t)
	}()
	go func() {
		defer wg.Done()
		Receiver(t)
	}()

	wg.Wait()
}
