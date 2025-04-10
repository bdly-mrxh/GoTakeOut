package user

import (
	"takeout/internal/control/user"
	"takeout/internal/middleware"
)

func (r *UserRouter) orderRouter() {
	order := r.user.Group("/order")
	order.Use(middleware.JwtUser())
	{
		orderController := user.NewOrderController()
		// 提交订单
		order.POST("/submit", orderController.Submit)
		// 订单支付
		order.PUT("/payment", orderController.Payment)
		// 历史订单查询
		order.GET("/historyOrders", orderController.Page)
		// 查询订单详细信息
		order.GET("/orderDetail/:id", orderController.Detail)
		// 用户取消订单
		order.PUT("/cancel/:id", orderController.Cancel)
		// 再来一单
		order.POST("/repetition/:id", orderController.Repetition)
		// 用户催单
		order.GET("/reminder/:id", orderController.Reminder)
	}
}
