package entity

import "takeout/model/wrap"

// Employee 员工数据模型
type Employee struct {
	ID         int            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Username   string         `json:"username" gorm:"unique;not null"`
	Password   string         `json:"password" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Phone      string         `json:"phone" gorm:"default:null"`
	Sex        string         `json:"sex" gorm:"default:null"`
	IdNumber   string         `json:"idNumber" gorm:"column:id_number;default:null"`
	Status     int            `json:"status" gorm:"default:1"`
	CreateTime wrap.LocalTime `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime wrap.LocalTime `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
	CreateUser int            `json:"createUser" gorm:"column:create_user;default:null"`
	UpdateUser int            `json:"updateUser" gorm:"column:update_user;default:null"`
}

// TableName 设置表名，gorm 通过这个函数来锁定表名
func (Employee) TableName() string {
	return "employee"
}
