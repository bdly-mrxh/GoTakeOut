package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"takeout/common/constant"
	"takeout/common/errs"
)

// GetId 获取当前账户ID，可以直接向上传递错误（自定义）
func GetId(ctx *gin.Context) (int, error) {
	idAny, exists := ctx.Get(constant.ID)
	if !exists {
		return 0, errs.New(constant.CodeUnauthorized, constant.MsgGetAccountInfoFail)
	}
	id, err := strconv.Atoi(idAny.(string))
	if err != nil {
		return 0, errs.Wrap(err, constant.CodeUnauthorized, constant.MsgGetAccountInfoFail)
	}
	return id, nil
}
