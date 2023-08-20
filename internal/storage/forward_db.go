package storage

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
)

// ForwardDB 存储器
type ForwardDB struct {
	db *bolt.DB
}

// AddForward add forward data
func (f *ForwardDB) AddForward(doc *Document) error {
	key := strconv.Itoa(int(doc.DocID))
	body, _ := json.Marshal(doc)
	return Put(f.db, forwardBucket, []byte(key), body)
}

// PutForward add forward data
func (f *ForwardDB) PutForward(key, value []byte) error {
	return Put(f.db, forwardBucket, []byte(key), value)
}

// ForwardCount 获取文档总数
func (f *ForwardDB) ForwardCount() (uint64, error) {
	body, err := Get(f.db, forwardBucket, []byte(ForwardCountKey))
	if err != nil {
		return 0, err
	}
	c, err := strconv.Atoi(string(body))
	return uint64(c), err
}

// UpdateForwardCount 获取文档总数
func (f *ForwardDB) UpdateForwardCount(count uint64) error {
	return Put(f.db, forwardBucket, []byte(ForwardCountKey), []byte(strconv.Itoa(int(count))))
}

// GetForward get forward data
func (f *ForwardDB) GetForward(docID uint64) ([]byte, error) {
	key := strconv.Itoa(int(docID))
	return Get(f.db, forwardBucket, []byte(key))
}

// GetForwardCursor 获取遍历游标
func (f *ForwardDB) GetForwardCursor(termCh chan KvInfo) {
	f.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(forwardBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			termCh <- KvInfo{k, v}
		}
		close(termCh)
		return nil
	})
}

// Close --
func (f *ForwardDB) Close() {
	f.db.Close()
}

// NewForwardDB --
func NewForwardDB(dbName string) *ForwardDB {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &ForwardDB{db}
}
