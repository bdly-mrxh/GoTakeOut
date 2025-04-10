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

type OrderController struct {
	orderService service.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{}
}

// Search 根据条件搜索订单
func (c *OrderController) Search(ctx *gin.Context) {
	var queryDTO dto.OrderPageQueryDTO
	if err := ctx.ShouldBindQuery(&queryDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	list, err := c.orderService.Search(&queryDTO)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, list)
}

// Statistics 各个状态的订单数量统计
func (c *OrderController) Statistics(ctx *gin.Context) {
	ret, err := c.orderService.Statistics()
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, ret)
}

// Detail 查询订单详情
func (c *OrderController) Detail(ctx *gin.Context) {
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

	orderVO, err := c.orderService.Detail(id)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, orderVO)
}

// Confirm 接单
func (c *OrderController) Confirm(ctx *gin.Context) {
	var confirmDTO dto.OrderConfirmDTO
	if err := ctx.ShouldBindJSON(&confirmDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.orderService.Confirm(&confirmDTO)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}

// Reject 拒单
func (c *OrderController) Reject(ctx *gin.Context) {
	var rejectDTO dto.OrderRejectionDTO
	if err := ctx.ShouldBindJSON(&rejectDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.orderService.Reject(&rejectDTO)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, nil)
}

// Cancel 取消订单
func (c *OrderController) Cancel(ctx *gin.Context) {
	var cancelDTO dto.OrderCancelDTO
	if err := ctx.ShouldBindJSON(&cancelDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.orderService.Cancel(&cancelDTO)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, nil)
}

// Delivery 派送订单
func (c *OrderController) Delivery(ctx *gin.Context) {
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

	err = c.orderService.Delivery(id)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, nil)
}

// Complete 完成订单
func (c *OrderController) Complete(ctx *gin.Context) {
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

	err = c.orderService.Complete(id)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, nil)
}
