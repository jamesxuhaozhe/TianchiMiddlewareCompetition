package backendprocess

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/backendprocess/engine"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/log"
	"net/http"
)

type setBadTraceIdsReq struct {
	Ids      []string `json:"ids"`
	BatchPos int      `json:"batchPos"`
}

func SetBadTraceIds(c *gin.Context) {
	var req setBadTraceIdsReq
	if err := c.Bind(&req); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	engine.SetBadTraceIds(req.Ids, req.BatchPos)
	c.String(http.StatusOK, "suc")
}

func MarkFinish(c *gin.Context) {
	engine.BumpProcessCount()
	c.String(http.StatusOK, "suc")
	log.Info("client marks finish.")
}
