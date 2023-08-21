package segment

import (
	"encoding/binary"
	"gsearch/pkg/utils/log"
	"reflect"
	"testing"
)

func TestEncodePostings(t *testing.T) {
	type args struct {
		postings    *PostingsList
		postingsLen uint64
	}

	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			"test1", args{
				&PostingsList{
					DocID:         234,
					Positions:     []uint64{8, 9},
					PositionCount: 2},
				0},
			[]uint64{0, 234, 2, 8, 9},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodePostings(tt.args.postings, tt.args.postingsLen)

			a := make([]uint64, 5)
			binary.Read(got, binary.LittleEndian, &a)
			log.Debug(a)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodePostings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(a, tt.want) {
				t.Errorf("encodePostings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPos(t *testing.T) {
	p := new(PostingsList)
	p.PositionCount = 1
	p.DocID = 234
	p.Positions = []uint64{8, 9}
	got, err := EncodePostings(p, 0)
	if err != nil {
		t.Errorf("encodePostings() error = %v, wantErr %v", err, nil)
		return
	}

	log.Debug(got.Len())
	log.Debug(string(got.Bytes()))

	a := make([]uint64, 4)
	binary.Read(got, binary.LittleEndian, &a)
	log.Debug(a)
}
