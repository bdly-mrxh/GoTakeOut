package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"takeout/common/constant"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/model/entity"
)

// EmployeeDAO 员工数据访问对象
type EmployeeDAO struct{}

// GetByUsername 根据用户名查询员工信息
func (dao *EmployeeDAO) GetByUsername(username string) (*entity.Employee, error) {
	var employee entity.Employee
	result := global.DB.Where("username = ?", username).First(&employee)
	return &employee, result.Error
}

// UpdateStatus 根据ID直接更新员工状态
func (dao *EmployeeDAO) UpdateStatus(id int, status int) error {
	result := global.DB.Model(&entity.Employee{}).Where("id = ?", id).UpdateColumn("status", status)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// UpdateById 根据ID更新员工信息
func (dao *EmployeeDAO) UpdateById(ctx *gin.Context, employee *entity.Employee) error {
	return utils.AutoFill(dao.updateById)(ctx, nil, employee, constant.Update)
}

func (dao *EmployeeDAO) updateById(_ *gin.Context, _ *gorm.DB, employee any, _ string) error {
	result := global.DB.Updates(employee)
	return result.Error
}

// Create 新增员工
func (dao *EmployeeDAO) Create(ctx *gin.Context, employee *entity.Employee) error {
	return utils.AutoFill(dao.create)(ctx, nil, employee, constant.Create)
}

// 真正的创建新员工操作
func (dao *EmployeeDAO) create(_ *gin.Context, _ *gorm.DB, employee any, _ string) error {
	result := global.DB.Create(employee)
	return result.Error
}

// CheckUsernameExists 检查用户名是否已存在
func (dao *EmployeeDAO) CheckUsernameExists(username string) (bool, error) {
	var count int64
	result := global.DB.Model(&entity.Employee{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// GetById 根据ID查询员工信息
func (dao *EmployeeDAO) GetById(id int) (*entity.Employee, error) {
	var employee entity.Employee
	result := global.DB.Where("id = ?", id).First(&employee)
	return &employee, result.Error
}

// PageQuery 分页查询员工信息
func (dao *EmployeeDAO) PageQuery(name string, page, pageSize int) ([]entity.Employee, int64, error) {
	var employees []entity.Employee
	var total int64

	// 构建查询条件
	query := global.DB.Model(&entity.Employee{})

	// 如果提供了姓名，添加模糊查询条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 排序
	if err := query.Order("create_time desc").Error; err != nil {
		return nil, 0, err
	}

	// 查询总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}
