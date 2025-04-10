package dto

// EmployeeLoginDTO 员工登录请求DTO
type EmployeeLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/*
  	两套DTO会增加冗余和管理难度
	共用DTO binding 需要放宽，要增加判断
*/

// EmployeeDTO 创建和更新共用的DTO
type EmployeeDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Sex      string `json:"sex"`
	IdNumber string `json:"idNumber"`
}

// EmployeeCreateDTO 员工创建请求DTO
type EmployeeCreateDTO struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Sex      string `json:"sex"`
	IdNumber string `json:"idNumber"`
}

// EmployeeUpdateDTO 员工信息更新请求DTO
type EmployeeUpdateDTO struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Sex      string `json:"sex"`
	IdNumber string `json:"idNumber"`
}

// EmployeePasswordDTO 员工密码修改请求DTO
type EmployeePasswordDTO struct {
	EmpId       int    `json:"empId" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// EmployeePageDTO 员工分页查询参数
type EmployeePageDTO struct {
	Name     string `form:"name"`     // 员工姓名，可选
	Page     int    `form:"page"`     // 页码
	PageSize int    `form:"pageSize"` // 每页记录数
}
