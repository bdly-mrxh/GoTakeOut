package dao

import (
	"gorm.io/gorm"
	"takeout/model/entity"
)

type ShoppingCartDAO struct{}

// List 动态查询购物车列表
func (d *ShoppingCartDAO) List(db *gorm.DB, cart *entity.ShoppingCart) ([]*entity.ShoppingCart, error) {
	query := db.Model(&entity.ShoppingCart{})
	if cart.UserID != 0 {
		query = query.Where("user_id = ?", cart.UserID)
	}
	if cart.SetmealID != 0 {
		query = query.Where("setmeal_id = ?", cart.SetmealID)
	}
	if cart.DishID != 0 {
		query = query.Where("dish_id = ?", cart.DishID)
	}
	if cart.DishFlavor != "" {
		query = query.Where("dish_flavor = ?", cart.DishFlavor)
	}

	var list []*entity.ShoppingCart
	result := query.Find(&list)
	return list, result.Error
}

// UpdateNumberByID 更新数量
func (d *ShoppingCartDAO) UpdateNumberByID(db *gorm.DB, c *entity.ShoppingCart) error {
	return db.Model(&entity.ShoppingCart{}).Where("id = ?", c.ID).Update("number", c.Number).Error
}

// Create 添加购物车
func (d *ShoppingCartDAO) Create(db *gorm.DB, cart *entity.ShoppingCart) error {
	return db.Model(&entity.ShoppingCart{}).Create(cart).Error
}

// CleanByUserID 根据用户ID清空购物车
func (d *ShoppingCartDAO) CleanByUserID(db *gorm.DB, userID int) error {
	return db.Where("user_id = ?", userID).Delete(&entity.ShoppingCart{}).Error
}

// DeleteByID 删除购物车中的一项
func (d *ShoppingCartDAO) DeleteByID(db *gorm.DB, cart *entity.ShoppingCart) error {
	return db.Where("id = ?", cart.ID).Delete(&entity.ShoppingCart{}).Error
}

// BatchInsert 批量插入购物车
func (d *ShoppingCartDAO) BatchInsert(db *gorm.DB, list []*entity.ShoppingCart) error {
	return db.Model(&entity.ShoppingCart{}).Create(&list).Error
}
