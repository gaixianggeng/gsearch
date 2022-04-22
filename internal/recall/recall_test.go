package recall

import (
	"brain/internal/engine"
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
		postHash engine.InvertedIndexHash
	}

	hash := engine.InvertedIndexHash{
		"北京": &engine.InvertedIndexValue{
			DocCount: 9,
		},
		"成都": &engine.InvertedIndexValue{
			DocCount: 4,
		},
		"上海": &engine.InvertedIndexValue{
			DocCount: 6,
		},
		"深圳": &engine.InvertedIndexValue{
			DocCount: 0,
		},
	}

	eng := engine.NewEngine(termDB, invertedDB, forwardDB)
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
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
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
			name: "test1", args: args{query: "数据"}, want: nil, wantErr: false,
		},
		{
			name: "test2", args: args{query: "五道"}, want: nil, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("before")
			eng := engine.NewEngine(termDB, invertedDB, forwardDB)
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
