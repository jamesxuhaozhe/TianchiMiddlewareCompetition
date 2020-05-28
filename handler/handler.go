package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/clientprocess/engine"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/log"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils"
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
func SetParameter(c *gin.Context) {
	port := c.Query("port")
	if port == "" {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	log.Infof("data source port is %s", port)

	// set the datasource port
	conf.SetDatasourcePort(port)

	if utils.IsClientProcess() {
		log.Info("Client process starts processing data...")
		go engine.ProcessData()
	}
	c.String(http.StatusOK, "suc")
}
