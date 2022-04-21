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
			DocsCount: 9,
		},
		"成都": &engine.InvertedIndexValue{
			DocsCount: 4,
		},
		"上海": &engine.InvertedIndexValue{
			DocsCount: 6,
		},
		"深圳": &engine.InvertedIndexValue{
			DocsCount: 0,
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
				t.Logf("token:%s,count:%d", v.token, v.invertedIndex.DocsCount)
			}

		})
	}
}
func init() {
	log.SetLevel(log.DebugLevel)
}
