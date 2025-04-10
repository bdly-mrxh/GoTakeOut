package router

import (
	"takeout/common/global"
	middleware2 "takeout/internal/middleware"
	"takeout/internal/websocket"
	"takeout/router/admin"
	"takeout/router/notify"
	"takeout/router/user"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.Server.Mode)

	// 创建gin实例
	r := gin.New()

	// 使用自定义中间件
	r.Use(middleware2.LoggerMiddleware(), middleware2.RecoveryMiddleware())

	// 健康检查
	//r.GET("/ping", func(c *gin.Context) {
	//	response.Success(c, "pong", nil)
	//})

	// API分组
	api := r.Group("")
	{
		// 管理端API
		adminGroup := api.Group("/admin")
		{
			adminRouter := admin.NewAdminRouter(adminGroup)
			// 注册管理端路由
			adminRouter.RegisterRoutes()
		}

		// 用户端API
		userGroup := api.Group("/user")
		{
			userRouter := user.NewUserRouter(userGroup)
			// 注册用户端路由
			userRouter.RegisterRouters()
		}

		// 微信回调API
		notifyGroup := api.Group("/notify")
		{
			notifyRouter := notify.NewNotifyRouter(notifyGroup)
			notifyRouter.RegisterRouters()
		}
	}

	// ws路由
	r.GET("/ws/:sid", websocket.WSHandler)

	return r
}
