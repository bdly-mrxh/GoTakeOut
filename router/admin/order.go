package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

func (r *AdminRouter) orderRouter() {
	order := r.admin.Group("order")
	order.Use(middleware.JwtAdmin())
	{
		orderController := admin.NewOrderController()
		// 搜索订单
		order.GET("/conditionSearch", orderController.Search)
		// 各个状态的订单统计
		order.GET("/statistics", orderController.Statistics)
		// 查询订单详情
		order.GET("/details/:id", orderController.Detail)
		// 接单
		order.PUT("/confirm", orderController.Confirm)
		// 拒单
		order.PUT("/rejection", orderController.Reject)
		// 取消订单
		order.PUT("/cancel", orderController.Cancel)
		// 派送订单
		order.PUT("/delivery/:id", orderController.Delivery)
		// 完成订单
		order.PUT("/complete/:id", orderController.Complete)
	}
}
