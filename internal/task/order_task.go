package task

import (
	"go.uber.org/zap"
	"takeout/common/constant"
	"takeout/common/global"
	"takeout/common/logger"
	"takeout/internal/dao"
	"takeout/model/wrap"
	"time"

	"github.com/robfig/cron/v3"
)

// Init 定时任务初始化，也可以改为 init()
func Init() error {
	orderTask := NewOrderTask()
	timerTask := cron.New(cron.WithSeconds())
	if _, err := timerTask.AddFunc("0 * * * * ?", orderTask.handleTimeoutOrder); err != nil {
		logger.Error("初始化定时任务失败", zap.Error(err))
		return err
	}
	if _, err := timerTask.AddFunc("0 0 1 * * ?", orderTask.handleDeliveryOrder); err != nil {
		logger.Error("初始化定时任务失败", zap.Error(err))
		return err
	}
	timerTask.Start()
	return nil
}

type OrderTask struct {
	orderDAO dao.OrderDAO
}

func NewOrderTask() *OrderTask {
	return &OrderTask{}
}

// 处理超时订单
func (t *OrderTask) handleTimeoutOrder() {
	logger.Info("处理超时订单", zap.Time("time", time.Now()))
	// 只能报告错误不能，抛出错误
	orders, err := t.orderDAO.GetByStatusLT(global.DB, constant.PendingPayment, time.Now().Add(-15*time.Minute))
	if err != nil {
		logger.Error(constant.MsgDatabaseError, zap.Error(err))
		return
	}
	if orders != nil && len(orders) != 0 {
		for _, order := range orders {
			order.Status = constant.Cancelled
			order.CancelReason = "订单超时"
			order.CancelTime = wrap.LocalTime(time.Now())
			err = t.orderDAO.Update(global.DB, order)
			if err != nil {
				logger.Error(constant.MsgDatabaseError, zap.Error(err))
				return
			}
		}
	}
}

// 处理一直在派送中的订单
func (t *OrderTask) handleDeliveryOrder() {
	logger.Info("处理未完成订单", zap.Time("time", time.Now()))
	orders, err := t.orderDAO.GetByStatusLT(global.DB, constant.DeliveryInProgress, time.Now().Add(-time.Hour))
	if err != nil {
		logger.Error(constant.MsgDatabaseError, zap.Error(err))
		return
	}
	if orders != nil && len(orders) != 0 {
		for _, order := range orders {
			order.Status = constant.Completed
			if err = t.orderDAO.Update(global.DB, order); err != nil {
				logger.Error(constant.MsgDatabaseError, zap.Error(err))
				return
			}
		}
	}
}
