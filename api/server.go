package api

import (
	"doraemon/conf"
	"doraemon/internal/engine"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// StartServ 启动服务
func StartServ(meta *engine.Meta, conf *conf.Config) {
	log.Info("start")

	recall := NewRecallServ(meta, conf)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/search", recall.Get)

	r.Run(":5168")
}
