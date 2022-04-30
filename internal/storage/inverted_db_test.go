package storage

import (
	"reflect"
	"testing"
)

const (
	termDB     = "../../data/0.term"
	invertedDB = "../../data/0.inverted"
	forwardDB  = "../../data/0.forward"
	sourceFile = "../../source/source.csv"
)

func TestInvertedDB_GetTermInfo(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *TermValue
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", args{"ab"}, &TermValue{3, 0, 80}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewInvertedDB(termDB, invertedDB)
			got, err := tr.GetTermInfo(tt.args.token)
			t.Logf("got:%d,%+v,err:%+v", got.DocCount, got, err)
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
