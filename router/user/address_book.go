package user

import (
	"takeout/internal/control/user"
	"takeout/internal/middleware"
)

func (r *UserRouter) addressBookRouter() {
	addressBook := r.user.Group("/addressBook")
	addressBook.Use(middleware.JwtUser())
	{
		addressBookController := user.NewAddressBookController()
		// 查询当前登录用户的所有地址
		addressBook.GET("/list", addressBookController.List)
		// 新增地址
		addressBook.POST("", addressBookController.Add)
		// 根据ID查询地址
		addressBook.GET("/:id", addressBookController.GetByID)
		// 根据ID修改地址
		addressBook.PUT("", addressBookController.Update)
		// 设置默认地址
		addressBook.PUT("/default", addressBookController.SetDefault)
		// 根据ID删除地址
		addressBook.DELETE("", addressBookController.DeleteByID)
		// 查询默认地址
		addressBook.GET("/default", addressBookController.GetDefault)
	}
}
