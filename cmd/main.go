package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("start")
	run()
	log.Info("end")
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
