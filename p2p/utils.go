package p2p

import (
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v4"
)

func encodeSDP(sdp *webrtc.SessionDescription) (string, error) {
	jsonSdp, err := json.Marshal(sdp)
	if err != nil {
		return "", fmt.Errorf("failed to encode sdp [%v] err [%v]", sdp, err)
	}
	return string(jsonSdp), nil
}
func encodeCandidate(candidate webrtc.ICECandidateInit) (string, error) {
	jsonCandidate, err := json.Marshal(candidate)
	if err != nil {
		return "", fmt.Errorf("failed to encode candidate [%v] err [%v]", candidate, err)
	}
	return string(jsonCandidate), nil
}

func encodeCandidates(candidates []webrtc.ICECandidateInit) ([]string, error) {
	result := make([]string, 0, len(candidates))
	for i := range candidates {
		stringCandidate, err := encodeCandidate(candidates[i])
		if err != nil {
			return nil, err
		}
		result = append(result, stringCandidate)
	}
	return result, nil
}

func decodeSDP(jsonSdp string) (*webrtc.SessionDescription, error) {
	var sdp webrtc.SessionDescription
	err := json.Unmarshal([]byte(jsonSdp), &sdp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sdp [%v] err [%v]", jsonSdp, err)
	}
	return &sdp, nil
}
func decodeCandidate(jsonCandidate string) (*webrtc.ICECandidateInit, error) {
	var candidate webrtc.ICECandidateInit
	err := json.Unmarshal([]byte(jsonCandidate), &candidate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode candidate [%v] err [%v]", jsonCandidate, err)
	}
	return &candidate, nil
}
func decodeCandidates(jsonCandidates []string) ([]webrtc.ICECandidateInit, error) {
	var candidates []webrtc.ICECandidateInit
	for i := range jsonCandidates {
		candidate, err := decodeCandidate(jsonCandidates[i])
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, *candidate)
	}
	return candidates, nil
}

//func put(req *common.PutSDPCandidatesRequest) {
//	mnClient, err := common.NewMNClient(GlobalConfig.MNAddr)
//	if err != nil {
//		log.Fatalf("Failed to create MN client: %v", err)
//	}
//	resp, err := mnClient.PutSDPCandidates(req)
//	if err != nil {
//		log.Errorf("PutSDPCandidates failed: %v", err)
//	}
//	log.Infof("PutSDPCandidates: %v", resp)
//}
//
//func get(req *common.GetSDPCandidatesRequest) (string, []string) {
//	mnClient, err := common.NewMNClient(GlobalConfig.MNAddr)
//	if err != nil {
//		log.Fatalf("Failed to create MN client: %v", err)
//	}
//	resp, err := mnClient.GetSDPCandidates(req)
//	if err != nil {
//		log.Errorf("GetSenderSDPCandidates failed: %v", err)
//	}
//	log.Infof("GetSenderSDPCandidates: %v", resp)
//	return resp.SDP, resp.Candidates
//}
