package vo

import (
	"github.com/shopspring/decimal"
	"takeout/model/entity"
	"takeout/model/wrap"
)

// DishVO 菜品视图对象
type DishVO struct {
	ID           int                  `json:"id"`
	Name         string               `json:"name"`
	CategoryID   int                  `json:"categoryId" gorm:"column:category_id"`
	Price        decimal.Decimal      `json:"price" gorm:"type:decimal(10,2)"`
	Image        string               `json:"image"`
	Description  string               `json:"description"`
	Status       int                  `json:"status"`
	UpdateTime   wrap.LocalTime       `json:"updateTime" gorm:"column:update_time"`
	CategoryName string               `json:"categoryName" gorm:"column:category_name"`
	Flavors      []*entity.DishFlavor `json:"flavors" gorm:"-"`
}

// DishItem 套餐内菜品信息
type DishItem struct {
	Name        string `json:"name"`
	Copies      int    `json:"copies"`
	Image       string `json:"image"`
	Description string `json:"description"`
}
