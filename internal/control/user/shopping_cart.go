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

type ShoppingCartController struct {
	shoppingCartService service.ShoppingCartService
}

func NewShoppingCartController() *ShoppingCartController {
	return &ShoppingCartController{}
}

// Add 添加购物车
func (c *ShoppingCartController) Add(ctx *gin.Context) {
	var shoppingCartDTO dto.ShoppingCartDTO
	if err := ctx.ShouldBindJSON(&shoppingCartDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.shoppingCartService.Add(ctx, &shoppingCartDTO)
	if err != nil {
		logger.Error(constant.MsgShoppingCartAddFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgShoppingCartAddSuccess, nil)
}

// List 查看购物车
func (c *ShoppingCartController) List(ctx *gin.Context) {
	list, err := c.shoppingCartService.List(ctx)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, list)
}

// Clean 清空购物车
func (c *ShoppingCartController) Clean(ctx *gin.Context) {
	err := c.shoppingCartService.Clean(ctx)
	if err != nil {
		logger.Error(constant.MsgDeleteFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgDeleteSuccess, nil)
}

// Sub 减少购物车中的一个商品
func (c *ShoppingCartController) Sub(ctx *gin.Context) {
	var shoppingCartDTO dto.ShoppingCartDTO
	if err := ctx.ShouldBindJSON(&shoppingCartDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.shoppingCartService.Sub(ctx, &shoppingCartDTO)
	if err != nil {
		logger.Error(constant.MsgDeleteFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgDeleteSuccess, nil)
}
