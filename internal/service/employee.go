package service

import (
	"errors"
	"strconv"
	"strings"
	"takeout/internal/dao"

	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/utils"
	"takeout/model/dto"
	"takeout/model/entity"
	"takeout/model/vo"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmployeeService 员工服务
type EmployeeService struct {
	employeeDAO dao.EmployeeDAO
}

// Login 员工登录
func (s *EmployeeService) Login(loginDTO *dto.EmployeeLoginDTO) (*vo.EmployeeLoginVO, error) {
	// 根据用户名查询员工
	employee, err := s.employeeDAO.GetByUsername(loginDTO.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(constant.CodeUserNotExist, constant.MsgUserNotExist)
		}
		return nil, errs.Wrap(err, constant.CodeDatabaseError, "查询员工信息失败")
	}

	// 检查员工状态
	if employee.Status == 0 {
		return nil, errs.New(constant.CodeUserDisabled, "账号已禁用")
	}

	// 密码加密比对
	md5PasswordStr := utils.Encrypt(loginDTO.Password)

	if employee.Password != md5PasswordStr {
		return nil, errs.New(constant.CodePasswordError, constant.MsgPasswordIncorrect)
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(constant.EmpID, strconv.Itoa(employee.ID))
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeServerError, "生成令牌失败")
	}

	// 构建返回数据
	loginVO := &vo.EmployeeLoginVO{
		ID:       employee.ID,
		Username: employee.Username,
		Name:     employee.Name,
		Token:    token,
	}

	return loginVO, nil
}

// Create 创建新员工
func (s *EmployeeService) Create(ctx *gin.Context, createDTO *dto.EmployeeCreateDTO) error {
	// 检查用户名是否已存在
	exists, err := s.employeeDAO.CheckUsernameExists(createDTO.Username)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, "检查用户名是否存在失败")
	}
	if exists {
		return errs.New(constant.CodeEmployeeCreateFail, "用户名已存在")
	}

	// 获取当前操作用户ID
	//id, err := utils.GetId(ctx)
	//if err != nil {
	//	return err
	//}

	// 创建员工实体
	employee := &entity.Employee{
		Status: constant.EmployeeStatusEnable, // 默认启用状态
		//CreateUser: id,
		//UpdateUser: id,
	}

	// 使用deepcopier复制DTO到实体
	if err = utils.CopyProperties(createDTO, employee); err != nil {
		return errs.Wrap(err, constant.CodeEmployeeCreateFail, constant.MsgCopyPropertiesFail)
	}

	// 设置默认密码并加密
	employee.Password = utils.Encrypt(constant.DefaultPassword)

	// 保存到数据库
	err = s.employeeDAO.Create(ctx, employee)
	if err != nil {
		var myErr *errs.Error
		if errors.As(err, &myErr) {
			return myErr
		}
		if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
			return errs.Wrap(err, constant.CodeCategoryCreateFail, constant.MsgNameConflict)
		}
		return errs.Wrap(err, constant.CodeEmployeeCreateFail, constant.MsgEmployeeCreateFail)
	}

	return nil
}

// GetById 根据ID查询员工信息
func (s *EmployeeService) GetById(id int) (*vo.EmployeeDetailVO, error) {
	// 根据ID查询员工
	employee, err := s.employeeDAO.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(constant.CodeUserNotExist, constant.MsgUserNotExist)
		}
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}

	// 转换为VO并返回
	employeeDetailVO, err := vo.EmployeeDetailVOFromEntity(employee)
	if err != nil {

	}
	return employeeDetailVO, nil
}

// UpdateStatusById 更新单个员工状态
func (s *EmployeeService) UpdateStatusById(status int, id int) error {
	// Update employee status directly in one operation
	err := s.employeeDAO.UpdateStatus(id, status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(constant.CodeUserNotExist, "员工ID不存在: "+strconv.Itoa(id))
		}
		return errs.Wrap(err, constant.CodeDatabaseError, "更新员工状态失败")
	}

	return nil
}

// UpdatePassword 修改员工密码
func (s *EmployeeService) UpdatePassword(ctx *gin.Context, passwordDTO *dto.EmployeePasswordDTO) error {
	// 根据ID查询员工
	employee, err := s.employeeDAO.GetById(passwordDTO.EmpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(constant.CodeUserNotExist, "员工ID不存在: "+strconv.Itoa(passwordDTO.EmpId))
		}
		return errs.Wrap(err, constant.CodeDatabaseError, "查询员工信息失败")
	}

	// 验证旧密码
	md5OldPasswordStr := utils.Encrypt(passwordDTO.OldPassword)

	if employee.Password != md5OldPasswordStr {
		return errs.New(constant.CodePasswordError, "原密码错误")
	}

	// 获取当前操作用户ID
	//id, err := utils.GetId(ctx)
	//if err != nil {
	//	return err
	//}

	// 加密新密码
	md5NewPasswordStr := utils.Encrypt(passwordDTO.NewPassword)

	// 更新密码
	employee.Password = md5NewPasswordStr
	// employee.UpdateUser = id

	// 保存到数据库
	err = s.employeeDAO.UpdateById(ctx, employee)
	if err != nil {
		var myErr *errs.Error
		if errors.As(err, &myErr) {
			return myErr
		}
		return errs.Wrap(err, constant.CodeDatabaseError, "修改密码失败")
	}

	return nil
}

// Update 更新员工信息
func (s *EmployeeService) Update(ctx *gin.Context, updateDTO *dto.EmployeeUpdateDTO) error {
	// 根据ID查询员工
	employee, err := s.employeeDAO.GetById(updateDTO.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(constant.CodeUserNotExist, constant.MsgUserNotExist)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, "查询员工信息失败")
	}

	// 获取当前操作用户ID
	//id, err := utils.GetId(ctx)
	//if err != nil {
	//	return err
	//}
	// 更新员工信息
	// 使用deepcopier复制DTO到实体
	err = utils.CopyProperties(updateDTO, employee)
	if err != nil {
		return errs.Wrap(err, constant.CodeEmployeeUpdateFail, constant.MsgCopyPropertiesFail)
	}
	// employee.UpdateUser = id

	// 保存到数据库
	err = s.employeeDAO.UpdateById(ctx, employee)
	if err != nil {
		var myErr *errs.Error
		if errors.As(err, &myErr) {
			return myErr
		}
		// Check if error is due to duplicate username
		if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
			return errs.Wrap(err, constant.CodeCategoryCreateFail, constant.MsgNameConflict)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, "更新员工信息失败")
	}

	return nil
}

// PageQuery 分页查询员工信息
func (s *EmployeeService) PageQuery(pageDTO *dto.EmployeePageDTO) (*vo.PageResult, error) {
	// 调用DAO层进行分页查询
	employees, total, err := s.employeeDAO.PageQuery(pageDTO.Name, pageDTO.Page, pageDTO.PageSize)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgPageQueryEmployeeFail)
	}

	// 转换为VO列表
	employeeVOs := make([]*vo.EmployeeDetailVO, 0, len(employees))
	for _, employee := range employees {
		var employeeDetail *vo.EmployeeDetailVO
		employeeDetail, err = vo.EmployeeDetailVOFromEntity(&employee)
		if err != nil {
			return nil, errs.Wrap(err, constant.CodeEmployeePageQueryFail, constant.MsgCopyPropertiesFail)
		}
		employeeVOs = append(employeeVOs, employeeDetail)
	}

	// 构建分页结果
	pageResult := &vo.PageResult{
		Total:   total,
		Records: employeeVOs,
	}

	return pageResult, nil
}
