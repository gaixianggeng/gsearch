package main

import (
	"doraemon/conf"
	"doraemon/internal/index"
	"log"
)

func run() {
	confPath := "../conf/conf.toml"
	c, err := conf.ReadConf(confPath)
	if err != nil {
		log.Fatal(err)
	}
	index.Run(c)
}
