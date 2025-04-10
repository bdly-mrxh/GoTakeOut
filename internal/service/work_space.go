package service

import (
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/internal/dao"
	"takeout/model/vo"
	"time"
)

type WorkSpaceService struct {
	orderDAO   dao.OrderDAO
	userDAO    dao.UserDAO
	dishDAO    dao.DishDAO
	setmealDAO dao.SetmealDAO
}

// GetBusinessData 获取今日营业数据
func (s *WorkSpaceService) GetBusinessData(begin *time.Time, end *time.Time) (*vo.BusinessDataVO, error) {
	turnover, err := s.orderDAO.GetAmount(global.DB, begin, end)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	amount, _ := turnover.Float64()
	validCnt, err := s.orderDAO.GetCount(global.DB, begin, end, constant.Completed)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	totalCnt, err := s.orderDAO.GetCount(global.DB, begin, end, 0)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	completeRate := 0.0
	if totalCnt != 0 {
		completeRate = float64(validCnt) / float64(totalCnt)
	}
	avgPrice := 0.0
	if validCnt != 0 {
		avgPrice = amount / float64(validCnt)
	}
	newUsers, err := s.userDAO.GetCount(global.DB, begin, end)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return &vo.BusinessDataVO{
		Turnover:            amount,
		ValidOrderCount:     validCnt,
		OrderCompletionRate: completeRate,
		UnitPrice:           avgPrice,
		NewUsers:            newUsers,
	}, nil
}

// GetOrderOverView 获取订单概览
func (s *WorkSpaceService) GetOrderOverView() (*vo.OrderOverViewVO, error) {
	begin, _ := s.getDateTime()
	waitCnt, err := s.orderDAO.GetCount(global.DB, begin, nil, constant.Confirmed)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	deliveryCnt, err := s.orderDAO.GetCount(global.DB, begin, nil, constant.DeliveryInProgress)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	completeCnt, err := s.orderDAO.GetCount(global.DB, begin, nil, constant.Completed)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	cancelCnt, err := s.orderDAO.GetCount(global.DB, begin, nil, constant.Cancelled)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	total, err := s.orderDAO.GetCount(global.DB, begin, nil, 0)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return &vo.OrderOverViewVO{
		WaitingOrders:   waitCnt,
		DeliveredOrders: deliveryCnt,
		CompletedOrders: completeCnt,
		CancelledOrders: cancelCnt,
		AllOrders:       total,
	}, nil
}

// 获取今日的始末时间
func (s *WorkSpaceService) getDateTime() (*time.Time, *time.Time) {
	date := time.Now()
	begin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Nanosecond*time.Second-time.Nanosecond), time.Local)
	return &begin, &end
}

// GetDishOverView 获取菜品总览
func (s *WorkSpaceService) GetDishOverView() (*vo.DishOverViewVO, error) {
	start, err := s.dishDAO.GetCount(global.DB, constant.DishEnable)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	stop, err := s.dishDAO.GetCount(global.DB, constant.DishDisable)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return &vo.DishOverViewVO{
		Sold:         start,
		Discontinued: stop,
	}, nil
}

// GetSetmealOverView 获取套餐概览
func (s *WorkSpaceService) GetSetmealOverView() (*vo.SetmealOverViewVO, error) {
	start, err := s.setmealDAO.GetCount(global.DB, constant.SetmealEnable)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	stop, err := s.setmealDAO.GetCount(global.DB, constant.SetmealDisable)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return &vo.SetmealOverViewVO{
		Sold:         start,
		Discontinued: stop,
	}, nil
}
