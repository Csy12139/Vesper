package sdk

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/p2p"
)

type Sdk struct {
	UUID         string
	MNAddr       string
	mnClient     *common.MNClient
	dataTransfer *p2p.DataTransfer
	filePath     string
	Files        map[string]uint64 // 需要fileID吗
}

func NewSdk(UUID string, MNAddr string, FilePath string, ConcurrencyNumber uint64) (*Sdk, error) {
	mnClient, err := common.NewMNClient(MNAddr)
	if err != nil {
		return nil, fmt.Errorf("NewMNClient error:%v", err)
	}
	trans, err := p2p.NewDataTransfer(UUID, MNAddr)
	if err != nil {
		return nil, fmt.Errorf("NewDataTransfer error:%v", err)
	}
	sdk := Sdk{
		UUID:         UUID,
		MNAddr:       MNAddr,
		mnClient:     mnClient,
		dataTransfer: trans,
		filePath:     FilePath,
		Files:        make(map[string]uint64),
	}
	return &sdk, nil
}

func (sdk *Sdk) Start() {}
