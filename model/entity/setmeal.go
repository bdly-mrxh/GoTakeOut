package entity

import (
	"github.com/shopspring/decimal"
	"takeout/model/wrap"
)

// Setmeal 套餐数据模型
type Setmeal struct {
	ID          int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CategoryID  int             `json:"categoryId" gorm:"column:category_id;not null"`
	Name        string          `json:"name" gorm:"not null;uniqueIndex:idx_setmeal_name"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2)"`
	Status      int             `json:"status" gorm:"default:1"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	CreateTime  wrap.LocalTime  `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  wrap.LocalTime  `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
	CreateUser  int             `json:"createUser" gorm:"column:create_user;default:null"`
	UpdateUser  int             `json:"updateUser" gorm:"column:update_user;default:null"`
}

// TableName 设置表名
func (Setmeal) TableName() string {
	return "setmeal"
}
