package entity

import "github.com/shopspring/decimal"

// OrderDetail 订单明细数据模型
type OrderDetail struct {
	ID         int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name       string          `json:"name"`
	OrderID    int             `json:"orderId" gorm:"column:order_id"`
	DishID     int             `json:"dishId" gorm:"column:dish_id"`
	SetmealID  int             `json:"setmealId" gorm:"column:setmeal_id"`
	DishFlavor string          `json:"dishFlavor" gorm:"column:dish_flavor"`
	Number     int             `json:"number"`
	Amount     decimal.Decimal `json:"amount"`
	Image      string          `json:"image"`
}

// TableName 指定表名
func (OrderDetail) TableName() string {
	return "order_detail"
}
