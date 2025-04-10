package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
)

// CategoryController 用户端分类接口
type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

// List 查询分类
func (c *CategoryController) List(ctx *gin.Context) {
	typeID := struct {
		ID int `form:"type"`
	}{}
	if err := ctx.ShouldBindQuery(&typeID); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	list, err := c.categoryService.List(typeID.ID)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, list)
}
