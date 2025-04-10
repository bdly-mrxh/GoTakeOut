package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/logger"
	"takeout/internal/dao"
	"takeout/model/vo"
	"time"

	"github.com/xuri/excelize/v2"
)

type ReportService struct {
	userDAO          dao.UserDAO
	orderDAO         dao.OrderDAO
	workSpaceService WorkSpaceService
}

// TurnoverStatistics 统计营业额
func (s *ReportService) TurnoverStatistics(begin time.Time, end time.Time) (*vo.TurnoverReportVO, error) {
	dates := s.getEveryDate(begin, end)
	var amounts []any
	for _, date := range dates {
		beginTime, endTime := s.getDateTime(&date)
		amount, err := s.orderDAO.GetAmount(global.DB, beginTime, endTime)
		if err != nil {
			return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		amounts = append(amounts, amount)
	}
	dateList := s.getDateList(dates)
	amountList := s.getList(amounts)
	return &vo.TurnoverReportVO{DateList: dateList, TurnoverList: amountList}, nil
}

// UserStatistics 用户量统计
func (s *ReportService) UserStatistics(begin time.Time, end time.Time) (*vo.UserReportVO, error) {
	dates := s.getEveryDate(begin, end)
	var newUsers []any
	var totalUsers []any
	for _, date := range dates {
		beginTime, endTime := s.getDateTime(&date)
		newCnt, err := s.userDAO.GetCount(global.DB, beginTime, endTime)
		if err != nil {
			return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		newUsers = append(newUsers, newCnt)
		totalCnt, err := s.userDAO.GetCount(global.DB, nil, endTime)
		if err != nil {
			return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		totalUsers = append(totalUsers, totalCnt)
	}
	dateList := s.getDateList(dates)
	newUserList := s.getList(newUsers)
	totalUserList := s.getList(totalUsers)
	return &vo.UserReportVO{DateList: dateList, NewUserList: newUserList, TotalUserList: totalUserList}, nil
}

// OrderStatistics 订单统计
func (s *ReportService) OrderStatistics(begin time.Time, end time.Time) (*vo.OrderReportVO, error) {
	dates := s.getEveryDate(begin, end)
	var totalNum []any
	var validNum []any
	var totalCnt, validCnt int64
	for _, date := range dates {
		beginTime, endTime := s.getDateTime(&date)
		tCnt, err := s.orderDAO.GetCount(global.DB, beginTime, endTime, 0)
		if err != nil {
			return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		totalCnt += tCnt
		totalNum = append(totalNum, tCnt)
		vCnt, err := s.orderDAO.GetCount(global.DB, beginTime, endTime, constant.Completed)
		if err != nil {
			return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		validCnt += vCnt
		validNum = append(validNum, vCnt)
	}
	dateList := s.getDateList(dates)
	totalList := s.getList(totalNum)
	validList := s.getList(validNum)
	rate := 0.0
	if totalCnt != 0 {
		rate = float64(validCnt) / float64(totalCnt)
	}
	return &vo.OrderReportVO{
		DateList:            dateList,
		ValidOrderCountList: validList,
		OrderCountList:      totalList,
		TotalOrderCount:     totalCnt,
		ValidOrderCount:     validCnt,
		OrderCompletionRate: rate,
	}, nil
}

// SalesTop10Statistics 统计销量前10
func (s *ReportService) SalesTop10Statistics(begin time.Time, end time.Time) (*vo.SalesTop10ReportVO, error) {
	beginTime := time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local)
	endTime := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, int(time.Nanosecond*time.Second-time.Nanosecond), time.Local)
	list, err := s.orderDAO.GetSalesTop10(global.DB, &beginTime, &endTime)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	var names []any
	var nums []any
	for _, item := range list {
		names = append(names, item.Name)
		nums = append(nums, item.Number)
	}
	nameList := s.getList(names)
	numberList := s.getList(nums)
	return &vo.SalesTop10ReportVO{NameList: nameList, NumberList: numberList}, nil
}

// 获取日期列表
func (s *ReportService) getDateList(dates []time.Time) string {
	n := len(dates)
	if n == 0 {
		return ""
	} else if n == 1 {
		return dates[0].Format("2006-01-02")
	}
	str := ""
	for i := 0; i < n-1; i++ {
		str += dates[i].Format("2006-01-02") + ","
	}
	return str + dates[n-1].Format("2006-01-02")
}

// 获取列表
func (s *ReportService) getList(objects []any) string {
	n := len(objects)
	if n == 0 {
		return ""
	} else if n == 1 {
		return fmt.Sprint(objects[0])
	}
	str := ""
	for i := 0; i < n-1; i++ {
		str += fmt.Sprint(objects[i]) + ","
	}
	return str + fmt.Sprint(objects[n-1])
}

// 返回一天的起始和结束时间
func (s *ReportService) getDateTime(date *time.Time) (*time.Time, *time.Time) {
	// 起始时间：00:00:00
	beginTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	// 结束时间：23:59:59.999999999
	endTime := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Nanosecond*time.Second-time.Nanosecond), time.Local)
	return &beginTime, &endTime
}

// 获取区间内的所有日期
func (s *ReportService) getEveryDate(begin, end time.Time) []time.Time {
	dates := make([]time.Time, 0)
	dates = append(dates, begin)
	for begin.Day() != end.Day() || begin.Month() != end.Month() || begin.Year() != end.Year() {
		begin = begin.Add(24 * time.Hour)
		dates = append(dates, begin)
	}
	return dates
}

// Export 导出数据表
func (s *ReportService) Export(ctx *gin.Context) {
	today := time.Now()
	beginDate := today.AddDate(0, 0, -30)
	endDate := today.Add(-24 * time.Hour)
	beginTime := time.Date(beginDate.Year(), beginDate.Month(), beginDate.Day(), 0, 0, 0, 0, time.Local)
	endTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Nanosecond*time.Second-time.Nanosecond), time.Local)
	businessData, err := s.workSpaceService.GetBusinessData(&beginTime, &endTime)
	if err != nil {
		logger.Error(constant.MsgDatabaseError, zap.Error(err))
		return
	}
	excel, err := excelize.OpenFile(global.Config.Template.Path)
	if err != nil {
		logger.Error(constant.MsgServerError, zap.Error(err))
		return
	}
	sheet := "Sheet1"
	timeStr := fmt.Sprintf("时间：%s 至 %s", beginDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	s.setCellValueSafe(excel, sheet, "B2", timeStr)
	s.setCellValueSafe(excel, sheet, "C4", businessData.Turnover)
	s.setCellValueSafe(excel, sheet, "E4", businessData.OrderCompletionRate)
	s.setCellValueSafe(excel, sheet, "G4", businessData.NewUsers)
	s.setCellValueSafe(excel, sheet, "C5", businessData.ValidOrderCount)
	s.setCellValueSafe(excel, sheet, "E5", businessData.UnitPrice)
	for i := 0; i < 30; i++ {
		date := beginDate.AddDate(0, 0, i)
		begin, end := s.getDateTime(&date)
		data, e := s.workSpaceService.GetBusinessData(begin, end)
		if e != nil {
			logger.Error(constant.MsgDatabaseError, zap.Error(err))
			return
		}
		row := i + 8
		s.setCellValueSafe(excel, sheet, fmt.Sprintf("B%d", row), date.Format("2006-01-02"))
		s.setCellValueSafe(excel, sheet, fmt.Sprintf("C%d", row), data.Turnover)
		s.setCellValueSafe(excel, sheet, fmt.Sprintf("D%d", row), data.ValidOrderCount)
		s.setCellValueSafe(excel, sheet, fmt.Sprintf("E%d", row), data.OrderCompletionRate)
		s.setCellValueSafe(excel, sheet, fmt.Sprintf("F%d", row), data.UnitPrice)
		s.setCellValueSafe(excel, sheet, fmt.Sprintf("G%d", row), data.NewUsers)
	}
	// 浏览器会知道该文件是一个 Excel 文件，并按照 Excel 文件的处理方式（如预览、下载或打开方式等）进行处理
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err = excel.Write(ctx.Writer); err != nil {
		logger.Error(constant.MsgServerError, zap.Error(err))
		return
	}
}

// setCellValueSafe 安全地设置单元格值，统一处理错误
func (s *ReportService) setCellValueSafe(excel *excelize.File, sheet, axis string, value any) {
	if err := excel.SetCellValue(sheet, axis, value); err != nil {
		logger.Error(constant.MsgServerError, zap.String("err", "填写excel出错"))
	}
}
