package api

import (
	"gsearch/conf"
	"gsearch/internal/meta"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Start 启动服务
func Start(meta *meta.Profile, conf *conf.Config) {
	log.Info("start")

	recall := NewRecall(meta, conf)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/search", recall.Search)

	r.Run(":5168")
}
