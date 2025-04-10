package dto

import (
	"github.com/shopspring/decimal"
	"takeout/model/entity"
)

// SetmealDTO 套餐数据传输对象
type SetmealDTO struct {
	ID            int                   `json:"id"`
	CategoryID    int                   `json:"categoryId" binding:"required"`
	Name          string                `json:"name" binding:"required"`
	Price         decimal.Decimal       `json:"price" binding:"required"`
	Status        int                   `json:"status"`
	Description   string                `json:"description"`
	Image         string                `json:"image"`
	SetmealDishes []*entity.SetmealDish `json:"setmealDishes"`
}

// SetmealPageQueryDTO 套餐分页查询参数
type SetmealPageQueryDTO struct {
	Name       string `form:"name"`       // 套餐名称，可选
	CategoryID int    `form:"categoryId"` // 分类ID，可选
	Status     int    `form:"status"`     // 状态，0表示禁用，1表示启用，可选
	Page       int    `form:"page"`       // 页码
	PageSize   int    `form:"pageSize"`   // 每页记录数
}
