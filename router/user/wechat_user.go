package user

import (
	"takeout/internal/control/user"
	"takeout/internal/middleware"
)

func (r *UserRouter) weChatUserRouter() {
	userController := user.NewWeChatUserController()
	r.user.POST("/user/login", userController.Login)
	weChatUser := r.user.Group("/user")
	weChatUser.Use(middleware.JwtUser())
	{
		weChatUser.POST("/logout", userController.Logout)
	}
}
