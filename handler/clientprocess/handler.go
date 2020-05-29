package clientprocess

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/clientprocess/engine"
	"net/http"
)

type getSpansForBadTraceIdsReq struct {
	Ids      []string `json:"ids"`
	BatchPos int      `json:"batchPos"`
}

type getSpansForBadTraceIdsResp struct {
	Map map[string]*[]string `json:"map"`
}

// GetSpansForBadTraceId takes request from the backend process server and return a map.
// map's key is the bad trace id and the key is a list of spans under the key trace id
func GetSpansForBadTraceId(c *gin.Context) {
	var req getSpansForBadTraceIdsReq
	if err := c.Bind(&req); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	result, err := engine.GetSpansForBadTraceId(req.Ids, req.BatchPos)
	if err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	resp := getSpansForBadTraceIdsResp{Map: result}
	c.JSON(http.StatusOK, resp)
}
