package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

func (r *AdminRouter) dishRouter() {
	dishController := admin.NewDishController()

	dish := r.admin.Group("/dish")
	dish.Use(middleware.JwtAdmin())
	{
		// 新增菜品
		dish.POST("", dishController.Create)
		// 分页查询菜品
		dish.GET("/page", dishController.PageQuery)
		// 批量删除菜品
		dish.DELETE("", dishController.Delete)
		// 根据 id 查询菜品信息
		dish.GET("/:id", dishController.GetByID)
		// 修改菜品信息
		dish.PUT("", dishController.Update)
		// 菜品起售、停售
		dish.POST("/status/:status", dishController.UpdateStatus)
		// 根据分类ID查询菜品列表
		dish.GET("/list", dishController.ListByCategoryID)
	}
}
