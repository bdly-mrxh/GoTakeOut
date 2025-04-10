package response

import (
	"net/http"
	"takeout/common/constant"
	"takeout/common/errs"

	"github.com/gin-gonic/gin"
)

// 定义业务错误码到 HTTP 状态码的映射表
var errorCodeToHttpStatus = map[int]int{
	constant.CodeBadRequest:    http.StatusBadRequest,          // 400 参数有误
	constant.CodeDatabaseError: http.StatusInternalServerError, // 500 数据库查询失败
	constant.CodeNotFound:      http.StatusNotFound,            // 404 用户未找到
	constant.CodeUnauthorized:  http.StatusUnauthorized,        // 401 用户未登录
	constant.CodeForbidden:     http.StatusForbidden,           // 403 无权限访问
	constant.CodeServerError:   http.StatusInternalServerError, // 500 服务器内部错误
	constant.CodeInternalError: http.StatusInternalServerError, // 系统内部错误
	constant.CodeCacheError:    http.StatusInternalServerError,

	constant.CodePasswordError: http.StatusUnauthorized,

	constant.CodeEmployLoginFail:       http.StatusUnauthorized,
	constant.CodeEmployeeCreateFail:    http.StatusBadRequest,
	constant.CodeEmployeeUpdateFail:    http.StatusBadRequest,
	constant.CodeEmployeePageQueryFail: http.StatusBadRequest,

	constant.CodeCategoryCreateFail: http.StatusBadRequest,
	constant.CodeDeleteCategoryFail: http.StatusBadRequest,

	constant.CodeBusinessError: http.StatusOK, // 对于业务错误，返回正常状态
}

// ErrorResponse 处理错误响应
func ErrorResponse(c *gin.Context, err error) {
	code := errs.GetCode(err)
	message := errs.GetMessage(err)

	// 根据错误码映射到HTTP状态码

	httpStatus, ok := errorCodeToHttpStatus[code]
	if !ok {
		// 未知错误码，默认服务器内部错误
		c.JSON(http.StatusInternalServerError, Response{
			Code: code,
			Msg:  message,
			Data: nil,
		})
		return
	}

	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}

// ErrorWithData 处理带有数据的错误响应
func ErrorWithData(c *gin.Context, err error, data any) {
	code := errs.GetCode(err)
	message := errs.GetMessage(err)

	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
}
