package user

import (
	"takeout/internal/control/user"
	"takeout/internal/middleware"
)

// 分类相关路由
func (r *UserRouter) categoryRouter() {
	category := r.user.Group("/category")
	category.Use(middleware.JwtUser())
	{
		categoryController := user.NewCategoryController()
		// 查询分类
		category.GET("/list", categoryController.List)
	}
}
