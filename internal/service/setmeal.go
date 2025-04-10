package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"takeout/common/aop"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/internal/dao"
	"takeout/model/dto"
	"takeout/model/entity"
	"takeout/model/vo"
)

// SetmealService 套餐服务
type SetmealService struct {
	setmealDAO     dao.SetmealDAO
	setmealDishDAO dao.SetmealDishDAO
	dishDAO        dao.DishDAO
}

// Create 新增套餐
func (s *SetmealService) Create(ctx *gin.Context, createDTO *dto.SetmealDTO) error {
	return aop.CacheEvict(func() error {
		return s.create(ctx, createDTO)
	}, &aop.CacheOptions{
		CacheName: constant.CacheSetmealKey,
		Key:       createDTO.CategoryID,
	})()
}

func (s *SetmealService) create(ctx *gin.Context, createDTO *dto.SetmealDTO) error {
	setmeal := &entity.Setmeal{}
	if err := utils.CopyProperties(createDTO, setmeal); err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.setmealDAO.Create(ctx, tx, setmeal); err != nil {
			if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
				return errs.Wrap(err, constant.CodeBusinessError, constant.MsgNameConflict)
			}
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}

		if createDTO.SetmealDishes != nil {
			for _, setmealDish := range createDTO.SetmealDishes {
				setmealDish.SetmealID = setmeal.ID
			}
			// 插入套餐菜品关系
			if err := s.setmealDishDAO.BatchInsert(tx, createDTO.SetmealDishes); err != nil {
				return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
			}
		}

		return nil
	})
}

// GetByID 根据ID查询套餐的详细信息
func (s *SetmealService) GetByID(id int) (*vo.SetmealVO, error) {
	setmeal, err := s.setmealDAO.GetByID(global.DB, id)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	setmealVO := &vo.SetmealVO{}
	if err = utils.CopyProperties(setmeal, setmealVO); err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	setmealDishes, err := s.setmealDishDAO.GetBySetmealID(global.DB, id)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	setmealVO.SetmealDishes = setmealDishes
	return setmealVO, nil
}

// BatchDelete 批量删除套餐，将查询操作也放到事务中，保证原子性
func (s *SetmealService) BatchDelete(ids []int) error {
	return aop.CacheEvict(func() error {
		return s.batchDelete(ids)
	}, &aop.CacheOptions{
		CacheName:  constant.CacheSetmealKey,
		AllEntries: true,
	})()
}

func (s *SetmealService) batchDelete(ids []int) error {
	if len(ids) == 0 {
		return errs.New(constant.CodeBusinessError, constant.MsgMissingRequest)
	}

	return global.DB.Transaction(func(db *gorm.DB) error {
		count, err := s.setmealDAO.CountOnSaleSetmealByIDs(db, ids)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
		}
		if count > 0 {
			return errs.New(constant.CodeBusinessError, constant.MsgSetmealOnSale)
		}

		err = s.setmealDAO.BatchDelete(db, ids)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDeleteFail)
		}

		err = s.setmealDishDAO.BatchDeleteBySetmealIDs(db, ids)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDeleteFail)
		}

		return nil
	})
}

// PageQuery 分页查询套餐
func (s *SetmealService) PageQuery(queryDTO *dto.SetmealPageQueryDTO) (*vo.PageResult, error) {
	total, page, err := s.setmealDAO.PageQuery(global.DB, queryDTO.Name, queryDTO.CategoryID, queryDTO.Status, queryDTO.Page, queryDTO.PageSize)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	return &vo.PageResult{
		Total:   total,
		Records: page,
	}, nil
}

// Update 更新套餐
func (s *SetmealService) Update(ctx *gin.Context, updateDTO *dto.SetmealDTO) error {
	return aop.CacheEvict(func() error {
		return s.update(ctx, updateDTO)
	}, &aop.CacheOptions{
		CacheName:  constant.CacheSetmealKey,
		AllEntries: true,
	})()
}

func (s *SetmealService) update(ctx *gin.Context, updateDTO *dto.SetmealDTO) error {
	var setmeal entity.Setmeal
	if err := utils.CopyProperties(updateDTO, &setmeal); err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	return global.DB.Transaction(func(db *gorm.DB) error {
		if err := s.setmealDAO.Update(ctx, db, &setmeal); err != nil {
			if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
				return errs.Wrap(err, constant.CodeBusinessError, constant.MsgNameConflict)
			}
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgUpdateFail)
		}

		if err := s.setmealDishDAO.BatchDeleteBySetmealID(db, updateDTO.ID); err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgUpdateFail)
		}

		if err := s.setmealDishDAO.BatchInsert(db, updateDTO.SetmealDishes); err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgUpdateFail)
		}

		return nil
	})
}

// UpdateStatus 更新套餐状态
func (s *SetmealService) UpdateStatus(id int, status int) error {
	return aop.CacheEvict(func() error {
		return s.updateStatus(id, status)
	}, &aop.CacheOptions{
		CacheName:  constant.CacheSetmealKey,
		AllEntries: true,
	})()
}

func (s *SetmealService) updateStatus(id int, status int) error {
	return global.DB.Transaction(func(db *gorm.DB) error {
		// 停售的套餐要开启需要所有关联菜品均起售
		if status == constant.SetmealEnable {
			dishIds, err := s.setmealDishDAO.GetDishIdsBySetmealId(db, id)
			if err != nil {
				return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
			}
			count, err := s.dishDAO.CountHaltSales(db, dishIds)
			if err != nil {
				return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
			}
			if count > 0 {
				return errs.New(constant.CodeBusinessError, constant.MsgSetmealAssociativeDishHalfSales)
			}
		}

		err := s.setmealDAO.UpdateStatus(db, id, status)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgUpdateFail)
		}
		return nil
	})
}

// ListByCategoryID 根据分类ID查询套餐列表
func (s *SetmealService) ListByCategoryID(categoryID int) ([]*entity.Setmeal, error) {
	return aop.Cacheable(func() ([]*entity.Setmeal, error) {
		return s.listByCategoryID(categoryID)
	}, &aop.CacheOptions{
		CacheName: constant.CacheSetmealKey,
		Key:       categoryID,
	})()
}

func (s *SetmealService) listByCategoryID(categoryID int) ([]*entity.Setmeal, error) {
	list, err := s.setmealDAO.ListByCategoryID(global.DB, categoryID)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	return list, nil
}

// GetDishItemBySetmealID 根据套餐ID查询DishItems
func (s *SetmealService) GetDishItemBySetmealID(id int) ([]*vo.DishItem, error) {
	dishItems, err := s.setmealDAO.GetDishItemBySetmealID(global.DB, id)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return dishItems, nil
}
