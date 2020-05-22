package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jamesxuhaozhe/tianchimiddlewarecompetition/handler"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(mw...)

	g.HEAD("/ready", handler.Ready)
	g.GET("/setParameter", handler.SetParameter)

	hg := g.Group("/check")
	{
		hg.GET("/health", handler.Check)
	}
	return g
}
