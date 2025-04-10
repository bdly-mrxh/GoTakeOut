package response

import (
	"net/http"
	"takeout/common/constant"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int    `json:"code"` // 业务状态码
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

// Success 成功响应
func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Code: constant.CodeSuccess, // 1 表示成功
		Msg:  message,
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}

// BadRequest 无效请求响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code: 400,
		Msg:  message,
		Data: nil,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: 401,
		Msg:  message,
		Data: nil,
	})
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: 500,
		Msg:  message,
		Data: nil,
	})
}
