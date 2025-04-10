package dao

import (
	"gorm.io/gorm"
	"takeout/model/entity"
)

type OrderDetailDAO struct{}

// BatchInsert 批量插入订单详情
func (d *OrderDetailDAO) BatchInsert(db *gorm.DB, list []*entity.OrderDetail) error {
	return db.Create(&list).Error
}

// GetByOrderID 根据订单ID查询订单详情
func (d *OrderDetailDAO) GetByOrderID(db *gorm.DB, orderID int) ([]*entity.OrderDetail, error) {
	var list []*entity.OrderDetail
	result := db.Model(&entity.OrderDetail{}).Where("order_id = ?", orderID).Find(&list)
	return list, result.Error
}
