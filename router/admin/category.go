package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

// CategoryRouter 分类相关路由
func (r *AdminRouter) categoryRouter() {
	category := r.admin.Group("category")
	category.Use(middleware.JwtAdmin())
	{
		categoryController := admin.NewCategoryController()
		// 新增分类
		category.POST("", categoryController.Create)
		// 启用/禁用分类
		category.POST("/status/:status", categoryController.UpdateStatus)
		// 修改分类
		category.PUT("", categoryController.Update)
		// 分类分页查询
		category.GET("/page", categoryController.PageQuery)
		// 根据类型查询分类
		category.GET("/list", categoryController.List)
		// 删除分类
		category.DELETE("", categoryController.Delete)
	}
}
