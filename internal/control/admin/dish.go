package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
	"takeout/model/dto"
)

// DishController 菜品控制器
type DishController struct {
	dishService service.DishService
}

// NewDishController 创建菜品控制器
func NewDishController() *DishController {
	return &DishController{}
}

// Create 创建菜品
func (c *DishController) Create(ctx *gin.Context) {
	var createDTO dto.DishDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	if err := c.dishService.CreateWithFlavors(ctx, &createDTO); err != nil {
		logger.Error(constant.MsgCreateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgCreateSuccess, nil)
}

// PageQuery 菜品分页查询
func (c *DishController) PageQuery(ctx *gin.Context) {
	var queryDTO dto.DishPageQueryDTO
	// ShouldBindQuery 无法判别有没有传值
	if err := ctx.ShouldBindQuery(&queryDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	// status 为 0 时是有意义的，要单独判断有没有传这个值
	statusStr := ctx.Query("status") // 传的是空，true
	if statusStr == "" {
		queryDTO.Status = constant.InvalidStatus
	}

	logger.Info("Info", zap.Any("DTO", queryDTO))
	page, err := c.dishService.PageQuery(&queryDTO)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, page)
}

// Delete 删除菜品
func (c *DishController) Delete(ctx *gin.Context) {
	// 解析 ids
	idsStr, isExist := ctx.GetQuery("ids")
	if !isExist {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	idStrS := strings.Split(idsStr, ",")
	var ids []int
	for _, idStr := range idStrS {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Error(constant.MsgBadRequest, zap.Error(err))
			response.BadRequest(ctx, constant.MsgBadRequest)
			return
		}
		ids = append(ids, id)
	}

	err := c.dishService.Delete(ids)
	if err != nil {
		logger.Error(constant.MsgDeleteFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgDeleteSuccess, nil)
}

// GetByID 根据ID获取菜品
func (c *DishController) GetByID(ctx *gin.Context) {
	// 获取路径参数
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层获取菜品
	dishVO, err := c.dishService.GetByID(id)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err), zap.Int("id", id))
		response.ErrorResponse(ctx, err)
		return
	}

	response.Success(ctx, constant.MsgQuerySuccess, dishVO)
}

// Update 修改菜品信息
func (c *DishController) Update(ctx *gin.Context) {
	var dishDTO dto.DishDTO
	if err := ctx.ShouldBindJSON(&dishDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.dishService.Update(ctx, &dishDTO)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}

// ☆ strconv.Atoi 函数会将空字符串 "" 转换为 0，并且不会返回错误。这个行为是符合 Go 语言的标准库设计的，空字符串被视为数值 0。

// UpdateStatus 菜品起售、停售
func (c *DishController) UpdateStatus(ctx *gin.Context) {
	statusStr := ctx.Param("status")
	if statusStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	status, err := strconv.Atoi(statusStr)
	if err != nil {
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

	err = c.dishService.UpdateStatus(id, status)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err), zap.Int("id", id), zap.Int("status", status))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}

// ListByCategoryID 根据分类ID查询菜品列表
func (c *DishController) ListByCategoryID(ctx *gin.Context) {
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
