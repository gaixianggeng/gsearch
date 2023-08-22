package index

import (
	"gsearch/conf"
	"gsearch/internal/meta"
	"gsearch/pkg/utils/log"
)

// Run 索引写入入口
func Run(meta *meta.Profile, conf *conf.Config) {

	log.Infof("index run...")
	index, err := NewIndexEngine(meta, conf)
	if err != nil {
		panic(err)
	}
	defer index.Close()
	index.IndexDoc()
	log.Infof("index run end")
}
