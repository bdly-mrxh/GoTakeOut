package admin

import (
	"takeout/internal/control/admin"
	"takeout/internal/middleware"
)

// EmployeeRouter 员工相关路由
func (r *AdminRouter) employeeRouter() {
	// 员工相关路由
	employeeController := admin.NewEmployeeController()
	r.admin.POST("/employee/login", employeeController.Login)

	// 需要JWT认证的路由
	employee := r.admin.Group("/employee")
	employee.Use(middleware.JwtAdmin())
	{
		// 新增员工
		employee.POST("", employeeController.Create)
		// 根据id查询员工信息
		employee.GET("/:id", employeeController.GetById)
		// 分页查询员工信息
		employee.GET("/page", employeeController.Page)
		// 更新员工信息
		employee.PUT("", employeeController.Update)
		// 更新员工状态 (通过查询参数接收员工ID)
		employee.POST("/status/:status", employeeController.UpdateStatus)
		// 修改密码
		employee.PUT("/editPassword", employeeController.UpdatePassword)
		// 退出登录
		employee.POST("/logout", employeeController.Logout)
	}
}
