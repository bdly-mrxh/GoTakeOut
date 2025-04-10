package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
	"takeout/model/dto"
)

// WeChatUserController 微信用户接口
type WeChatUserController struct {
	userService service.UserService
}

func NewWeChatUserController() *WeChatUserController {
	return &WeChatUserController{}
}

// Login 实现微信登录接口
func (c *WeChatUserController) Login(ctx *gin.Context) {
	var userLoginDTO dto.UserLoginDTO
	if err := ctx.ShouldBindJSON(&userLoginDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	userLoginVO, err := c.userService.Login(&userLoginDTO)
	if err != nil {
		logger.Error(constant.MsgUserLoginFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUserLoginSuccess, userLoginVO)
}

// Logout 实现用户退出
func (c *WeChatUserController) Logout(ctx *gin.Context) {
	response.Success(ctx, constant.MsgUserLogoutSuccess, nil)
}
