package tests

import (
	"os"
	"testing"
	"time"
)

// TestFileDel 测试文件读取快照
func TestFileDel(t *testing.T) {

	f, err := os.Open("test.txt")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(f.Fd())
	time.Sleep(10e9)
	s, _ := f.Stat()
	t.Log("size:", s.Size(), s.Size())
	b := make([]byte, s.Size(), s.Size())
	_, err = f.ReadAt(b, 0)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(b))
	t.Log(f.Fd())
}
