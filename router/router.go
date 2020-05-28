package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler"
	backend "github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/backendprocess"
	client "github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/clientprocess"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(mw...)

	g.HEAD("/ready", handler.Ready)
	g.GET("/setParameter", handler.SetParameter)

	if utils.IsBackendProcess() {
		g.POST("/setBadTraceIds", backend.SetBadTraceIds)
		g.GET("/markFinish", backend.MarkFinish)
	}

	if utils.IsClientProcess() {
		g.POST("/getSpansForBadTraceIds", client.GetSpansForBadTraceId)
	}

	hg := g.Group("/check")
	{
		hg.GET("/health", handler.Check)
	}
	return g
}
