package middleware

import (
	"takeout/common/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 请求后
		statusCode := c.Writer.Status()
		logger.Info("HTTP Request",
			zap.String("path", path),
			zap.String("method", method),
			zap.Int("status", statusCode),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
