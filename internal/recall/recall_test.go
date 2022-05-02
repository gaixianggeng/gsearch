package recall

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"doraemon/internal/segment"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
)

const (
	termDB     = "../../data/term.db"
	invertedDB = "../../data/inverted.db"
	forwardDB  = "../../data/forward.db"
	sourceFile = "../../data/source.csv"
)

func TestRecall_sortToken(t *testing.T) {

	type args struct {
		postHash segment.InvertedIndexHash
	}

	hash := segment.InvertedIndexHash{
		"北京": &segment.InvertedIndexValue{
			DocCount: 9,
		},
		"成都": &segment.InvertedIndexValue{
			DocCount: 4,
		},
		"上海": &segment.InvertedIndexValue{
			DocCount: 6,
		},
		"深圳": &segment.InvertedIndexValue{
			DocCount: 0,
		},
	}

	eng := newEng()
	defer eng.Close()

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "test1", args: args{postHash: hash}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recall{
				Engine: eng,
			}
			r.sortToken(tt.args.postHash)
			if r.queryToken == nil || len(r.queryToken) == 0 {
				t.Errorf("sortToken() error")
			}
			t.Log("after")
			for _, v := range r.queryToken {
				t.Logf("token:%s,count:%d", v.token, v.invertedIndex.DocCount)
			}

		})
	}
}
func TestRecall_Search(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchItem
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1", args: args{query: "五道口"}, want: nil, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("before")

			eng := newEng()
			r := NewRecall(eng)
			defer r.Close()

			got, err := r.Search(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Recall.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range got {
				t.Logf("docid:%v,score:%v", v.DocID, v.Score)
			}
		})
	}
}
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)
			return frame.Function, fmt.Sprintf("%v:%d", fileName, frame.Line)
		},
	})
}

func newEng() *engine.Engine {
	c, err := conf.ReadConf("../../conf/conf.toml")
	if err != nil {
		log.Fatal(err)
	}
	c.Storage.Path = "../../data/"

	meta, err := engine.ParseMeta(c)
	if err != nil {
		log.Fatal(err)
	}
	eng := engine.NewEngine(meta, c, segment.SearchMode)
	return eng

}
