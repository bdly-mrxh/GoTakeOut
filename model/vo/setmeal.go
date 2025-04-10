package vo

import (
	"github.com/shopspring/decimal"
	"takeout/model/entity"
	"takeout/model/wrap"
)

// SetmealVO 套餐视图对象
type SetmealVO struct {
	ID            int                   `json:"id"`
	CategoryID    int                   `json:"categoryId" gorm:"column:category_id"`
	Name          string                `json:"name"`
	Price         decimal.Decimal       `json:"price"`
	Status        int                   `json:"status"`
	Description   string                `json:"description"`
	Image         string                `json:"image"`
	UpdateTime    wrap.LocalTime        `json:"updateTime" gorm:"column:update_time"`
	CategoryName  string                `json:"categoryName" gorm:"column:category_name"`
	SetmealDishes []*entity.SetmealDish `json:"setmealDishes" gorm:"-"`
}
