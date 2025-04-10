package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"takeout/common/logger"
)

// RecoveryMiddleware 自定义全局异常恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer + recover 捕获 panic
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool

				// 判断是否是网络错误（Broken Pipe 或 Connection Reset）
				if e, ok := err.(error); ok {
					var netErr *net.OpError
					if errors.As(e, &netErr) { // 判断是否是网络错误
						var systemErr *os.SyscallError
						if errors.As(netErr.Err, &systemErr) { // 判断是否是系统调用错误
							errMsg := strings.ToLower(systemErr.Error())
							if strings.Contains(errMsg, "broken pipe") || strings.Contains(errMsg, "connection reset by peer") {
								brokenPipe = true
							}
						}
					}
				}

				// 获取 HTTP 请求数据（不包含 body）
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// Broken Pipe 错误：表示客户端在服务端响应之前主动关闭了连接。
				// 如果不处理，Gin 会触发 c.JSON 等操作，导致 二次 panic
				// 如果是 Broken Pipe 错误，特殊处理，防止多余日志污染
				if brokenPipe {
					logger.Error("Broken Pipe", zap.Any("error", err), zap.String("request", string(httpRequest)))
					// 直接返回错误，避免触发 Gin 的默认错误处理机制
					_ = c.Error(err.(error)) // 将错误传递给 Gin 的上下文对象，供其处理，但并不让 Gin 进入默认的错误处理流程。
					c.Abort()
					return
				}

				// 其它 panic 错误，记录堆栈信息
				logger.Error("[Recovery from panic]", zap.Any("error", err), zap.String("request", string(httpRequest)), zap.String("stack", string(debug.Stack())))

				// 返回 500 错误响应
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		// 执行下一个中间件
		c.Next()
	}
}
