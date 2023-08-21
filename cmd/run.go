package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gsearch/api"
	"gsearch/conf"
	"gsearch/internal/index"
	"gsearch/internal/meta"
	"gsearch/pkg/utils/log"
	"time"
)

var (
	action int64
)

const (
	START_SERVER = 1
	START_INDEX  = 2
)

const (
	startCommand = `
Usage: gsearch [options]
Options:
	--flag 1:online service 2:index
`
)

func run() {
	if action != START_SERVER && action != START_INDEX {
		fmt.Printf("%s\n", startCommand)
		return
	}
	log.Info("start...")
	// TODO: 命令行启动参数
	confPath := "./conf/conf.toml"
	c, err := conf.ReadConf(confPath)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := json.Marshal(c)
	fmt.Printf("conf:%s\n", t)
	meta, err := meta.ParseProfile(c)
	if err != nil {
		panic(err)
	}
	if meta == nil {
		panic("meta is nil")
	}

	// 定时同步meta数据
	ticker := time.NewTicker(time.Second * 1)
	go meta.SyncByTicker(ticker)

	if c.Version != meta.Version {
		panic(fmt.Sprintf("version not match, %s != %s", c.Version, meta.Version))
	}

	start(c, meta)

	// close
	func() {
		// 最后同步元数据至文件
		log.Info("close")
		meta.SyncMeta()
		log.Info("close")
		ticker.Stop()
		log.Info("close")
	}()
}

func start(c *conf.Config, profile *meta.Profile) {
	if action == START_SERVER {
		log.Debugf("start server...")
		api.Start(profile, c)
	} else if action == START_INDEX {
		log.Debugf("start")
		index.Run(profile, c)
		log.Debug("end")
	}
}

// 获取flag参数
func flagInit() {
	flag.Int64Var(&action, "flag", 0, "start flag:\n[1:online service]\n [2:index]")
	flag.Parse()
}

func init() {
	flagInit()
}
