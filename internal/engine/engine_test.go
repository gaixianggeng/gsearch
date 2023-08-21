package engine

import (
	"encoding/json"
	"fmt"
	"gsearch/conf"
	"gsearch/internal/meta"
	"gsearch/internal/segment"
	"gsearch/pkg/utils/jstool"
	"testing"
)

func TestEngine_Text2PostingsLists(t *testing.T) {
	type args struct {
		text  string
		docID uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				text:  "北京的冬天，北京的夏天，北京的秋天，北京的春天春天",
				docID: 123,
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				text:  "北京的冬天，北京的夏天，北京的秋天，北京的春天",
				docID: 12,
			},
			wantErr: false,
		},
	}
	e := newEng()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := e.Text2PostingsLists(tt.args.text, tt.args.docID); (err != nil) != tt.wantErr {
				t.Errorf("Engine.Text2PostingsLists() error = %v, wantErr %v", err, tt.wantErr)
			}
			// t.Log("buf:", jstool.StructToStr(e.PostingsHashBuf))
		})
	}
	// 打印出 next 链路
	invert := e.PostingsHashBuf["北京"]
	t.Logf("posting list:%s", jstool.StructToStr(invert.PostingsList))
}

func newEng() *Engine {
	confPath := "../../conf/conf.toml"
	c, err := conf.ReadConf(confPath)
	if err != nil {
		panic(err)
	}
	t, _ := json.Marshal(c)
	fmt.Printf("conf:%s\n", t)
	c.Storage.Path = "../../data/"
	meta, err := meta.ParseProfile(c)
	if err != nil {
		panic(err)
	}
	return NewEngine(meta, c, segment.IndexMode)
}
