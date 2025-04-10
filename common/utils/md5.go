package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(str string) string {
	// 创建一个MD5哈希对象
	h := md5.New()
	// 计算字符串的MD5哈希值
	h.Write([]byte(str))
	sum := h.Sum(nil)
	// 将哈希值转为字符串返回
	return hex.EncodeToString(sum)
}
