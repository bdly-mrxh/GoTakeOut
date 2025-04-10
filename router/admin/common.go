package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

// CommonRouter 通用路由
func (r *AdminRouter) commonRouter() {
	// 通用控制器
	commonController := admin.NewCommonController()

	// 需要JWT认证的路由
	common := r.admin.Group("/common")
	common.Use(middleware.JwtAdmin())
	{
		// 文件上传
		common.POST("/upload", commonController.Upload)
	}
}
