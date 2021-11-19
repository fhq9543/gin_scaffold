package logger

import (
	"baseFrame/pkg/util"
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

func Logrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("trace_id", util.RandUUID())
		startTime := time.Now()
		rawData, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		Infof("|%3d|%13v|%15s|%s|%s|%s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
			string(rawData),
			//c.MustGet("trace_id"),
		)
	}
}
