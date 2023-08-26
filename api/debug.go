package api

import (
	"gsearch/conf"
	"gsearch/internal/meta"
	"gsearch/pkg/utils/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DebugController 管理接口
type DebugController struct {
}

// List 汇总信息
func (a *DebugController) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

func (a *DebugController) Doc(c *gin.Context) {
	docID := c.Param("docID")
	log.Debugf("doc id:%v", docID)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

// Term
func (a *DebugController) Term(c *gin.Context) {
	term := c.Param("term")
	log.Debugf("term:%v", term)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

func NewDebug(profile *meta.Profile, c *conf.Config) *DebugController {
	return &DebugController{}
}
