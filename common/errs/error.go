package errs

import (
	"errors"
	"fmt"
)

// Error 自定义错误类型
type Error struct {
	Code    int    // 错误码
	Message string // 错误消息
	Err     error  // 原始错误
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("错误码: %d, 错误信息: %s, 原始错误: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Message)
}

// New 创建一个新的错误
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     nil,
	}
}

// Wrap 包装一个已有的错误
func Wrap(err error, code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// GetCode 获取错误码
func GetCode(err error) int {
	if err == nil {
		return 0
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return -1 // 未知错误码
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}
	return err.Error()
}
