package vo

type BusinessDataVO struct {
	Turnover            float64 `json:"turnover"`            // 营业额
	ValidOrderCount     int64   `json:"validOrderCount"`     // 有效订单数量
	OrderCompletionRate float64 `json:"orderCompletionRate"` // 订单完成率
	UnitPrice           float64 `json:"unitPrice"`           // 平均客单价
	NewUsers            int64   `json:"newUsers"`            // 新增用户数
}

type OrderOverViewVO struct {
	WaitingOrders   int64 `json:"waitingOrders"`
	DeliveredOrders int64 `json:"deliveredOrders"`
	CompletedOrders int64 `json:"completedOrders"`
	CancelledOrders int64 `json:"cancelledOrders"`
	AllOrders       int64 `json:"allOrders"`
}

type DishOverViewVO struct {
	Sold         int64 `json:"sold"`         // 起售数量
	Discontinued int64 `json:"discontinued"` // 停售数量
}

type SetmealOverViewVO struct {
	Sold         int64 `json:"sold"`
	Discontinued int64 `json:"discontinued"`
}
