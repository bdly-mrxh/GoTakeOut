package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/utils"
	"takeout/internal/dao"
	"takeout/model/dto"
	"takeout/model/entity"
	"takeout/model/vo"
)

// CategoryService 分类服务
type CategoryService struct {
	categoryDAO dao.CategoryDAO
	dishDAO     dao.DishDAO
	setmealDAO  dao.SetmealDAO
}

// Create 创建分类
func (s *CategoryService) Create(ctx *gin.Context, createDTO *dto.CategoryDTO) error {
	// 前端前来的 sort, type 是 string，要转化为 int
	sort, err := strconv.Atoi(createDTO.Sort)
	if err != nil {
		return errs.New(constant.CodeBadRequest, constant.MsgTypeConversionFail)
	}
	typeId, err := strconv.Atoi(createDTO.Type)
	if err != nil {
		return errs.New(constant.CodeBadRequest, constant.MsgTypeConversionFail)
	}
	// 创建分类实体
	category := &entity.Category{
		Type:   typeId,
		Name:   createDTO.Name,
		Sort:   sort,
		Status: constant.DefaultStatus,
	}

	// 保存分类
	if err = s.categoryDAO.Create(ctx, category); err != nil {
		var myErr *errs.Error
		if errors.As(err, &myErr) {
			return myErr
		}
		if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
			return errs.Wrap(err, constant.CodeEmployeeUpdateFail, "菜名已存在")
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}

	return nil
}

// UpdateStatus 更新分类状态
func (s *CategoryService) UpdateStatus(id, status int) error {
	err := s.categoryDAO.UpdateStatus(id, status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.Wrap(err, constant.CodeNotFound, constant.MsgNotFound)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Update 更新分类信息
func (s *CategoryService) Update(ctx *gin.Context, updateDTO *dto.CategoryDTO) error {
	category := &entity.Category{}
	// ☆ error: 字段名哪怕一点儿不一样也不会拷贝，ID 和 Id 就无法拷贝，甚至连 int64 到 int 都无法拷贝
	err := utils.CopyProperties(updateDTO, category)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	// 前端没定义好
	sort, err := strconv.Atoi(updateDTO.Sort)
	if err != nil {
		return errs.New(constant.CodeBadRequest, constant.MsgTypeConversionFail)
	}
	// type 可能没设置
	if updateDTO.Type != "" {
		var typeId int
		typeId, err = strconv.Atoi(updateDTO.Type)
		if err != nil {
			return errs.New(constant.CodeBadRequest, constant.MsgTypeConversionFail)
		}
		category.Type = typeId
	}
	category.Sort = sort

	err = s.categoryDAO.Update(ctx, category)
	if err != nil {
		var myErr *errs.Error
		if errors.As(err, &myErr) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.Wrap(err, constant.CodeNotFound, constant.MsgNotFound+"id: "+strconv.Itoa(updateDTO.ID))
		}
		if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
			return errs.Wrap(err, constant.CodeCategoryCreateFail, constant.MsgNameConflict)
		}
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// PageQuery 分类分页查询
func (s *CategoryService) PageQuery(queryDTO *dto.CategoryPageDTO) (*vo.PageResult, error) {
	categories, total, err := s.categoryDAO.PageQuery(queryDTO.Name, queryDTO.Type, queryDTO.Page, queryDTO.PageSize)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}

	page := &vo.PageResult{
		Total:   total,
		Records: categories,
	}

	return page, nil
}

// List 按类型查询分类
func (s *CategoryService) List(typeId int) ([]*entity.Category, error) {
	categories, err := s.categoryDAO.List(typeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(err, constant.CodeNotFound, constant.MsgNotFound)
		}
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	return categories, nil
}

// Delete 删除分类
func (s *CategoryService) Delete(id int) error {
	// 删除的分类不能有任何关联菜品和套餐
	count, err := s.dishDAO.CountByCategoryID(id)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	if count > 0 {
		return errs.New(constant.CodeDeleteCategoryFail, constant.MsgExistAssociativeDishOrSetmeal)
	}

	count, err = s.setmealDAO.CountByCategoryID(id)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	if count > 0 {
		return errs.New(constant.CodeDeleteCategoryFail, constant.MsgExistAssociativeDishOrSetmeal)
	}

	// 删除分类
	err = s.categoryDAO.Delete(id)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDeleteFail)
	}
	return nil
}
