package tests

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

func TestWrite(t *testing.T) {
	c, _ := os.ReadFile("../data/forward.db")
	log.Debug(len(c))
	buf := bytes.NewBuffer(c)
	a := make([]uint64, 21)
	binary.Read(buf, binary.LittleEndian, &a)
	log.Debug(a)
}
func TestBucket_Get_FromNode(t *testing.T) {
	db := MustOpenDB()
	defer db.MustClose()
	log.Debug(db.Path())

	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("widgets"))
		if err != nil {
			t.Fatal(err)
		}
		if err := b.Put([]byte("foo"), []byte("bar")); err != nil {
			t.Fatal(err)
		}
		if v := b.Get([]byte("foo")); !bytes.Equal(v, []byte("bar")) {
			t.Fatalf("unexpected value: %v", v)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func init() {
	log.SetLevel(log.DebugLevel)

}
