package sdk

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"github.com/Csy12139/Vesper/p2p"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// TODO：文件怎么加密
func (sdk *Sdk) StoreFile(FileName string, ReplicaNumber int, IsEncrypted bool, concurrencyNumber int) error {
	fileId, err := sdk.mnClient.AddFileMeta()
	if err != nil {
		return fmt.Errorf("AddFileMeta error:%v", err)
	}
	filePath := filepath.Join(sdk.filePath, FileName)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file error [%v]", err)
	}
	defer file.Close() //TODO:注意文件关闭时间

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat file error [%v]", err)
	}
	fileSize := fileInfo.Size()
	// Control the number of concurrent tasks
	sem := make(chan struct{}, concurrencyNumber)
	wg := sync.WaitGroup{}
	for offset := int64(0); offset < fileSize; offset += common.ChunkSize {
		if err != nil {
			return fmt.Errorf("seek file error [%v]", err)
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(chunkOffset int64) {
			defer func() {
				<-sem
				wg.Done()
			}()
			chunkBuffer := make([]byte, common.ChunkSize)
			_, err = file.Seek(chunkOffset, 0)
			if err != nil {
				log.Fatalf("seek file error [%v], chunkOffset[%v]", err, chunkOffset)
				return
			}
			bytes, err := file.Read(chunkBuffer)
			if err != nil {
				if err == io.EOF {
					// no action
				} else {
					log.Fatalf("read file error [%v], chunkOffset[%v]", err, chunkOffset)
					return
				}
			}
			//Fill the data with a size less than 32MB.
			if bytes < common.ChunkSize {
				for j := bytes; j < common.ChunkSize; j++ {
					chunkBuffer[j] = 0
				}
			}
			chunkId, err := sdk.mnClient.AddChunkMeta()
			if err != nil {
				log.Fatalf("Failed to add chunk meta: %v", err)
			}
			var excludes []string
			for i := 0; i < ReplicaNumber; i++ {
				dnId, err := sdk.mnClient.AllocateDnForChunk(chunkId, excludes)
				if err != nil {
					log.Fatalf("Failed to allocate dn for chunk: %v", err)
				}
				log.Debugf("Allocated dn %v for chunk %v", dnId, chunkId)
				excludes = append(excludes, dnId)
				err = sdk.mnClient.AddChunkOnDN(sdk.UUID, dnId)
				if err != nil {
					log.Fatalf("Failed to add chunk on dn: %v", err)
				}
				trans, err := p2p.NewDataTransfer(sdk.UUID, sdk.MNAddr)
				if err != nil {
					log.Fatalf("Failed to create DataTransfer: %v", err)
				}
			}

		}(offset)
	}
	wg.Wait()
	return nil
}

func (sdk *Sdk) ReadFile(FileName string) error {

}

func (sdk *Sdk) DeleteFile(FileName string) error {

}
