package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/conf"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/clientprocess/engine"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils"
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
func SetParameter(c *gin.Context) {
	var r SetParameterReq
	if err := c.Bind(&r); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	log.Printf("data source port is %s\n", r.DatasourcePort)

	// set the datasource port
	conf.SetDatasourcePort(r.DatasourcePort)

	if utils.IsClientProcess() {
		fmt.Println("Client process starts processing data...")
		go engine.ProcessData()
	}
	c.String(http.StatusOK, "suc")
}
