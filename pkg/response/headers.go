package response

import (
	"baseFrame/pkg/config"
	"github.com/gin-gonic/gin"
	"time"
)

func ResponseHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("app-time", time.Now().Format("2006-01-02 15:04:05"))
		c.Header("api-env", config.GetConfig("", "env"))
		c.Header("Access-Control-Expose-Headers", "app-time,api-env")
		c.Next()
	}
}
