package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

func (r *AdminRouter) workSpaceRouter() {
	workspace := r.admin.Group("/workspace")
	workspace.Use(middleware.JwtAdmin())
	{
		workSpaceController := admin.NewWorkSpaceController()
		workspace.GET("/businessData", workSpaceController.BusinessData)
		workspace.GET("/overviewOrders", workSpaceController.OrderOverView)
		workspace.GET("/overviewDishes", workSpaceController.DishOverView)
		workspace.GET("/overviewSetmeals", workSpaceController.SetmealOverView)
	}
}
