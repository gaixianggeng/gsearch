package meta

import (
	"fmt"
	"gsearch/conf"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"
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
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File)
			return frame.Function, fmt.Sprintf("%v:%d", fileName, frame.Line)
		},
	})
}
