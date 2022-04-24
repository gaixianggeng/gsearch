package main

import (
	"doraemon/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("start")
	run()
	log.Info("end")
}

func init() {
	utils.LogInit()
}
