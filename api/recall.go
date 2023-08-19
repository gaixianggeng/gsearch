package api

import (
	"gsearch/conf"
	"gsearch/internal/engine"
	"gsearch/internal/recall"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecallController 召回
type RecallController struct {
	engine *recall.Recall
}

// Search 搜索入口
func (r *RecallController) Search(c *gin.Context) {
	query := c.Query("query") // shortcut for c.Request.URL.Query().Get("lastname")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is empty"})
		return
	}
	res, err := r.engine.Search(query)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"query": query,
		"data":  res,
	})

}

// NewRecall 创建召回服务
func NewRecall(meta *engine.Meta, c *conf.Config) *RecallController {
	r := recall.NewRecall(meta, c)
	return &RecallController{r}
}
