package main

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"doraemon/internal/index"
	"fmt"
	"log"
	"time"
)

func run() {
	confPath := "../conf/conf.toml"
	c, err := conf.ReadConf(confPath)
	if err != nil {
		log.Fatal(err)
	}

	meta, err := engine.ParseMeta(c)
	if err != nil {
		panic(err)
	}
	if meta == nil {
		panic("meta is nil")
	}

	ticker := time.NewTicker(time.Second * 1)
	go meta.SyncByTicker(ticker)
	defer ticker.Stop()

	if c.Version != meta.Version {
		panic(fmt.Sprintf("version not match, %s != %s", c.Version, meta.Version))
	}

	index.Run(meta, c)

	// 最后同步元数据至文件
	meta.SyncMeta()

	time.Sleep(time.Second * 3)
}
