package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler/backendprocess"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/utils"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(mw...)

	g.HEAD("/ready", handler.Ready)
	g.GET("/setParameter", handler.SetParameter)

	if utils.IsBackendProcess() {
		g.POST("/setBadTraceIds", backendprocess.SetBadTraceIds)
	}

	hg := g.Group("/check")
	{
		hg.GET("/health", handler.Check)
	}
	return g
}
