package segment

import (
	"doraemon/conf"
	"log"
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
		Engine *Segment
	}
	type args struct {
		bufInvertHash InvertedIndexHash
		token         string
		position      uint64
		docID         uint64
	}

	e := newEng(IndexMode)
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
			docCount := tt.args.bufInvertHash[tt.args.token].DocCount
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

func TestEngineFetchPostings(t *testing.T) {

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *PostingsList
		want1   uint64
		wantErr bool
	}{
		{
			name: "test1", args: args{"据数"}, want: nil, want1: 1, wantErr: false,
		},
		{
			name: "test2", args: args{"数据"}, want: nil, want1: 2, wantErr: false,
		},
		{
			name: "test3", args: args{"北京"}, want: nil, want1: 2, wantErr: false,
		},
		{
			name: "test4", args: args{"道口"}, want: nil, want1: 2, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eng := newEng(SearchMode)
			got, got1, err := eng.FetchPostings(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Engine.FetchPostings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for got != nil {
				t.Logf("got:%v", got)
				got = got.Next
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Engine.FetchPostings() got = %v, want %v", got, tt.want)
			// }
			if got1 != tt.want1 {
				t.Errorf("Engine.FetchPostings() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func newEng(mode Mode) *Segment {
	c, err := conf.ReadConf("../../conf/conf.toml")
	if err != nil {
		log.Fatal(err)
	}
	eng := NewSegment(0, c)
	return eng

}
