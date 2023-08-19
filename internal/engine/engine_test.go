package engine

import (
	"gsearch/conf"
	"gsearch/internal/segment"
	"testing"
)

func TestEngine_Text2PostingsLists(t *testing.T) {
	type fields struct {
		meta            *Meta
		conf            *conf.Config
		Scheduler       *MergeScheduler
		BufCount        uint64
		BufSize         uint64
		PostingsHashBuf segment.InvertedIndexHash
		CurrSegID       segment.SegID
		Seg             map[segment.SegID]*segment.Segment
		N               int32
	}
	type args struct {
		text  string
		docID uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				meta:            tt.fields.meta,
				conf:            tt.fields.conf,
				Scheduler:       tt.fields.Scheduler,
				BufCount:        tt.fields.BufCount,
				BufSize:         tt.fields.BufSize,
				PostingsHashBuf: tt.fields.PostingsHashBuf,
				CurrSegID:       tt.fields.CurrSegID,
				Seg:             tt.fields.Seg,
				N:               tt.fields.N,
			}
			if err := e.Text2PostingsLists(tt.args.text, tt.args.docID); (err != nil) != tt.wantErr {
				t.Errorf("Engine.Text2PostingsLists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
