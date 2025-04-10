package vo

import (
	"takeout/common/utils"
	"takeout/model/entity"
)

// EmployeeLoginVO 员工登录响应VO
type EmployeeLoginVO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Token    string `json:"token"`
}

// EmployeeDetailVO 员工详情响应VO
type EmployeeDetailVO struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"` // 会被掩码处理
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Sex        string `json:"sex"`
	IdNumber   string `json:"idNumber"`
	Status     int    `json:"status"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	CreateUser int    `json:"createUser"`
	UpdateUser int    `json:"updateUser"`
}

// EmployeeDetailVOFromEntity FromEntity 从实体转换为VO
func EmployeeDetailVOFromEntity(employee *entity.Employee) (*EmployeeDetailVO, error) {
	vo := &EmployeeDetailVO{}

	// 使用deepcopier复制实体到VO
	err := utils.CopyProperties(employee, vo)
	if err != nil {
		return nil, err
	}

	// 特殊处理
	vo.Password = "******" // 密码掩码
	vo.CreateTime = employee.CreateTime.String()
	vo.UpdateTime = employee.UpdateTime.String()

	return vo, nil
}
