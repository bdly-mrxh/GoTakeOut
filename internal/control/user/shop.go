package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
