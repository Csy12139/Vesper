package SDK

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/p2p"
	"sync"
)

type Sdk struct {
	UUID         string
	MNAddr       string
	mnClient     *common.MNClient
	dataTransfer *p2p.DataTransfer
	filePath     string
	filesMap     map[string]uint64
	filesMapMux  sync.Mutex
}

func NewSdk(UUID string, MNAddr string, FilePath string) (*Sdk, error) {
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
		filesMap:     make(map[string]uint64),
		filesMapMux:  sync.Mutex{},
	}
	return &sdk, nil
}

func (sdk *Sdk) Start() {}
