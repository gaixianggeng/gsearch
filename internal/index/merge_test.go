package index

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"doraemon/pkg/utils"
	"encoding/json"
	"testing"
	"time"
)

func TestMergeScheduler_mayMerge(t *testing.T) {
	confPath := "../../conf/conf.toml"
	c, err := conf.ReadConf(confPath)
	if err != nil {
		t.Error(err)
		return
	}
	c.Storage.Path = "../../data/"

	meta, err := engine.ParseMeta(c)
	if err != nil {
		t.Error(err)
		return
	}
	if meta == nil {
		t.Error("meta is nil")
		return
	}
	cont, _ := json.Marshal(meta.SegInfo)
	t.Logf("seg info:%s:", cont)

	m := NewScheduleer(meta, c)

	go m.Merge()

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"test1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.mayMerge()
		})
	}
	time.Sleep(3e9)
}

func init() {
	utils.LogInit()
}
