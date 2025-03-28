package p2p

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"github.com/pion/webrtc/v4"
	"io"
)

type SdpInfo struct {
	Sdp *webrtc.SessionDescription `json:"sdp"`
}
type CandidatesInfo struct {
	Candidates []webrtc.ICECandidateInit `json:"candidates"`
}

func encodeSDP(sdp *webrtc.SessionDescription) (string, error) {
	info := SdpInfo{
		Sdp: sdp,
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	g, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	defer g.Close()
	if _, err = g.Write(infoJSON); err != nil {
		return "", err
	}

	if err = g.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func encodeCandidates(candidates []webrtc.ICECandidateInit) ([]string, error) {
	result := make([]string, 0, len(candidates))

	for _, candidate := range candidates {
		info := CandidatesInfo{
			Candidates: []webrtc.ICECandidateInit{candidate},
		}
		infoJSON, err := json.Marshal(info)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		g, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		if err != nil {
			return nil, err
		}
		defer g.Close()
		if _, err = g.Write(infoJSON); err != nil {
			return nil, err
		}
		if err = g.Close(); err != nil {
			return nil, err
		}

		encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
		result = append(result, encoded)
	}

	return result, nil
}

func decodeSDP(in string) (*webrtc.SessionDescription, error) {
	buf, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	infoBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var info SdpInfo
	err = json.Unmarshal(infoBytes, &info)
	if err != nil {
		return nil, err
	}
	return info.Sdp, nil
}

func decodeCandidates(in []string) ([]webrtc.ICECandidateInit, error) {
	var candidates []webrtc.ICECandidateInit

	for _, encoded := range in {
		buf, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, err
		}
		r, err := gzip.NewReader(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		infoBytes, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		var info CandidatesInfo
		err = json.Unmarshal(infoBytes, &info)
		if err != nil {
			return nil, err
		}
		if len(info.Candidates) > 0 {
			candidates = append(candidates, info.Candidates[0])
		}
	}
	return candidates, nil
}
