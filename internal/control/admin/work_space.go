package admin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
	"time"
)

type WorkSpaceController struct {
	workSpaceService service.WorkSpaceService
}

func NewWorkSpaceController() *WorkSpaceController {
	return &WorkSpaceController{}
}

// BusinessData 工作台今日数据查询
func (c *WorkSpaceController) BusinessData(ctx *gin.Context) {
	date := time.Now()
	begin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Nanosecond*time.Second-time.Nanosecond), time.Local)

	data, err := c.workSpaceService.GetBusinessData(&begin, &end)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, data)
}

// OrderOverView 获取订单概览
func (c *WorkSpaceController) OrderOverView(ctx *gin.Context) {
	data, err := c.workSpaceService.GetOrderOverView()
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, data)
}

// DishOverView 查询菜品
func (c *WorkSpaceController) DishOverView(ctx *gin.Context) {
	data, err := c.workSpaceService.GetDishOverView()
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, data)
}

// SetmealOverView 查询套餐
func (c *WorkSpaceController) SetmealOverView(ctx *gin.Context) {
	data, err := c.workSpaceService.GetSetmealOverView()
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, data)
}
