package conf

import (
	"reflect"
	"testing"
)

func TestReadConf(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", &Config{
			Project: "gsearch",
			Version: "0.0.1",
			Storage: struct {
				Path string "toml:\"path\""
			}{Path: "../data/"},
			Source: struct {
				Files []string "toml:\"files\""
			}{
				Files: []string{"../data/source.csv"}}},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConf("../conf/conf.toml")
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
