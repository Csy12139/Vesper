package SDK

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// TODO：文件怎么加密
func (sdk *Sdk) StoreFile(FileName string, ReplicaNumber int, IsEncrypted bool, concurrencyNumber int) error {
	filePath := filepath.Join(sdk.filePath, FileName)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file error [%v]", err)
	}
	defer file.Close()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat file error [%v]", err)
	}
	fileSize := fileInfo.Size()
	// Control the number of concurrent tasks
	sem := make(chan struct{}, concurrencyNumber)
	wg := sync.WaitGroup{}
	ChunkIds := make(map[uint64]int64)
	for offset, index := int64(0), int64(0); offset < fileSize; offset += common.ChunkSize {
		chunkId, err := sdk.mnClient.AddChunkMeta()
		if err != nil {
			log.Fatalf("Failed to add chunk meta: %v", err)
		}
		ChunkIds[chunkId] = index
		wg.Add(1)
		sem <- struct{}{}
		go func(chunkId uint64, chunkOffset int64) {
			defer func() {
				<-sem
				wg.Done()
			}()
			chunkBuffer := make([]byte, common.ChunkSize)
			bytes, err := file.ReadAt(chunkBuffer, chunkOffset)
			if err != nil {
				if err == io.EOF {
					//Fill the data with a size less than 32MB.
					if bytes < common.ChunkSize {
						for j := bytes; j < common.ChunkSize; j++ {
							chunkBuffer[j] = 0
						}
					}
				} else {
					log.Fatalf("read file error [%v], chunkOffset[%v]", err, chunkOffset)
					return
				}
			}
			var excludes []string
			for i := 0; i < ReplicaNumber; i++ {
				for {
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
					err = sdk.dataTransfer.Send(dnId, strconv.FormatUint(chunkId, 10), chunkBuffer)
					if err != nil {
						log.Errorf("Failed to send data: %v", err)
						return
					}
					// Get the chunk meta to verify whether this replica is successfully stored.
					isBreak := false
					for times := 0; times < 3; times++ {
						log.Debugf("ask for ChunkMeta for chunkId:[%v]", chunkId)
						metaResp, err := sdk.mnClient.GetChunkMeta(chunkId)
						if err != nil {
							log.Errorf("Failed to get chunk meta: %v", err)
							return
						}
						_, ok := metaResp.DnUuids[dnId]
						if ok {
							isBreak = true
							break
						}
						time.Sleep(1 * time.Second)
					}
					if isBreak {
						log.Debugf("Chunk [%v] has stored success", chunkId)
						break
					}
				}
			}
		}(chunkId, offset)
	}
	wg.Wait()
	//TODO:重试逻辑
	fileId, err := sdk.mnClient.AddFileMeta(FileName, fileSize, ChunkIds)
	if err != nil {
		log.Fatalf("Failed to add file meta: %v", err)
	}
	sdk.filesMapMux.Lock()
	sdk.filesMap[FileName] = fileId
	sdk.filesMapMux.Unlock()
	return nil
}

func (sdk *Sdk) ReadFile(FileName string, concurrencyNumber int) error {
	sdk.filesMapMux.Lock()
	fileId, ok := sdk.filesMap[FileName]
	if !ok {
		return fmt.Errorf("file not exist")
	}
	sdk.filesMapMux.Unlock()
	fileMeta, err := sdk.mnClient.GetFileMeta(fileId)
	if err != nil {
		return fmt.Errorf("get file meta error [%v]", err)
	}
	filePath := filepath.Join(sdk.filePath, FileName)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("open file error [%v]", err)
	}
	sem := make(chan struct{}, concurrencyNumber)
	wg := sync.WaitGroup{}
	for key, value := range fileMeta.ChunkIds {
		wg.Add(1)
		sem <- struct{}{}
		go func(chunkId uint64, chunkOffset int64) {
			defer func() {
				<-sem
				wg.Done()
			}()
			chunkMeta, err := sdk.mnClient.GetChunkMeta(chunkId)
			if err != nil {
				log.Errorf("Failed to get chunk meta: %v", err)
			}
			readChunkOk := false
			for dnId := range chunkMeta.DnUuids {
				err = sdk.mnClient.ReadChunkFromDN(sdk.UUID, dnId, chunkId)
				if err != nil {
					log.Fatalf("Failed to add chunk on dn: %v", err)
				}
				ch := make(chan bool, 1)
				err := sdk.dataTransfer.StartToReceive(dnId, func(targetUUID string, label string, data []byte) {
					log.Info("Receive data !")
					bytes, err := file.WriteAt(data, chunkOffset*common.ChunkSize)
					if err != nil {
						ch <- false
						log.Fatalf("WriteAt error [%v]", err)
					}
					if bytes < common.ChunkSize {
						ch <- false
						log.Errorf("WriteAt length error [%v]", err)
					}
					ch <- true
				})
				if err != nil {
					ch <- false
					log.Fatalf("Failed to receive data: %v", err)
				}
				if <-ch {
					readChunkOk = true
					break
				}
			}
			if !readChunkOk {
				log.Fatalf("Read chunk fail")
			}
		}(key, value)
	}
	wg.Wait()
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat file error [%v]", err)
	}
	currentSize := fileInfo.Size()
	if currentSize < fileMeta.FileSize {
		return fmt.Errorf("file size [%v] less than file size [%v]", currentSize, fileMeta.FileSize)
	} else if currentSize > fileMeta.FileSize {
		err = file.Truncate(fileMeta.FileSize)
		if err != nil {
			return fmt.Errorf("truncate file error [%v]", err)
		}
		log.Debugf("Truncated file size [%v]", fileMeta.FileSize)
	}
	return nil
}

func (sdk *Sdk) DeleteFile(FileName string) error {
	return nil
}
