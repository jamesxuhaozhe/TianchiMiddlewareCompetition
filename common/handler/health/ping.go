package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthCheck shows `OK` as the ping-pong result
func Check(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
}
