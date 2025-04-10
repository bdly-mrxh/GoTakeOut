package user

import (
	"takeout/internal/control/user"
	"takeout/internal/middleware"
)

func (r *UserRouter) dishRouter() {
	dish := r.user.Group("/dish")
	dish.Use(middleware.JwtUser())
	{
		dishController := user.NewDishController()
		// 根据分类ID查询菜品
		dish.GET("/list", dishController.List)
	}
}
