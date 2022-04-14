package tests

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestWrite(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c, _ := os.ReadFile("../data/forward.db")
	log.Debug(len(c))
	buf := bytes.NewBuffer(c)
	a := make([]uint64, 21)
	binary.Read(buf, binary.LittleEndian, &a)
	log.Debug(a)
}
