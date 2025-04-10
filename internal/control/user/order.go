package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"takeout/common/constant"
	"takeout/common/errs"
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

// Submit 提交订单
func (c *OrderController) Submit(ctx *gin.Context) {
	var submitDTO dto.OrderSubmitDTO
	if err := ctx.ShouldBindJSON(&submitDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	submitVO, err := c.orderService.Submit(ctx, &submitDTO)
	if err != nil {
		var e *errs.Error
		if errors.As(err, &e); e.Code == constant.CodeBusinessError {
			logger.Info(constant.MsgOrderSubmitFail, zap.Error(err))
		} else {
			logger.Error(constant.MsgOrderSubmitFail, zap.Error(err))
		}
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgOrderSubmitSuccess, submitVO)
}

// Payment 订单支付
func (c *OrderController) Payment(ctx *gin.Context) {
	var payDTO dto.OrderPaymentDTO
	if err := ctx.ShouldBindJSON(&payDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	payVO, err := c.orderService.Payment(&payDTO)
	if err != nil {
		var e *errs.Error
		if errors.As(err, &e); e.Code == constant.CodeBusinessError {
			logger.Info(constant.MsgBusinessError, zap.Error(err))
		} else {
			logger.Error(constant.MsgPayFail, zap.Error(err))
		}
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, payVO)
}

// RealPayment 真实的微信支付
func (c *OrderController) RealPayment(ctx *gin.Context) {
	var payDTO dto.OrderPaymentDTO
	if err := ctx.ShouldBindJSON(&payDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	payVO, err := c.orderService.RealPayment(ctx, &payDTO)
	if err != nil {
		var e *errs.Error
		if errors.As(err, &e); e.Code == constant.CodeBusinessError {
			logger.Info(constant.MsgBusinessError, zap.Error(err))
		} else {
			logger.Error(constant.MsgPayFail, zap.Error(err))
		}
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, payVO)
}

// Page 查询历史订单
func (c *OrderController) Page(ctx *gin.Context) {
	var queryDTO dto.OrderPageQueryDTO
	if err := ctx.ShouldBindQuery(&queryDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	page, err := c.orderService.Page(ctx, &queryDTO)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, page)
}

// Detail 查看订单信息
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

	dishVO, err := c.orderService.Detail(id)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, dishVO)
}

// Cancel 用户取消订单
func (c *OrderController) Cancel(ctx *gin.Context) {
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

	err = c.orderService.CancelByUser(id)
	if err != nil {
		logger.Error(constant.MsgOrderCancelFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgOrderCancelSuccess, nil)
}

// Repetition 再来一单
func (c *OrderController) Repetition(ctx *gin.Context) {
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

	err = c.orderService.Repetition(ctx, id)
	if err != nil {
		logger.Error(constant.MsgCreateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, nil)
}

// Reminder 用户催单
func (c *OrderController) Reminder(ctx *gin.Context) {
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

	err = c.orderService.Reminder(id)
	if err != nil {
		logger.Error(constant.MsgServerError, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, nil)
}
