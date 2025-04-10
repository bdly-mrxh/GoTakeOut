package middleware

import (
	"go.uber.org/zap"
	"takeout/common/constant"
	"takeout/common/global"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/common/utils"

	"github.com/gin-gonic/gin"
)

// JwtAdmin 管理端认证
func JwtAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头 header（不是 query）获取 token
		token := ctx.GetHeader(global.Config.JWT.AdminTokenName)
		// 未携带 token 就是未登录状态
		if token == "" {
			logger.Error(constant.MsgJWTWithoutToken)
			response.Unauthorized(ctx, constant.MsgJWTWithoutToken)
			ctx.Abort()
			return
		}
		// jwt 校验
		id, err := utils.ParseToken(token, constant.EmpID)
		if err != nil {
			// token 解析失败
			logger.Error(constant.MsgJWTParseFail, zap.Error(err))
			response.ErrorResponse(ctx, err)
			ctx.Abort()
			return
		}
		// 存储当前员工 ID
		ctx.Set(constant.ID, id)
		ctx.Next()
	}
}

// JwtUser 用户端认证
func JwtUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(global.Config.JWT.UserTokenName)
		if token == "" {
			logger.Error(constant.MsgJWTWithoutToken)
			response.Unauthorized(ctx, constant.MsgJWTWithoutToken)
			ctx.Abort()
			return
		}
		id, err := utils.ParseToken(token, constant.UserID)
		if err != nil {
			logger.Error(constant.MsgJWTParseFail, zap.Error(err))
			response.ErrorResponse(ctx, err)
			ctx.Abort()
			return
		}
		ctx.Set(constant.ID, id)
		ctx.Next()
	}
}
