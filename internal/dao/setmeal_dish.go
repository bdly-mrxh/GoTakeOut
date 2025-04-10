package dao

import (
	"gorm.io/gorm"
	"takeout/common/global"
	"takeout/model/entity"
)

// SetmealDishDAO 套餐菜品关系数据访问对象
type SetmealDishDAO struct{}

// CountByDishIDs 根据菜品ID列表统计关联的套餐数量
func (dao *SetmealDishDAO) CountByDishIDs(dishIDs []int) (int64, error) {
	var count int64
	result := global.DB.Table("setmeal_dish").Where("dish_id in ?", dishIDs).Count(&count)
	return count, result.Error
}

// BatchInsert 批量插入
func (dao *SetmealDishDAO) BatchInsert(db *gorm.DB, dishes []*entity.SetmealDish) error {
	return db.Model(&entity.SetmealDish{}).Create(&dishes).Error
}

// GetBySetmealID 根据套餐ID查找
func (dao *SetmealDishDAO) GetBySetmealID(db *gorm.DB, setmealID int) ([]*entity.SetmealDish, error) {
	var setmealDishes []*entity.SetmealDish
	result := db.Model(&entity.SetmealDish{}).Where("setmeal_id = ?", setmealID).Find(&setmealDishes)
	return setmealDishes, result.Error
}

// BatchDeleteBySetmealIDs 根据套餐IDs批量删除
func (dao *SetmealDishDAO) BatchDeleteBySetmealIDs(db *gorm.DB, setmealIds []int) error {
	result := db.Where("setmeal_id in ?", setmealIds).Delete(&entity.SetmealDish{})
	return result.Error
}

// BatchDeleteBySetmealID 根据套餐ID删除
func (dao *SetmealDishDAO) BatchDeleteBySetmealID(db *gorm.DB, setmealId int) error {
	result := db.Where("setmeal_id = ?", setmealId).Delete(&entity.SetmealDish{})
	return result.Error
}

// GetDishIdsBySetmealId 根据套餐ID获取菜品的IDs
func (dao *SetmealDishDAO) GetDishIdsBySetmealId(db *gorm.DB, setmealId int) ([]int, error) {
	var dishIds []int
	result := db.Model(&entity.SetmealDish{}).Where("setmeal_id = ?", setmealId).Pluck("dish_id", &dishIds)
	return dishIds, result.Error
}
