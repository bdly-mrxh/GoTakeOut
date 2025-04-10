package dto

import (
	"github.com/shopspring/decimal"
	"takeout/model/entity"
	"takeout/model/wrap"
)

// OrderSubmitDTO 用户下单接口参数
type OrderSubmitDTO struct {
	AddressBookID         int             `json:"addressBookId"`
	Amount                decimal.Decimal `json:"amount"`
	DeliveryStatus        int             `json:"deliveryStatus"`
	EstimatedDeliveryTime wrap.LocalTime  `json:"estimatedDeliveryTime"`
	PackAmount            decimal.Decimal `json:"packAmount"`
	PayMethod             int             `json:"payMethod"`
	Remark                string          `json:"remark"`
	TablewareNumber       int             `json:"tablewareNumber"`
	TablewareStatus       int             `json:"tablewareStatus"`
}

type OrderDTO struct {
	ID            int                   `json:"id"`
	Number        string                `json:"number"`
	Status        int                   `json:"status"`
	UserID        int                   `json:"userId"`
	AddressBookID int                   `json:"addressBookId"`
	OrderTime     wrap.LocalTime        `json:"orderTime"`
	CheckoutTime  wrap.LocalTime        `json:"checkoutTime"`
	PayMethod     int                   `json:"payMethod"`
	Amount        decimal.Decimal       `json:"amount"`
	Remark        string                `json:"remark"`
	UserName      string                `json:"userName"`
	Phone         string                `json:"phone"`
	Address       string                `json:"address"`
	Consignee     string                `json:"consignee"`
	OrderDetails  []*entity.OrderDetail `json:"orderDetails"`
}

// OrderPaymentDTO 订单支付DTO
type OrderPaymentDTO struct {
	OrderNumber string `json:"orderNumber"`
	PayMethod   int    `json:"payMethod"`
}

// OrderPageQueryDTO 订单分页查询数据模型，时间是string
type OrderPageQueryDTO struct {
	Page      int    `form:"page" binding:"required"`
	PageSize  int    `form:"pageSize" binding:"required"`
	UserID    int    `form:"userId"`
	Number    string `form:"number"`
	Phone     string `form:"phone"`
	Status    int    `form:"status"`
	BeginTime string `form:"beginTime"`
	EndTime   string `form:"endTime"` // query中的时间不太好绑定
}

// OrderPageQueryNonDTO 订单分页查询数据模型
type OrderPageQueryNonDTO struct {
	Page      int            `form:"page" binding:"required"`
	PageSize  int            `form:"pageSize" binding:"required"`
	UserID    int            `form:"userId"`
	Number    string         `form:"number"`
	Phone     string         `form:"phone"`
	Status    int            `form:"status"`
	BeginTime wrap.LocalTime `form:"beginTime"`
	EndTime   wrap.LocalTime `form:"endTime"` // query中的时间不太好绑定
}

// OrderConfirmDTO 接单接收数据模型
type OrderConfirmDTO struct {
	OrderID int `json:"id"`
	Status  int `json:"status"`
}

// OrderRejectionDTO 拒单接收数据模型
type OrderRejectionDTO struct {
	OrderID         int    `json:"id"`
	RejectionReason string `json:"rejectionReason"`
}

// OrderCancelDTO 商家取消订单接收数据模型
type OrderCancelDTO struct {
	OrderID      int    `json:"id"`
	CancelReason string `json:"cancelReason"`
}
