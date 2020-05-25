package backendprocess

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SetBadTraceIdsReq struct {
	Ids      []string `json:"ids"`
	BatchPos int      `json:"batchPos"`
}

func SetBadTraceIds(c *gin.Context) {
	var req SetBadTraceIdsReq
	if err := c.Bind(&req); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	fmt.Printf("%v\n", req)
	for _, id := range req.Ids {
		fmt.Printf("id: %s ", id)
	}
	fmt.Println()
	c.String(http.StatusOK, "suc")
}
