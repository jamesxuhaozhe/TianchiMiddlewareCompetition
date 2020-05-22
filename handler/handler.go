package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"log"
	"net/http"
)

// HealthCheck shows `OK` as the ping-pong result.
func Check(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
}

// Ready signals the remote data source that the server is ready.
func Ready(c *gin.Context) {
	c.String(http.StatusOK, "suc")
}

// setParameter notifies server what the remote data source port is.
func SetParameter(c *gin.Context)  {
	var r SetParameterReq
	if err := c.Bind(&r); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	log.Printf("data source port is %s", r.Dataport)

	// set the datasource port
	conf.SetDatasourcePort(r.Dataport)
	c.String(http.StatusOK, "suc")
}
