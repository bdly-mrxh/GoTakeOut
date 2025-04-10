package entity

import (
	"github.com/shopspring/decimal"
	"takeout/model/wrap"
)

// 可以使用 *int 插入零值

// Order 订单数据模型
type Order struct {
	ID                    int             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Number                string          `json:"number"`
	Status                int             `json:"status"`
	UserID                int             `json:"userId" gorm:"column:user_id"`
	AddressBookID         int             `json:"addressBookId" gorm:"column:address_book_id"`
	OrderTime             wrap.LocalTime  `json:"orderTime" gorm:"column:order_time;autoCreateTime"`
	CheckoutTime          wrap.LocalTime  `json:"checkoutTime" gorm:"column:checkout_time"`
	CancelTime            wrap.LocalTime  `json:"cancelTime" gorm:"column:cancel_time"`
	EstimatedDeliveryTime wrap.LocalTime  `json:"estimatedDeliveryTime" gorm:"column:estimated_delivery_time"`
	DeliveryTime          wrap.LocalTime  `json:"deliveryTime" gorm:"column:delivery_time"`
	PayMethod             int             `json:"payMethod" gorm:"column:pay_method"`
	PayStatus             int             `json:"payStatus" gorm:"column:pay_status;default:0"`
	Amount                decimal.Decimal `json:"amount"`
	Remark                string          `json:"remark"`
	Username              string          `json:"username" gorm:"column:user_name"`
	Phone                 string          `json:"phone"`
	Address               string          `json:"address"`
	Consignee             string          `json:"consignee"`
	CancelReason          string          `json:"cancelReason" gorm:"column:cancel_reason"`
	RejectionReason       string          `json:"rejectionReason" gorm:"column:rejection_reason"`
	DeliveryStatus        int             `json:"deliveryStatus" gorm:"column:delivery_status"`
	PackAmount            decimal.Decimal `json:"packAmount" gorm:"column:pack_amount"`
	TablewareNumber       int             `json:"tablewareNumber" gorm:"column:tableware_number"`
	TablewareStatus       int             `json:"tablewareStatus" gorm:"column:tableware_status"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}
