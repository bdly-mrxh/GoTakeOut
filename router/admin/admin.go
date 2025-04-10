package admin

import "github.com/gin-gonic/gin"

// AdminRouter 管理端路由
type AdminRouter struct {
	admin *gin.RouterGroup
}

// NewAdminRouter 初始化管理端路由
func NewAdminRouter(admin *gin.RouterGroup) *AdminRouter {
	return &AdminRouter{admin: admin}
}

// RegisterRoutes 注册所有路由
func (r *AdminRouter) RegisterRoutes() {
	// 注册员工相关接口
	r.employeeRouter()
	// 注册分类路由
	r.categoryRouter()
	// 注册通用路由
	r.commonRouter()
	// 注册菜品路由
	r.dishRouter()
	// 注册套餐路由
	r.setmealRouter()
	// 注册店铺路由
	r.shopRouter()
	// 注册订单路由
	r.orderRouter()
	// 注册统计路由
	r.reportRouter()
	// 注册工作台路由
	r.workSpaceRouter()
}
