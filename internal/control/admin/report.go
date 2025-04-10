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

type ReportController struct {
	reportService service.ReportService
}

func NewReportController() *ReportController {
	return &ReportController{}
}

// Date 标签 uri(param) form(query) json(body)
type Date struct {
	Begin time.Time `form:"begin" time_format:"2006-01-02"`
	End   time.Time `form:"end" time_format:"2006-01-02"`
}

// TurnoverStatistics 营业额统计
func (c *ReportController) TurnoverStatistics(ctx *gin.Context) {
	var date Date
	if err := ctx.ShouldBindQuery(&date); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	ret, err := c.reportService.TurnoverStatistics(date.Begin, date.End)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, ret)
}

// UserStatistics 用户量统计
func (c *ReportController) UserStatistics(ctx *gin.Context) {
	var date Date
	if err := ctx.ShouldBindQuery(&date); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	ret, err := c.reportService.UserStatistics(date.Begin, date.End)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, ret)
}

// OrderStatistics 订单统计
func (c *ReportController) OrderStatistics(ctx *gin.Context) {
	var date Date
	if err := ctx.ShouldBindQuery(&date); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	ret, err := c.reportService.OrderStatistics(date.Begin, date.End)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, ret)
}

// SalesTop10Statistics 销量top10
func (c *ReportController) SalesTop10Statistics(ctx *gin.Context) {
	var date Date
	if err := ctx.ShouldBindQuery(&date); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	ret, err := c.reportService.SalesTop10Statistics(date.Begin, date.End)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgSuccess, ret)
}

// Export 导出运营数据表
func (c *ReportController) Export(ctx *gin.Context) {
	c.reportService.Export(ctx)
}
