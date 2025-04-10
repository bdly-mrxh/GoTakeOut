package utils

import (
	"github.com/ulule/deepcopier"
)

// CopyProperties 使用deepcopier库复制对象属性
// from: 源对象
// to: 目标对象
// 返回错误信息
func CopyProperties(from any, to any) error {
	return deepcopier.Copy(from).To(to)
}
