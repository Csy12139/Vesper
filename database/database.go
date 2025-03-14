package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

type DataBase interface {
	Close() error
	Get(key []byte, opts *opt.ReadOptions) ([]byte, error)
	Put(key, value []byte, opts *opt.WriteOptions) error
	Delete(key []byte, opts *opt.WriteOptions) error
	Has(key []byte, opts *opt.ReadOptions) (bool, error)
	WriteBatch(batch *leveldb.Batch, opts *opt.WriteOptions) error
}
type DBHandler struct {
	db   *leveldb.DB
	lock sync.Mutex
}

func NewDB(path string, opts *opt.Options) (*DBHandler, error) {
	db, err := leveldb.OpenFile(path, opts)
	if err != nil {
		return nil, err
	}
	return &DBHandler{db: db}, nil
}

func (d *DBHandler) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.db.Close()
}

func (d *DBHandler) Get(key []byte, opts *opt.ReadOptions) ([]byte, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.db.Get(key, opts)
}

func (d *DBHandler) Put(key, value []byte, opts *opt.WriteOptions) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.db.Put(key, value, opts)
}
func (d *DBHandler) Delete(key []byte, opts *opt.WriteOptions) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.db.Delete(key, opts)
}
func (d *DBHandler) Has(key []byte, opts *opt.ReadOptions) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.db.Has(key, opts)
}
func (d *DBHandler) WriteBatch(batch *leveldb.Batch, opts *opt.WriteOptions) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.db.Write(batch, opts)
}
