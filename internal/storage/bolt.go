package storage

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// Put 通过bolt写入数据
func Put(db *bolt.DB, bucket string, key []byte, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put(key, value)
		return err
	})
}

// Get 通过bolt获取数据
func Get(db *bolt.DB, bucket string, key []byte) ([]byte, error) {
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			// 单独测试 Text2PostingsLists 时会出现这个错误，因为没有添加正排数据
			return fmt.Errorf("bucket:%s not exist", bucket)
		}
		v = b.Get(key)
		if v == nil {
			return fmt.Errorf("key not found")
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("get token:%s err:%v", string(key), err)
	}
	return v, nil
}
