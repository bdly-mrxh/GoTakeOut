package dao

import (
	"gorm.io/gorm"
	"takeout/common/constant"
	"takeout/model/entity"
)

type AddressBookDAO struct{}

// List 动态查询地址
func (d *AddressBookDAO) List(db *gorm.DB, book *entity.AddressBook) ([]*entity.AddressBook, error) {
	query := db.Model(&entity.AddressBook{})
	if book.UserID != 0 {
		query = query.Where("user_id = ?", book.UserID)
	}
	if book.Phone != "" {
		query = query.Where("phone = ?", book.Phone)
	}
	if book.IsDefault != constant.NotSetAddress {
		query = query.Where("is_default = ?", book.IsDefault)
	}

	var list []*entity.AddressBook
	result := query.Find(&list)
	return list, result.Error
}

// Add 新增地址
func (d *AddressBookDAO) Add(db *gorm.DB, book *entity.AddressBook) error {
	return db.Model(&entity.AddressBook{}).Create(book).Error
}

// GetByID 根据ID查询地址
func (d *AddressBookDAO) GetByID(db *gorm.DB, id int) (*entity.AddressBook, error) {
	var address entity.AddressBook
	result := db.Model(&entity.AddressBook{}).Where("id = ?", id).First(&address)
	return &address, result.Error
}

// Update 更新地址信息
func (d *AddressBookDAO) Update(db *gorm.DB, address *entity.AddressBook) error {
	return db.Table("address_book").Where("id = ?", address.ID).Updates(address).Error
}

// SetNonDefault 将指定用户的所有地址设置为非默认
func (d *AddressBookDAO) SetNonDefault(db *gorm.DB, userID int) error {
	return db.Model(&entity.AddressBook{}).Where("user_id = ?", userID).Update("is_default", constant.NonDefaultAddress).Error
}

// SetDefault 将对应地址ID设置为默认地址
func (d *AddressBookDAO) SetDefault(db *gorm.DB, id int) error {
	return db.Model(&entity.AddressBook{}).Where("id = ?", id).Update("is_default", constant.DefaultAddress).Error
}

func (d *AddressBookDAO) DeleteByID(db *gorm.DB, id int) error {
	return db.Where("id = ?", id).Delete(&entity.AddressBook{}).Error
}
