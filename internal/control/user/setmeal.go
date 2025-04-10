package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
)

// SetmealController 用户端套餐接口
type SetmealController struct {
	setmealService service.SetmealService
}

func NewSetmealController() *SetmealController {
	return &SetmealController{}
}

// List 根据分类ID查询套餐
func (c *SetmealController) List(ctx *gin.Context) {
	categoryIDStr := ctx.Query("categoryId")
	if categoryIDStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	list, err := c.setmealService.ListByCategoryID(categoryID)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, list)
}

// DishList 根据套餐ID查询所包含的菜品列表
func (c *SetmealController) DishList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	dishItems, err := c.setmealService.GetDishItemBySetmealID(id)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, dishItems)
}
