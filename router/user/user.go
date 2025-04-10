package user

import "github.com/gin-gonic/gin"

type UserRouter struct {
	user *gin.RouterGroup
}

func NewUserRouter(user *gin.RouterGroup) *UserRouter {
	return &UserRouter{user: user}
}

// RegisterRouters 用户端注册路由
func (r *UserRouter) RegisterRouters() {
	// 店铺相关接口
	r.shopRouter()
	// 用户微信登录接口
	r.weChatUserRouter()
	// 用户端分类接口
	r.categoryRouter()
	// 用户端菜品路由
	r.dishRouter()
	// 套餐路由
	r.setmealRouter()
	// 购物车路由
	r.shoppingCartRouter()
	// 地址簿路由
	r.addressBookRouter()
	// 订单路由
	r.orderRouter()
}
