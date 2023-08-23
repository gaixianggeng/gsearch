package api

import (
	"gsearch/conf"
	"gsearch/internal/meta"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminController 管理接口
type AdminController struct {
}

// Summary 汇总信息
func (a *AdminController) Summary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

func NewAdmin(profile *meta.Profile, c *conf.Config) *AdminController {
	return &AdminController{}
}
