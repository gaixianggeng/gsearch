package api

import (
	"doraemon/conf"
	"doraemon/internal/engine"
	"doraemon/internal/recall"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recall 召回
type Recall struct {
	*recall.Recall
}

// Get 搜索入口
func (r *Recall) Get(c *gin.Context) {
	query := c.Query("query") // shortcut for c.Request.URL.Query().Get("lastname")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is empty"})
		return
	}
	res, err := r.Search(query)
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

// NewRecallServ 创建召回服务
func NewRecallServ(meta *engine.Meta, c *conf.Config) *Recall {
	eng := engine.NewEngine(meta, c, engine.SearchMode)
	r := recall.NewRecall(eng)
	return &Recall{r}

}
