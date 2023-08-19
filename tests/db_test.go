package tests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gsearch/internal/storage"
	"os"
	"sort"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
)

func TestReadDB(t *testing.T) {
	c, _ := os.ReadFile("../data/inverted.db")
	log.Debug(len(c))
	buf := bytes.NewBuffer(c)
	a := make([]uint64, len(c)/8)
	binary.Read(buf, binary.LittleEndian, &a)
	log.Debug(a)
}
func TestBucketGetFromNode(t *testing.T) {

	termName := "../data/term.db"
	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// for i := 0; i < 100; i++ {
	// 	err = storage.Put(db, "term",
	// 		[]byte(fmt.Sprintf("test%d", i)),
	// 		[]byte(fmt.Sprintf("%d", i)))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	b, err := storage.Get(db, "term", []byte("test1"))
	if err != nil {
		log.Fatal(err)
	}
	buf := bytes.NewBuffer(b)
	a := make([]byte, len(b))
	binary.Read(buf, binary.LittleEndian, &a)
	log.Debugf("%s", a)
}

func TestSort(t *testing.T) {
	key := "a"
	nodes := []string{"b", "d", "e", "f"}

	exact := false
	index := sort.Search(len(nodes), func(i int) bool {
		// TODO(benbjohnson): Optimize this range search. It's a bit hacky right now.
		// sort.Search() finds the lowest index where f() != -1 but we need the highest index.
		ret := bytes.Compare([]byte(nodes[i]), []byte(key))
		if ret == 0 {
			exact = true
		}
		return ret != -1
	})

	fmt.Printf("sort:%d, %v\n", index, exact)

}

func TestCom(t *testing.T) {

	ret := bytes.Compare([]byte("a"), []byte("b"))
	fmt.Println(ret)
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
