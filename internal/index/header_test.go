package index

import (
	"reflect"
	"testing"
)

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name string
		want *Header
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseHeader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
