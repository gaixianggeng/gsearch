package engine

import (
	"encoding/json"
	"gsearch/conf"
	"gsearch/internal/meta"
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

	meta, err := meta.ParseProfile(c)
	if err != nil {
		t.Error(err)
		return
	}
	if meta == nil {
		t.Error("meta is nil")
		return
	}
	cont, _ := json.Marshal(meta.SegMeta.SegInfo)
	t.Logf("seg info:%s:", cont)

	m := NewScheduler(meta, c)

	go m.Merge()

	ticker := time.NewTicker(time.Second * 1)
	go meta.SyncByTicker(ticker)
	defer ticker.Stop()

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"test1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.MayMerge()
		})
	}
	time.Sleep(30e9)
}
