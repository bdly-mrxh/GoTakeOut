package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
)

// ShopController 店铺控制器
type ShopController struct {
	shopService service.ShopService
}

// NewShopController 创建店铺控制器
func NewShopController() *ShopController {
	return &ShopController{}
}

// SetStatus 设置店铺状态
func (c *ShopController) SetStatus(ctx *gin.Context) {
	statusStr := ctx.Param("status")
	if statusStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err = c.shopService.SetStatus(status)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}

// GetStatus 获取店铺状态
func (c *ShopController) GetStatus(ctx *gin.Context) {
	status, err := c.shopService.GetStatus()
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, status)
}
