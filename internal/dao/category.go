package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/model/entity"
)

// CategoryDAO 分类数据访问对象
type CategoryDAO struct{}

// Create 创建分类
func (dao *CategoryDAO) Create(ctx *gin.Context, category *entity.Category) error {
	// return utils.AutoFill(dao.create)(ctx, nil, category, constant.Create)
	return utils.AutoFill(func(_ *gin.Context, _ *gorm.DB, _ any, _ string) error {
		return dao.create(category) // 传递的是指针，可以这样做
	})(ctx, nil, category, constant.Create)
}

func (dao *CategoryDAO) create(category any) error {
	result := global.DB.Model(&entity.Category{}).Create(category)
	return result.Error
}

// Update 更新分类
func (dao *CategoryDAO) Update(ctx *gin.Context, category *entity.Category) error {
	return utils.AutoFill(dao.update)(ctx, nil, category, constant.Update)
}

func (dao *CategoryDAO) update(_ *gin.Context, _ *gorm.DB, category any, _ string) error {
	if c, ok := category.(*entity.Category); ok {
		result := global.DB.Model(&entity.Category{}).Where("id = ?", c.ID).Updates(c)
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return result.Error
	}
	return errs.New(constant.CodeInternalError, constant.MsgTypeConversionFail)
}

func (dao *CategoryDAO) UpdateStatus(id, status int) error {
	result := global.DB.Model(&entity.Category{}).Where("id = ?", id).UpdateColumn("status", status)
	// Update 操作 id 不存在 不会报错，靠影响的行数来判断
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (dao *CategoryDAO) PageQuery(name string, typeId, page, pageSize int) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64

	// 构建查询条件
	query := global.DB.Model(&entity.Category{})

	// 如果提供了名称，添加模糊查询
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 如果指定了类型，按类型查询
	if typeId > 0 {
		query = query.Where("type = ?", typeId)
	}

	// 查询记录总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询的数据
	offset := (page - 1) * pageSize
	if err := query.Order("sort asc, create_time desc").Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// List 按类型查询
func (dao *CategoryDAO) List(typeId int) ([]*entity.Category, error) {
	var categories []*entity.Category
	query := global.DB.Model(entity.Category{})
	if typeId != 0 {
		query = query.Where("type = ?", typeId)
	}
	result := query.Order("sort ASC, create_time DESC").Find(&categories)
	return categories, result.Error
}

// Delete 删除分类
func (dao *CategoryDAO) Delete(id int) error {
	result := global.DB.Delete(&entity.Category{}, id)
	return result.Error
}
