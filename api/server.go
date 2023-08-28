package api

import (
	"gsearch/conf"
	"gsearch/internal/meta"
	"gsearch/pkg/utils/log"

	"github.com/gin-gonic/gin"
)

// Start 启动服务
func Start(meta *meta.Profile, conf *conf.Config) {
	log.Info("start")
	recallAPI := NewRecall(meta, conf)
	debugAPI := NewDebug(meta, conf)

	r := gin.Default()
	r.GET("/search", recallAPI.Search)
	// admin 相关的接口
	admin := r.Group("/debug")
	{
		admin.GET("/list", debugAPI.List)
		admin.GET("/doc/:docID", debugAPI.Doc)
		admin.GET("/index/:term", debugAPI.Term)
	}

	r.Run(":5168")
}
