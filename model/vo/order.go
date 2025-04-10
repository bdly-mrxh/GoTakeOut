package vo

import (
	"github.com/shopspring/decimal"
	"takeout/model/entity"
	"takeout/model/wrap"
)

type OrderSubmitVO struct {
	ID          int             `json:"id"`
	OrderNumber string          `json:"orderNumber"`
	OrderAmount decimal.Decimal `json:"orderAmount"`
	OrderTime   wrap.LocalTime  `json:"orderTime"`
}

type OrderPaymentVO struct {
	NonceStr   string `json:"nonceStr"`   // 随机字符串
	PaySign    string `json:"paySign"`    // 签名
	TimeStamp  string `json:"timeStamp"`  // 时间戳
	SignType   string `json:"signType"`   // 签名算法
	PackageStr string `json:"packageStr"` // 统一下单接口返回的 prepay_id 参数值
}

// OrderVO 查询订单详情返回数据模型
type OrderVO struct {
	entity.Order
	OrderDishes     string                `json:"orderDishes"`
	OrderDetailList []*entity.OrderDetail `json:"orderDetailList"`
}

// OrderStatisticsVO 订单数量统计返回数据模型
type OrderStatisticsVO struct {
	ToBeConfirmed      int64 `json:"toBeConfirmed"`
	Confirmed          int64 `json:"confirmed"`
	DeliveryInProgress int64 `json:"deliveryInProgress"`
}
