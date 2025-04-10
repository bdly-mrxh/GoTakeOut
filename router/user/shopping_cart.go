package user

import (
	"takeout/internal/control/user"
	"takeout/internal/middleware"
)

func (r *UserRouter) shoppingCartRouter() {
	shoppingCart := r.user.Group("/shoppingCart")
	shoppingCart.Use(middleware.JwtUser())
	{
		shoppingCartController := user.NewShoppingCartController()
		// 添加购物车
		shoppingCart.POST("/add", shoppingCartController.Add)
		// 查看购物车
		shoppingCart.GET("/list", shoppingCartController.List)
		// 清空购物车
		shoppingCart.DELETE("/clean", shoppingCartController.Clean)
		// 删除购物车中的一个商品
		shoppingCart.POST("/sub", shoppingCartController.Sub)
	}
}
