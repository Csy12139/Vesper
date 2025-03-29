package p2p

import (
	"github.com/Csy12139/Vesper/log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

var offerReady = make(chan bool)
var answerReady = make(chan bool)
var callbackReady = make(chan bool, 1)
var receiveData = make(chan bool)
var DataSize = 1 * 1024 * 1024

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
	pc, err := NewP2PConnection()
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
	<-receiveData
}

func Receiver(t *testing.T) {
	pc, err := NewP2PConnection()
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
	err = pc.RegisterReceiveDataCallback("test", func(label string, data []byte) {
		for _, b := range data {
			if b != 0 {
				t.Error("[Receiver]The received data is incorrect.")
			}
		}
		if len(data) != DataSize {
			t.Error("[Receiver]The received data length is incorrect.")
		}
		t.Logf("[Receiver]The received data lable is %s.", label)
		callbackReady <- true
	})
	if err != nil {
		t.Fatal("[Receiver]RegisterReceiveDataCallback failed: " + err.Error())
	}
	<-callbackReady
	receiveData <- true
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
