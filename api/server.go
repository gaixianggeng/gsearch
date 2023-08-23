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
	adminAPI := NewAdmin(meta, conf)

	r := gin.Default()
	r.GET("/search", recallAPI.Search)
	// admin 相关的接口
	admin := r.Group("/admin")
	{
		admin.GET("/summary", adminAPI.Summary)
	}

	r.Run(":5168")
}
