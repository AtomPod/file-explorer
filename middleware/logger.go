package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phantom-atom/file-explorer/internal/log"
)

//Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		log.Info(
			"time",
			time.Now().Format("2006/01/02 - 15:04:05"),
			"status_code", statusCode,
			"method", method,
			"latency", latency,
			"path", path,
			"remote_ip", clientIP,
			"body_size", bodySize,
			"error_msg", errorMsg,
		)
	}
}
