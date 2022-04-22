package storage

import (
	"testing"

	"github.com/boltdb/bolt"
)

func TestForwardDB_Count(t *testing.T) {
	db, err := bolt.Open("../../data/forward.db", 0600, nil)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ForwardDB{
				db: db,
			}
			got, err := f.Count()
			if (err != nil) != tt.wantErr {
				t.Errorf("ForwardDB.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ForwardDB.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}
