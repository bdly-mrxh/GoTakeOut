package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

func (r *AdminRouter) reportRouter() {
	report := r.admin.Group("/report")
	report.Use(middleware.JwtAdmin())
	{
		reportController := admin.NewReportController()
		report.GET("/turnoverStatistics", reportController.TurnoverStatistics)
		report.GET("/userStatistics", reportController.UserStatistics)
		report.GET("/ordersStatistics", reportController.OrderStatistics)
		report.GET("/top10", reportController.SalesTop10Statistics)
		report.GET("/export", reportController.Export)
	}
}
