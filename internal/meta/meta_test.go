package meta

import (
	"gsearch/conf"
	"reflect"
	"testing"
)

func TestParseHeader(t *testing.T) {

	tests := []struct {
		name    string
		want    *Profile
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c, err := conf.ReadConf("../../conf/conf.toml")
			if err != nil {
				t.Errorf("ReadConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			c.Storage.Path = "../../data/"
			got, err := ParseProfile(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
