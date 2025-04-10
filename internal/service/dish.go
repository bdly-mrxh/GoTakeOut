package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/internal/dao"
	"takeout/model/dto"
	"takeout/model/entity"
	"takeout/model/vo"
)

// DishService 菜品服务
type DishService struct {
	dishDAO        dao.DishDAO
	dishFlavorDAO  dao.DishFlavorDAO
	setmealDishDAO dao.SetmealDishDAO
}

// CreateWithFlavors 创建菜品及其口味（事务操作）
func (s *DishService) CreateWithFlavors(ctx *gin.Context, createDTO *dto.DishDTO) error {
	dish := &entity.Dish{}
	// ☆ deepcopier.Copy(from).To(to): 需要传入非 nil 指针
	err := utils.CopyProperties(createDTO, dish)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	// dish.Status = constant.DefaultStatus
	// 开启事务
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err = s.dishDAO.CreateWithTx(ctx, dish, tx); err != nil {
			if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
				return errs.Wrap(err, constant.CodeBadRequest, constant.MsgNameConflict)
			}
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgCreateFail)
		}

		// 如果存在口味信息，批量创建口味
		if len(createDTO.Flavors) > 0 {
			// 设置菜品ID
			for i := range createDTO.Flavors {
				createDTO.Flavors[i].DishID = dish.ID
			}
			// 批量创建口味
			if err = s.dishFlavorDAO.BatchCreateWithTx(createDTO.Flavors, tx); err != nil {
				return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgCreateFail)
			}
		}

		err = utils.CleanCache(constant.CacheDishKey + strconv.Itoa(createDTO.CategoryID))
		if err != nil {
			return errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
		}

		return nil
	})
}

// PageQuery 分页查询菜品
func (s *DishService) PageQuery(queryDTO *dto.DishPageQueryDTO) (*vo.PageResult, error) {
	dishVOs, total, err := s.dishDAO.PageQuery(queryDTO.Name, queryDTO.CategoryID, queryDTO.Status, queryDTO.Page, queryDTO.PageSize)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	//for _, dishVO := range dishVOs {
	//	dishVO.Flavors = make([]*entity.DishFlavor, 0)
	//}
	return &vo.PageResult{
		Total:   total,
		Records: dishVOs,
	}, nil
}

// Delete 删除菜品，可以多删除
func (s *DishService) Delete(ids []int) error {
	// 首先判断能不能删除
	// 1.判断菜品是否在售
	for _, id := range ids {
		dish, err := s.dishDAO.GetById(id)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
		}
		if dish.Status == constant.DishEnable {
			return errs.New(constant.CodeBusinessError, constant.MsgDishOnSale)
		}
	}
	// 2.判断菜品是否与套餐关联
	count, err := s.setmealDishDAO.CountByDishIDs(ids)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	if count > 0 {
		return errs.New(constant.CodeBusinessError, constant.MsgDishAssociativeWithSetmeal)
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		err = s.dishDAO.DeleteByIDsTx(ids, tx)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDeleteFail)
		}
		err = s.dishFlavorDAO.DeleteByDishIDsTx(ids, tx)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDeleteFail)
		}

		err = utils.CleanCache(constant.CacheDishKey + "*")
		if err != nil {
			return errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
		}

		return nil
	})
}

// GetByID 根据 ID 查询菜品信息
func (s *DishService) GetByID(id int) (*vo.DishVO, error) {
	// 查询菜品基本信息
	dish, err := s.dishDAO.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(constant.CodeBadRequest, constant.MsgNotFound)
		}
		return nil, errs.New(constant.CodeDatabaseError, constant.MsgQueryFail)
	}

	dishVO := &vo.DishVO{}
	err = utils.CopyProperties(dish, dishVO)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}

	// 查询口味数据
	flavors, err := s.dishFlavorDAO.GetByDishID(id)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}
	dishVO.Flavors = flavors
	return dishVO, nil
}

// Update 更新菜品信息
func (s *DishService) Update(ctx *gin.Context, dishDTO *dto.DishDTO) error {
	dish := &entity.Dish{}
	if err := utils.CopyProperties(dishDTO, dish); err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	// 开启事务
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 更新菜品基本信息
		if err := s.dishDAO.UpdateTx(ctx, dish, tx); err != nil {
			if strings.Contains(err.Error(), constant.MsgKeyDuplicateError) {
				return errs.Wrap(err, constant.CodeBusinessError, constant.MsgNameConflict)
			}
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgUpdateFail)
		}
		// 删除口味数据
		if err := s.dishFlavorDAO.DeleteByDishIDTx(dishDTO.ID, tx); err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDeleteFail)
		}
		// 新增口味数据
		// 如果存在口味信息，批量创建口味
		if len(dishDTO.Flavors) > 0 {
			// 设置菜品ID
			for i := range dishDTO.Flavors {
				dishDTO.Flavors[i].DishID = dish.ID
			}
			// 批量创建口味
			if err := s.dishFlavorDAO.BatchCreateWithTx(dishDTO.Flavors, tx); err != nil {
				return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgCreateFail)
			}
		}
		// 修改菜品的话，缓存也是要全部删除，可能涉及多个分类
		err := utils.CleanCache(constant.CacheDishKey + "*")
		if err != nil {
			return errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
		}

		return nil
	})
}

// UpdateStatus 更新菜品状态
func (s *DishService) UpdateStatus(id int, status int) error {
	err := s.dishDAO.UpdateStatus(id, status)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgUpdateFail)
	}
	// TODO: 可以改为取出dish，清除对应categoryID的缓存
	err = utils.CleanCache(constant.CacheDishKey + "*")
	if err != nil {
		return errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
	}

	return nil
}

// ListByCategoryID 根据分类ID查询菜品
func (s *DishService) ListByCategoryID(categoryID int) ([]*vo.DishVO, error) {
	// 如果在缓存中，直接返回
	list := make([]*vo.DishVO, 0)
	ctx := context.Background()
	cacheData, err := global.Redis.Get(ctx, constant.CacheDishKey+strconv.Itoa(categoryID)).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
		}
	} else {
		if err = json.Unmarshal([]byte(cacheData), &list); err != nil {
			return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
		}
		// logger.Infof("缓存命中，list : %v", list)
		return list, nil
	}

	// 缓存未命中
	// logger.Infof("缓存未命中，categoryID : %v", categoryID)
	dishes, err := s.dishDAO.GetByCategoryID(categoryID)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgQueryFail)
	}

	for _, dish := range dishes {
		var dishVO vo.DishVO
		if err = utils.CopyProperties(dish, &dishVO); err != nil {
			return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
		}
		flavors, e := s.dishFlavorDAO.GetByDishID(dish.ID)
		if e != nil {
			return nil, errs.Wrap(e, constant.CodeDatabaseError, constant.MsgQueryFail)
		}
		dishVO.Flavors = flavors
		list = append(list, &dishVO)
	}

	// 数据存入缓存
	listJson, err := json.Marshal(list)
	if err != nil {
		// 其实只要WARN一下错误即可
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgMarshalFail)
	}
	err = global.Redis.Set(ctx, constant.CacheDishKey+strconv.Itoa(categoryID), string(listJson), 0).Err()
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
	}

	return list, nil
}
