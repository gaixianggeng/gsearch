package index

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

func Test_encodePostings(t *testing.T) {
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
		// TODO: Add test cases.
		{
			"test1", args{
				&PostingsList{
					DocID:         234,
					positions:     []uint64{8, 9},
					positionCount: 2},
				0},
			[]uint64{234, 2, 8, 9},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodePostings(tt.args.postings, tt.args.postingsLen)

			a := make([]uint64, 4)
			binary.Read(got, binary.BigEndian, &a)
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
	p.positionCount = 1
	p.DocID = 234
	p.positions = []uint64{8, 9}
	got, err := encodePostings(p, 0)
	if err != nil {
		t.Errorf("encodePostings() error = %v, wantErr %v", err, nil)
		return
	}

	fmt.Println(got.Len())
	fmt.Println(string(got.Bytes()))

	a := make([]uint64, 4)
	binary.Read(got, binary.BigEndian, &a)
	fmt.Println(a)
}
