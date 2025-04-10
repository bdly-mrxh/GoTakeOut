package entity

import (
	"github.com/shopspring/decimal"
	"takeout/model/wrap"
)

// ShoppingCart 购物车实体
type ShoppingCart struct {
	ID         int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name       string          `json:"name"`
	Image      string          `json:"image"`
	UserID     int             `json:"userId" gorm:"column:user_id;not null"`
	DishID     int             `json:"dishId" gorm:"column:dish_id"`
	SetmealID  int             `json:"setmealId" gorm:"column:setmeal_id"`
	DishFlavor string          `json:"dishFlavor" gorm:"column:dish_flavor"`
	Number     int             `json:"number" gorm:"not null;default:1"`
	Amount     decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not null"`
	CreateTime wrap.LocalTime  `json:"createTime" gorm:"column:create_time;autoCreateTime"`
}

// TableName 指定表名
func (ShoppingCart) TableName() string {
	return "shopping_cart"
}
