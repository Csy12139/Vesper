package MN

import (
	"errors"
	"fmt"

	"github.com/Csy12139/Vesper/common"
	"github.com/Csy12139/Vesper/log"
	pb "github.com/Csy12139/Vesper/proto"
	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

type KVTable struct {
	db               *badger.DB
	ChunkIdSeq       *badger.Sequence
	chunkTablePrefix []byte
	fileTablePrefix  []byte
	chunkIdPrefix    []byte
}

func NewKVTable() *KVTable {
	return &KVTable{
		db:               nil,
		chunkTablePrefix: []byte("chunk_"),
		fileTablePrefix:  []byte("file_"),
		chunkIdPrefix:    []byte("ChunkId_"),
	}
}

func (kv *KVTable) Open(path string) error {
	var err error
	kv.db, err = badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return err
	}
	kv.ChunkIdSeq, err = kv.db.GetSequence(kv.chunkTablePrefix, 16)
	if err != nil {
		return err
	}
	return nil
}

func (kv *KVTable) Close() {
	err := kv.db.Close()
	if err != nil {
		log.Errorf("failed to close DB: %s", err)
	}
}

func (kv *KVTable) PutChunkMeta(meta *common.ChunkMeta) error {
	key := kv.chunkTablePrefix
	key = append(key, common.Uint64ToBytes(meta.ID)...)
	value, err := proto.Marshal(common.ChunkMeta2Proto(meta))
	if err != nil {
		return fmt.Errorf("marshal chunk meta failed [%s]", err)
	}
	err = kv.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (kv *KVTable) GetChunkMeta(ChunkId uint64) (*common.ChunkMeta, error) {
	key := kv.chunkTablePrefix
	key = append(key, common.Uint64ToBytes(ChunkId)...)
	var meta *common.ChunkMeta
	err := kv.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			var chunkMetaProto pb.ChunkMeta
			err := proto.Unmarshal(val, &chunkMetaProto)
			if err != nil {
				return err
			}
			meta = common.Proto2ChunkMeta(&chunkMetaProto)
			return nil
		})
		return err
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, common.ErrChunkNotFound
		}
		return nil, err
	}
	return meta, nil
}

func (kv *KVTable) AllocateChunkId() (uint64, error) {
	id, err := kv.ChunkIdSeq.Next()
	if err != nil {
		return 0, fmt.Errorf("allocate chunk id failed [%s]", err)
	}
	return id, nil
}
