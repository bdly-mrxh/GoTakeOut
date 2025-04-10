package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
	"takeout/model/dto"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService service.CategoryService
}

// NewCategoryController 创建分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

// Create 创建分类
func (c *CategoryController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var createDTO dto.CategoryDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层创建分类
	if err := c.categoryService.Create(ctx, &createDTO); err != nil {
		logger.Error(constant.MsgCategoryCreateFail, zap.Error(err), zap.String("name", createDTO.Name))
		response.ErrorResponse(ctx, err)
		return
	}

	response.Success(ctx, constant.MsgCategoryCreateSuccess, nil)
}

// UpdateStatus 更新分类状态
func (c *CategoryController) UpdateStatus(ctx *gin.Context) {
	statusStr := ctx.Param("status")
	if statusStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	status, err := strconv.Atoi(statusStr)
	if err != nil || (status != 0 && status != 1) {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	idStr := ctx.Query("id")
	if idStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err = c.categoryService.UpdateStatus(id, status)
	if err != nil {
		logger.Error(constant.MsgCategoryStatusFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgCategoryStatusSuccess, nil)
}

// Update 更新分类
func (c *CategoryController) Update(ctx *gin.Context) {
	var updateDTO dto.CategoryDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.categoryService.Update(ctx, &updateDTO)
	if err != nil {
		logger.Error(constant.MsgCategoryUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgCategoryUpdateSuccess, nil)
}

// PageQuery 分类分页查询
func (c *CategoryController) PageQuery(ctx *gin.Context) {
	var queryDTO dto.CategoryPageDTO
	if err := ctx.ShouldBindQuery(&queryDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	page, err := c.categoryService.PageQuery(&queryDTO)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, page)
}

// List 根据类型查询分类
func (c *CategoryController) List(ctx *gin.Context) {
	typeIdStr := ctx.Query("type")
	typeId, err := strconv.Atoi(typeIdStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	list, err := c.categoryService.List(typeId)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, list)
}

// Delete 删除分类
func (c *CategoryController) Delete(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err = c.categoryService.Delete(id)
	if err != nil {
		logger.Error(constant.MsgDeleteFail, zap.Error(err), zap.Int("id", id))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgDeleteSuccess, nil)
}
