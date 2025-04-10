package constant

// 系统常量定义
const (
	// DefaultTimeFormat 时间格式
	DefaultTimeFormat = "2006-01-02 15:04:05"

	// DefaultPassword 默认密码
	DefaultPassword = "123456"

	ID     = "ID"
	EmpID  = "empId"
	UserID = "userId"

	RedisKeyShopStatus = "shop::status"

	DefaultPageSize = 10 // 默认分页大小
	DefaultPageNum  = 1  // 默认页码

	WeChatLoginUrl = "https://api.weixin.qq.com/sns/jscode2session"

	CacheDishKey    = "dish::"
	CacheSetmealKey = "setmeal::"
)

// 员工状态常量
const (
	EmployeeStatusEnable  = 1 // 员工启用
	EmployeeStatusDisable = 0 // 员工禁用
)

// Operator 自动填充的标签
const (
	Create     = "create"
	Update     = "update"
	CreateUser = "CreateUser"
	UpdateUser = "UpdateUser"
)

const (
	DefaultStatus = 0
	InvalidStatus = -1

	DishEnable     = 1 // 菜品起售
	DishDisable    = 0
	SetmealEnable  = 1
	SetmealDisable = 0
)

const (
	DefaultAddress    = 1  // 默认地址
	NonDefaultAddress = 0  // 非默认地址
	NotSetAddress     = -1 // 未设置地址
)

// 订单相关常量
const (
	PendingPayment     = 1 // 待付款
	ToBeConfirmed      = 2 // 待接单
	Confirmed          = 3 // 已接单
	DeliveryInProgress = 4 // 派送中
	Completed          = 5 // 已完成
	Cancelled          = 6 // 已取消

	UnPaid = 0 // 未支付
	Paid   = 1 // 已支付
	ReFund = 2 // 退款

	NotifyOrder = 1 // 通知接单
	UserRemind  = 2 // 用户催单
)
