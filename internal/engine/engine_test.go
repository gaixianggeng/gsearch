package engine

import (
	"testing"
)

const (
	termDB     = "../../data/term.db"
	invertedDB = "../../data/inverted.db"
	forwardDB  = "../../data/forward.db"
	sourceFile = "../../data/source.csv"
)

func TestIndex_token2PostingsLists(t *testing.T) {
	type fields struct {
		Engine *Engine
	}
	type args struct {
		bufInvertHash InvertedIndexHash
		token         string
		position      uint64
		docID         uint64
	}
	e := NewEngine(termDB, invertedDB, forwardDB)
	if e == nil {
		t.Errorf("new engine is nil")
	}
	defer e.Close()
	bufInvertedHash := make(InvertedIndexHash)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				Engine: e,
			},
			args: args{
				bufInvertHash: bufInvertedHash,
				token:         "北京",
				position:      0,
				docID:         123,
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				Engine: e,
			},
			args: args{
				bufInvertHash: bufInvertedHash,
				token:         "北京",
				position:      2,
				docID:         123,
			},
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				Engine: e,
			},
			args: args{
				bufInvertHash: bufInvertedHash,
				token:         "北京",
				position:      0,
				docID:         123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := e.Token2PostingsLists(tt.args.bufInvertHash, tt.args.token, tt.args.position, tt.args.docID); (err != nil) != tt.wantErr {
				t.Errorf("Index.token2PostingsLists() error = %v, wantErr %v", err, tt.wantErr)
			}
			count := tt.args.bufInvertHash[tt.args.token].PositionCount
			docCount := tt.args.bufInvertHash[tt.args.token].DocsCount
			if tt.name == "test1" && (count != 1 || docCount != 1) {
				t.Errorf("count:%v,docCount:%v", count, docCount)
			}
			if tt.name == "test2" && (count != 2 || docCount != 1) {
				t.Errorf("count:%v,docCount:%v", count, docCount)
			}
			if tt.name == "test3" && (count != 3 || docCount != 1) {
				t.Errorf("count:%v,docCount:%v", count, docCount)
			}
		})
	}
}
