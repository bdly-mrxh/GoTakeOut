package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

func (r *AdminRouter) shopRouter() {
	shop := r.admin.Group("/shop")
	shop.Use(middleware.JwtAdmin())
	{
		shopController := admin.NewShopController()
		// 设置店铺状态
		shop.PUT("/:status", shopController.SetStatus)
		// 获取店铺状态
		shop.GET("/status", shopController.GetStatus)
	}
}
