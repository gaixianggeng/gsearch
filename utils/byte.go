package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BinaryWrite(buf *bytes.Buffer, v any) error {
	size := binary.Size(v)
	// log.Debug("docid size:", size)
	if size <= 0 {
		return fmt.Errorf("encodePostings binary.Size err,size: %v", size)
	}
	return binary.Write(buf, binary.LittleEndian, v)
}
