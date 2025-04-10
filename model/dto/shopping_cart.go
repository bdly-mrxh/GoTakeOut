package dto

// ShoppingCartDTO 添加购物车DTO
type ShoppingCartDTO struct {
	DishID     int    `json:"dishId"`
	SetmealID  int    `json:"setmealId"`
	DishFlavor string `json:"dishFlavor"`
}
