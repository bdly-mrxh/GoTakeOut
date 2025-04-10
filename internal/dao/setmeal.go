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

// SetmealDAO 套餐数据访问对象
type SetmealDAO struct{}

// CountByCategoryID 根据分类ID统计套餐数量
func (dao *SetmealDAO) CountByCategoryID(categoryID int) (int64, error) {
	var count int64
	result := global.DB.Model(&entity.Setmeal{}).Where("category_id = ?", categoryID).Count(&count)
	return count, result.Error
}

// Create 新增套餐
func (dao *SetmealDAO) Create(ctx *gin.Context, tx *gorm.DB, setmeal *entity.Setmeal) error {
	return utils.AutoFill(dao.create)(ctx, tx, setmeal, constant.Create)
}

func (dao *SetmealDAO) create(_ *gin.Context, tx *gorm.DB, setmeal any, _ string) error {
	s, ok := setmeal.(*entity.Setmeal)
	if !ok {
		return errs.New(constant.CodeInternalError, constant.MsgTypeConversionFail)
	}
	return tx.Model(&entity.Setmeal{}).Create(&s).Error
}

// GetByID 获取套餐详细信息
func (dao *SetmealDAO) GetByID(db *gorm.DB, id int) (*entity.Setmeal, error) {
	var setmeal entity.Setmeal
	result := db.Model(&entity.Setmeal{}).Where("id = ?", id).First(&setmeal)
	return &setmeal, result.Error
}

// CountOnSaleSetmealByIDs 查看IDs中在售套餐的个数
func (dao *SetmealDAO) CountOnSaleSetmealByIDs(db *gorm.DB, ids []int) (int64, error) {
	var count int64
	result := db.Model(&entity.Setmeal{}).Where("id in ? and status = ?", ids, constant.SetmealEnable).Count(&count)
	return count, result.Error
}

// BatchDelete 批量删除
func (dao *SetmealDAO) BatchDelete(db *gorm.DB, ids []int) error {
	result := db.Where("id = ?", ids).Delete(&entity.Setmeal{})
	return result.Error
}

// PageQuery 分页查询
func (dao *SetmealDAO) PageQuery(db *gorm.DB, name string, categoryId int, status int, page int, size int) (int64, []*vo.SetmealVO, error) {
	var (
		setmealVOs []*vo.SetmealVO
		total      int64
		err        error
	)

	query := db.Table("setmeal s").
		Joins("left join category c on s.category_id = c.id").
		Select("s.*, c.name category_name")
	if name != "" {
		query = query.Where("s.name like ?", "%"+name+"%")
	}
	if categoryId != 0 {
		query = query.Where("s.category_id = ?", categoryId)
	}
	if status == 0 || status == 1 {
		query = query.Where("s.status = ?", status)
	}

	if err = query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	offset := (page - 1) * size
	if err = query.Offset(offset).Limit(size).Order("s.create_time desc").Find(&setmealVOs).Error; err != nil {
		return 0, nil, err
	}
	return total, setmealVOs, nil
}

// Update 更新套餐信息
func (dao *SetmealDAO) Update(ctx *gin.Context, db *gorm.DB, setmeal *entity.Setmeal) error {
	return utils.AutoFill(dao.update)(ctx, db, setmeal, constant.Update)
}

func (dao *SetmealDAO) update(_ *gin.Context, db *gorm.DB, setmeal any, _ string) error {
	s, ok := setmeal.(*entity.Setmeal)
	if !ok {
		return errs.New(constant.CodeInternalError, constant.MsgTypeConversionFail)
	}
	result := db.Model(&entity.Setmeal{}).Where("id = ?", s.ID).Updates(&s)
	return result.Error
}

// UpdateStatus 更新套餐状态
func (dao *SetmealDAO) UpdateStatus(db *gorm.DB, id int, status int) error {
	result := db.Model(&entity.Setmeal{}).Where("id = ?", id).UpdateColumn("status", status)
	return result.Error
}

// ListByCategoryID 根据分类ID查询起售的套餐列表
func (dao *SetmealDAO) ListByCategoryID(db *gorm.DB, categoryID int) ([]*entity.Setmeal, error) {
	var setmeals []*entity.Setmeal
	result := db.Model(&entity.Setmeal{}).Where("category_id = ? and status = ?", categoryID, constant.SetmealEnable).Find(&setmeals)
	return setmeals, result.Error
}

// GetDishItemBySetmealID 根据套餐id查询菜品选项
func (dao *SetmealDAO) GetDishItemBySetmealID(db *gorm.DB, setmealId int) ([]*vo.DishItem, error) {
	var dishItems []*vo.DishItem
	result := db.Table("setmeal_dish sd").
		Joins("left join dish d on sd.dish_id = d.id").
		Where("sd.setmeal_id = ?", setmealId).
		Select("sd.name, sd.copies, d.image, d.description").
		Find(&dishItems)
	return dishItems, result.Error
}

// GetCount 获取起售停售套餐数量
func (dao *SetmealDAO) GetCount(db *gorm.DB, status int) (int64, error) {
	var cnt int64
	err := db.Model(&entity.Setmeal{}).Where("status = ?", status).Count(&cnt).Error
	return cnt, err
}
