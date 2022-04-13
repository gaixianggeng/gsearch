package main

import log "github.com/sirupsen/logrus"

func main() {
	log.Debug("hello world!")
}

func init() {
	log.SetLevel(log.DebugLevel)
}
