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

// DishController 客户端菜品接口
type DishController struct {
	dishService service.DishService
}

func NewDishController() *DishController {
	return &DishController{}
}

// List 根据分类ID查询菜品
func (c *DishController) List(ctx *gin.Context) {
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

	list, err := c.dishService.ListByCategoryID(categoryID)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, list)
}
