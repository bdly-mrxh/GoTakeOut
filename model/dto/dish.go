package dto

import (
	"github.com/shopspring/decimal"
	"takeout/model/entity"
)

// DishDTO 新增菜品DTO
type DishDTO struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	CategoryID  int                  `json:"categoryId"`
	Price       decimal.Decimal      `json:"price"`
	Image       string               `json:"image"`
	Description string               `json:"description"`
	Status      int                  `json:"status"`
	Flavors     []*entity.DishFlavor `json:"flavors"`
}

// DishPageQueryDTO 菜品分页查询DTO
type DishPageQueryDTO struct {
	CategoryID int    `form:"categoryId"`
	Name       string `form:"name"`
	Page       int    `form:"page" bind:"required"`
	PageSize   int    `form:"pageSize" bind:"required"`
	Status     int    `form:"status"`
}
