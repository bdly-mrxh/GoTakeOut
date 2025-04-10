package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

// 套餐相关路由
func (r *AdminRouter) setmealRouter() {
	setmealService := admin.NewSetmealController()
	setmeal := r.admin.Group("/setmeal")
	setmeal.Use(middleware.JwtAdmin())
	{
		// 新增套餐
		setmeal.POST("", setmealService.Create)
		// 根据ID获取套餐的详细信息
		setmeal.GET("/:id", setmealService.GetByID)
		// 批量删除
		setmeal.DELETE("", setmealService.BatchDelete)
		// 分页查询
		setmeal.GET("/page", setmealService.PageQuery)
		// 修改套餐
		setmeal.PUT("", setmealService.Update)
		// 更改套餐状态
		setmeal.POST("/status/:status", setmealService.UpdateStatus)
	}
}
