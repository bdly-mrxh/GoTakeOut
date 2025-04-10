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

// SetmealController 套餐相关接口
type SetmealController struct {
	setmealService service.SetmealService
}

func NewSetmealController() *SetmealController {
	return &SetmealController{}
}

// Create 新增套餐
func (c *SetmealController) Create(ctx *gin.Context) {
	var createDTO dto.SetmealDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.setmealService.Create(ctx, &createDTO)
	if err != nil {
		logger.Error(constant.MsgCreateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgCreateSuccess, nil)
}

// GetByID 根据ID获取套餐的详细信息
func (c *SetmealController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgTypeConversionFail, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	setmealVO, err := c.setmealService.GetByID(id)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, setmealVO)
}

// BatchDelete 批量删除套餐
func (c *SetmealController) BatchDelete(ctx *gin.Context) {
	idsStr := ctx.Query("ids")
	if idsStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	// 解析ids列表
	idStrs := strings.Split(idsStr, ",")
	ids := make([]int, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.Error(constant.MsgBadRequest)
			response.BadRequest(ctx, constant.MsgBadRequest)
			return
		}
		ids = append(ids, id)
	}

	if err := c.setmealService.BatchDelete(ids); err != nil {
		logger.Error(constant.MsgDeleteFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgDeleteSuccess, nil)
}

// PageQuery 套餐分页查询
func (c *SetmealController) PageQuery(ctx *gin.Context) {
	var queryDTO dto.SetmealPageQueryDTO
	if err := ctx.ShouldBindQuery(&queryDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	status := ctx.Query("status")
	if status == "" {
		queryDTO.Status = constant.InvalidStatus
	}

	page, err := c.setmealService.PageQuery(&queryDTO)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, page)
}

// Update 更新套餐
func (c *SetmealController) Update(ctx *gin.Context) {
	var updateDTO dto.SetmealDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.setmealService.Update(ctx, &updateDTO)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}

// UpdateStatus 更新套餐状态
func (c *SetmealController) UpdateStatus(ctx *gin.Context) {
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
	idStr := ctx.Query("id")
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

	err = c.setmealService.UpdateStatus(id, status)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}
