package recall

import (
	"fmt"
	"gsearch/conf"
	"gsearch/internal/meta"
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

// func TestRecall_sortToken(t *testing.T) {

// 	type args struct {
// 		postHash segment.InvertedIndexHash
// 	}

// 	hash := segment.InvertedIndexHash{
// 		"北京": &segment.InvertedIndexValue{
// 			DocCount: 9,
// 		},
// 		"成都": &segment.InvertedIndexValue{
// 			DocCount: 4,
// 		},
// 		"上海": &segment.InvertedIndexValue{
// 			DocCount: 6,
// 		},
// 		"深圳": &segment.InvertedIndexValue{
// 			DocCount: 0,
// 		},
// 	}

// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 		{name: "test1", args: args{postHash: hash}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			r := newRecall()
// 			defer r.Close()

// 			tokens := r.sortToken(tt.args.postHash)
// 			if tokens == nil || len(tokens) == 0 {
// 				t.Errorf("sortToken() error")
// 			}
// 			t.Log("after")
// 			for _, v := range tokens {
// 				t.Logf("token:%s,count:%d", v.token, v.invertedIndex.DocCount)
// 			}

//			})
//		}
//	}
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
			r := newRecall()
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

func newRecall() *Recall {
	c, err := conf.ReadConf("../../conf/conf.toml")
	if err != nil {
		log.Fatal(err)
	}
	c.Storage.Path = "../../data/"

	meta, err := meta.ParseProfile(c)
	if err != nil {
		log.Fatal(err)
	}

	r := NewRecall(meta, c)
	return r
}
