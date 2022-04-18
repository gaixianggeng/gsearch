package tests

import (
	"brain/internal/storage"
	"bytes"
	"encoding/binary"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

func TestReadDB(t *testing.T) {
	c, _ := os.ReadFile("../data/forward.db")
	log.Debug(len(c))
	buf := bytes.NewBuffer(c)
	a := make([]uint64, len(c)/8)
	binary.Read(buf, binary.LittleEndian, &a)
	log.Debug(a)
}
func TestBucket_Get_FromNode(t *testing.T) {

	termName := "../data/term.db"
	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	b, err := storage.Get(db, "term", []byte("据数"))
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(b)
	a := make([]uint64, 2)
	binary.Read(buf, binary.LittleEndian, &a)
	log.Debugf("%v", a)
}

func TestGetForward(t *testing.T) {

	termName := "../data/forward.db"
	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	b, err := storage.Get(db, "forward", []byte("56291828"))
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(b)
	a := make([]byte, len(b))
	binary.Read(buf, binary.LittleEndian, &a)
	log.Debugf("%s", a)
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}
