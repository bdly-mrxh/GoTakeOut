package dao

import (
	"gorm.io/gorm"
	"takeout/common/global"
	"takeout/model/entity"
)

// DishFlavorDAO 菜品口味数据访问对象
type DishFlavorDAO struct{}

// BatchCreateWithTx 使用事务批量创建菜品口味
func (dao *DishFlavorDAO) BatchCreateWithTx(flavors []*entity.DishFlavor, tx *gorm.DB) error {
	result := tx.Create(&flavors)
	return result.Error
}

// DeleteByDishIDsTx 根据 dishIDs 删除关联的口味数据
func (dao *DishFlavorDAO) DeleteByDishIDsTx(dishIds []int, tx *gorm.DB) error {
	result := tx.Where("dish_id in ?", dishIds).Delete(&entity.DishFlavor{})
	return result.Error
}

// GetByDishID 根据菜品ID查询口味数据
func (dao *DishFlavorDAO) GetByDishID(dishId int) ([]*entity.DishFlavor, error) {
	var flavors []*entity.DishFlavor
	result := global.DB.Model(&entity.DishFlavor{}).Where("dish_id = ?", dishId).Find(&flavors)
	return flavors, result.Error
}

// DeleteByDishIDTx 根据dishID删除口味数据
func (dao *DishFlavorDAO) DeleteByDishIDTx(dishId int, tx *gorm.DB) error {
	return tx.Where("dish_id = ?", dishId).Delete(&entity.DishFlavor{}).Error
}
