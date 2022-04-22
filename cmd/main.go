package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	run()
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	// log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}
