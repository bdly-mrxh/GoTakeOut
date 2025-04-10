package service

import (
	"github.com/gin-gonic/gin"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/internal/dao"
	"takeout/model/dto"
	"takeout/model/entity"
)

type ShoppingCartService struct {
	shoppingCartDAO dao.ShoppingCartDAO
	dishDAO         dao.DishDAO
	setmealDAO      dao.SetmealDAO
}

// Add 添加购物车
func (s *ShoppingCartService) Add(ctx *gin.Context, addDTO *dto.ShoppingCartDTO) error {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgGetIDFail)
	}
	cart := &entity.ShoppingCart{}
	if err = utils.CopyProperties(addDTO, cart); err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	cart.UserID = userID

	list, err := s.shoppingCartDAO.List(global.DB, cart)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	if len(list) > 0 {
		c := list[0]
		c.Number++
		err = s.shoppingCartDAO.UpdateNumberByID(global.DB, c)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
	} else {
		if addDTO.DishID != 0 {
			// 添加的是菜品
			dish, e := s.dishDAO.GetById(addDTO.DishID)
			if e != nil {
				return errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
			}
			cart.Name = dish.Name
			cart.Image = dish.Image
			cart.Amount = dish.Price
		} else {
			// 添加的是套餐
			setmeal, e := s.setmealDAO.GetByID(global.DB, addDTO.SetmealID)
			if e != nil {
				return errs.Wrap(e, constant.CodeDatabaseError, constant.MsgDatabaseError)
			}
			cart.Name = setmeal.Name
			cart.Image = setmeal.Image
			cart.Amount = setmeal.Price
		}
		cart.Number = 1
		err = s.shoppingCartDAO.Create(global.DB, cart)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
	}

	return nil
}

// List 查看购物车
func (s *ShoppingCartService) List(ctx *gin.Context) ([]*entity.ShoppingCart, error) {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return nil, err
	}
	list, err := s.shoppingCartDAO.List(global.DB, &entity.ShoppingCart{UserID: userID})
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return list, nil
}

// Clean 清空购物车
func (s *ShoppingCartService) Clean(ctx *gin.Context) error {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return err
	}
	err = s.shoppingCartDAO.CleanByUserID(global.DB, userID)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// Sub 减少购物车中的一个商品
func (s *ShoppingCartService) Sub(ctx *gin.Context, subDTO *dto.ShoppingCartDTO) error {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgGetIDFail)
	}
	cart := &entity.ShoppingCart{}
	if err = utils.CopyProperties(subDTO, cart); err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	cart.UserID = userID

	list, err := s.shoppingCartDAO.List(global.DB, cart)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	if len(list) == 0 {
		return errs.New(constant.CodeBusinessError, constant.MsgBadRequest)
	}
	c := list[0]
	if c.Number != 1 {
		c.Number--
		err = s.shoppingCartDAO.UpdateNumberByID(global.DB, c)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
	} else {
		err = s.shoppingCartDAO.DeleteByID(global.DB, c)
		if err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
	}
	return nil
}
