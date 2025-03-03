package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v4"
	"io"
	"strings"
)

func MustReadStream(stream io.Reader) (string, error) { return "", nil }

// MustReadStream blocks until input is received from the stream
func MustReadStream(stream io.Reader) (string, error) {
	r := bufio.NewReader(stream)
	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				return "", err
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}
	fmt.Println("")
	return in, nil
}

func Encode(sdp *webrtc.SessionDescription) (string, error) {
	sdpJSON, err := json.Marshal(sdp)
	if err != nil {
		return "", err
	}
	var gzbuf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&gzbuf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err := gz.Write(sdpJSON); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(gzbuf.Bytes()), nil
}
func Decode(in string) (*webrtc.SessionDescription, error) {
	buf, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	gz, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	sdpBytes, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	var sdp webrtc.SessionDescription
	err = json.Unmarshal(sdpBytes, &sdp)
	if err != nil {
		return nil, err
	}

	return &sdp, nil
}
