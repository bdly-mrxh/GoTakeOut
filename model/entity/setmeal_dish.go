package entity

import (
	"github.com/shopspring/decimal"
)

// SetmealDish 套餐菜品关系数据模型
type SetmealDish struct {
	ID        int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	SetmealID int             `json:"setmealId" gorm:"column:setmeal_id"`
	DishID    int             `json:"dishId" gorm:"column:dish_id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price" gorm:"type:decimal(10,2)"`
	Copies    int             `json:"copies"`
}

// TableName 设置表名
func (SetmealDish) TableName() string {
	return "setmeal_dish"
}
