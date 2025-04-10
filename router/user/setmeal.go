package user

import (
	"takeout/internal/control/user"
	"takeout/internal/middleware"
)

// 用户端套餐路由
func (r *UserRouter) setmealRouter() {
	setmeal := r.user.Group("/setmeal")
	setmeal.Use(middleware.JwtUser())
	{
		setmealController := user.NewSetmealController()
		// 根据分类ID查询起售的套餐列表
		setmeal.GET("/list", setmealController.List)
		// 根据套餐ID查询包含的菜品列表
		setmeal.GET("/dish/:id", setmealController.DishList)
	}
}
