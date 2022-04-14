package tests

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	os.ReadFile("../data/forward.db")
}
