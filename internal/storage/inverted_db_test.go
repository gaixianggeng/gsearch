package storage

import (
	"reflect"
	"testing"
)

const (
	termDB     = "../../data/term.db"
	invertedDB = "../../data/inverted.db"
	forwardDB  = "../../data/forward.db"
	sourceFile = "../../data/source.csv"
)

func TestInvertedDB_GetTermInfo(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *TermInfo
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", args{"数据"}, &TermInfo{2, 0, 64}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewInvertedDB(termDB, invertedDB)
			got, err := tr.GetTermInfo(tt.args.token)
			t.Logf("got:%+v,err:%+v", got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("InvertedDB.GetTermInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvertedDB.GetTermInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
