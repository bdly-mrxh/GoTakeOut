package entity

import "takeout/model/wrap"

// Category 分类数据模型
type Category struct {
	ID         int            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Type       int            `json:"type" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Sort       int            `json:"sort"`
	Status     int            `json:"status"`
	CreateTime wrap.LocalTime `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime wrap.LocalTime `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
	CreateUser int            `json:"createUser" gorm:"column:create_user;default:null"`
	UpdateUser int            `json:"updateUser" gorm:"column:update_user;default:null"`
}

// TableName 设置表名
func (Category) TableName() string {
	return "category"
}
