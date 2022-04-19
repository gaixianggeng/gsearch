package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BinaryWrite any -> bytes.Buffer
func BinaryWrite(buf *bytes.Buffer, v any) error {
	size := binary.Size(v)
	// log.Debug("docid size:", size)
	if size <= 0 {
		return fmt.Errorf("encodePostings binary.Size err,size: %v", size)
	}
	return binary.Write(buf, binary.LittleEndian, v)
}

// BinaryRead bytes.Buffer -> any
func BinaryRead(buf *bytes.Buffer, v any) error {
	return binary.Read(buf, binary.LittleEndian, &v)
}
