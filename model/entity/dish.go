package entity

import (
	"github.com/shopspring/decimal"
	"takeout/model/wrap"
)

// Dish 菜品数据模型
type Dish struct {
	ID          int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"not null;uniqueIndex:idx_dish_name"`
	CategoryID  int             `json:"categoryId" gorm:"column:category_id;not null"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2)"`
	Image       string          `json:"image"`
	Description string          `json:"description"`
	Status      int             `json:"status" gorm:"default:1"`
	CreateTime  wrap.LocalTime  `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	UpdateTime  wrap.LocalTime  `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
	CreateUser  int             `json:"createUser" gorm:"column:create_user;default:null"`
	UpdateUser  int             `json:"updateUser" gorm:"column:update_user;default:null"`
}

// TableName 设置表名
func (Dish) TableName() string {
	return "dish"
}
