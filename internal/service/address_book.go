package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/internal/dao"
	"takeout/model/dto"
	"takeout/model/entity"
)

type AddressBookService struct {
	addressBookDAO dao.AddressBookDAO
}

// List 查询当前账户的所有地址信息
func (s *AddressBookService) List(ctx *gin.Context) ([]*entity.AddressBook, error) {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return nil, err
	}
	addressBook := &entity.AddressBook{UserID: userID, IsDefault: constant.NotSetAddress}

	list, err := s.addressBookDAO.List(global.DB, addressBook)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return list, nil
}

// Add 新增地址
func (s *AddressBookService) Add(ctx *gin.Context, addDTO *dto.AddressBookDTO) error {
	addressBook := &entity.AddressBook{}
	err := utils.CopyProperties(addDTO, addressBook)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}
	userID, err := utils.GetId(ctx)
	if err != nil {
		return err
	}
	addressBook.UserID = userID
	addressBook.IsDefault = 0
	err = s.addressBookDAO.Add(global.DB, addressBook)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// GetByID 根据ID获取地址详细信息
func (s *AddressBookService) GetByID(id int) (*entity.AddressBook, error) {
	address, err := s.addressBookDAO.GetByID(global.DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(err, constant.CodeBadRequest, constant.MsgBadRequest)
		}
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return address, nil
}

// Update 跟新地址信息
func (s *AddressBookService) Update(addressDTO *dto.AddressBookDTO) error {
	address := &entity.AddressBook{}
	err := utils.CopyProperties(addressDTO, address)
	if err != nil {
		return errs.Wrap(err, constant.CodeInternalError, constant.MsgCopyPropertiesFail)
	}

	err = s.addressBookDAO.Update(global.DB, address)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// SetDefault 设置默认地址
func (s *AddressBookService) SetDefault(ctx *gin.Context, addressDTO *dto.AddressBookDTO) error {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return err
	}
	return global.DB.Transaction(func(db *gorm.DB) error {
		// 将该用户的所有地址置为非默认
		if err = s.addressBookDAO.SetNonDefault(db, userID); err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		// 将对应地址ID设置为默认地址
		if err = s.addressBookDAO.SetDefault(db, addressDTO.ID); err != nil {
			return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
		return nil
	})
}

// DeleteByID 根据ID删除地址
func (s *AddressBookService) DeleteByID(id int) error {
	err := s.addressBookDAO.DeleteByID(global.DB, id)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDeleteSuccess)
	}
	return nil
}

// GetDefault 查询默认地址
func (s *AddressBookService) GetDefault(ctx *gin.Context) (*entity.AddressBook, error) {
	userID, err := utils.GetId(ctx)
	if err != nil {
		return nil, err
	}
	list, err := s.addressBookDAO.List(global.DB, &entity.AddressBook{UserID: userID, IsDefault: constant.DefaultAddress})
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	if len(list) == 0 {
		return nil, errs.New(constant.CodeBusinessError, constant.MsgNotExistDefaultAddress)
	}
	return list[0], nil
}
