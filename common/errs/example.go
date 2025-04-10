package errs

import (
	"errors"
	"takeout/common/constant"
)

// 示例：如何使用错误处理系统

// ExampleCreateError 创建一个简单的错误
func ExampleCreateError() error {
	return New(constant.CodeBadRequest, constant.MsgBadRequest)
}

// ExampleWrapError 包装一个已有的错误
func ExampleWrapError() error {
	originalErr := errors.New("数据库连接失败")
	return Wrap(originalErr, constant.CodeServerError, constant.MsgServerError)
}

// ExampleGetMessageByCode 使用错误码获取消息
func ExampleGetMessageByCode(code int) string {
	switch code {
	case constant.CodeSuccess:
		return constant.MsgSuccess
	case constant.CodeBadRequest:
		return constant.MsgBadRequest
	case constant.CodeUnauthorized:
		return constant.MsgUnauthorized
	case constant.CodeServerError:
		return constant.MsgServerError
	default:
		return "未知错误"
	}
}
