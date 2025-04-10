package user

import (
	"takeout/internal/control/user"
)

func (r *UserRouter) shopRouter() {
	shop := r.user.Group("/shop")
	{
		shopController := user.NewShopController()
		// 获取店铺状态
		shop.GET("/status", shopController.GetStatus)
	}
}
