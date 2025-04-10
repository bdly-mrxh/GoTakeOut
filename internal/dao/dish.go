package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/model/entity"
	"takeout/model/vo"
)

// DishDAO 菜品数据访问对象
type DishDAO struct{}

// CountByCategoryID 根据分类ID统计菜品数量
func (dao *DishDAO) CountByCategoryID(categoryID int) (int64, error) {
	var count int64
	result := global.DB.Model(&entity.Dish{}).Where("category_id = ?", categoryID).Count(&count)
	return count, result.Error
}

// CreateWithTx 可以进行事务操作的创建新菜品
func (dao *DishDAO) CreateWithTx(ctx *gin.Context, dish *entity.Dish, tx *gorm.DB) error {
	if tx != nil {
		// 如果存在事务
		return utils.AutoFill(dao.createWithTx)(ctx, tx, dish, constant.Create)
	}
	// 不存在事务
	return errs.New(constant.CodeInternalError, constant.MsgDatabaseTransactionFail)
}

func (dao *DishDAO) createWithTx(_ *gin.Context, tx *gorm.DB, dish any, _ string) error {
	result := tx.Create(dish)
	return result.Error
}

// PageQuery 分页查询
func (dao *DishDAO) PageQuery(name string, categoryId, status, page, pageSize int) ([]*vo.DishVO, int64, error) {
	var (
		dishVOs []*vo.DishVO
		total   int64
	)
	query := global.DB.Model(&entity.Dish{}).Joins("left join category on dish.category_id = category.id").Select("dish.*, category.name as category_name")

	if name != "" {
		// ☆ like 不能写成 "="
		query = query.Where("dish.name like ?", "%"+name+"%")
	}
	if categoryId > 0 {
		query = query.Where("dish.category_id = ?", categoryId)
	}
	if status == 0 || status == 1 {
		query = query.Where("dish.status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("dish.create_time desc").Offset(offset).Limit(pageSize).Find(&dishVOs).Error; err != nil {
		return nil, 0, err
	}

	return dishVOs, total, nil
}

// GetById 根据 id 查找菜品信息
func (dao *DishDAO) GetById(id int) (*entity.Dish, error) {
	var dish entity.Dish
	result := global.DB.Model(&entity.Dish{}).Where("id = ?", id).First(&dish)
	return &dish, result.Error
}

// DeleteByIDsTx 按 ids 批量删除菜品，事务版
func (dao *DishDAO) DeleteByIDsTx(ids []int, tx *gorm.DB) error {
	result := tx.Delete(&entity.Dish{}, ids) // 只有 id 是主键的时候才能这样用
	return result.Error
}

// UpdateTx 事务更新菜品基本信息
func (dao *DishDAO) UpdateTx(ctx *gin.Context, dish *entity.Dish, tx *gorm.DB) error {
	return utils.AutoFill(dao.updateTx)(ctx, tx, dish, constant.Update)
}

func (dao *DishDAO) updateTx(_ *gin.Context, tx *gorm.DB, dish any, _ string) error {
	d, ok := dish.(*entity.Dish)
	if !ok {
		return errs.New(constant.CodeInternalError, constant.MsgTypeConversionFail)
	}
	result := tx.Model(&entity.Dish{}).Where("id = ?", d.ID).Updates(d)
	return result.Error
}

// UpdateStatus 更新菜品状态
func (dao *DishDAO) UpdateStatus(id int, status int) error {
	return global.DB.Table("dish").Where("id = ?", id).UpdateColumn("status", status).Error
}

// GetByCategoryID 根据分类ID获取菜品列表
func (dao *DishDAO) GetByCategoryID(categoryID int) ([]*entity.Dish, error) {
	var list []*entity.Dish
	result := global.DB.Model(&entity.Dish{}).Where("category_id = ?", categoryID).Order("create_time desc").Find(&list)
	return list, result.Error
}

// CountHaltSales 查看所给的ID中有多少停售的套餐
func (dao *DishDAO) CountHaltSales(db *gorm.DB, ids []int) (int64, error) {
	var count int64
	result := db.Model(&entity.Dish{}).Where("id in ? and status = ?", ids, constant.DishDisable).Count(&count)
	return count, result.Error
}

// GetCount 获取起售或停售菜品数量
func (dao *DishDAO) GetCount(db *gorm.DB, status int) (int64, error) {
	var cnt int64
	err := db.Model(&entity.Dish{}).Where("status = ?", status).Count(&cnt).Error
	return cnt, err
}
