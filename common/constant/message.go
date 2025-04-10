package constant

// 响应消息常量
const (
	MsgSuccess        = "操作成功"
	MsgServerError    = "服务器内部错误"
	MsgCacheError     = "缓存错误"
	MsgDatabaseError  = "数据库错误"
	MsgBadRequest     = "无效的请求参数"
	MsgMissingRequest = "缺少请求参数"
	MsgUnauthorized   = "未授权访问"
	MsgNotFound       = "资源不存在"
	MsgNameConflict   = "名称冲突"
	MsgBusinessError  = "业务错误"

	MsgJWTUnKnownSigningMethod = "JWT未知的签名方法"
	MsgJWTParseFail            = "JWT解析失败"
	MsgJWTWithoutToken         = "JWT未携带token"

	MsgGetAccountInfoFail = "未能获取当前账户信息"
	MsgGetIDFail          = "获取ID失败"

	MsgCopyPropertiesFail = "拷贝属性失败"
	MsgTypeConversionFail = "类型转换失败"
	MsgUnmarshalFail      = "反序列化失败"
	MsgMarshalFail        = "序列化失败"

	MsgKeyDuplicateError = "Duplicate entry"

	MsgQuerySuccess = "查询成功"
	MsgQueryFail    = "查询失败"

	MsgDeleteFail    = "删除失败"
	MsgDeleteSuccess = "删除成功"
	MsgCreateFail    = "创建失败"
	MsgCreateSuccess = "创建成功"
	MsgUpdateFail    = "更新失败"
	MsgUpdateSuccess = "更新成功"

	MsgDatabaseTransactionFail = "数据库事务操作失败"
	MsgPasswordIncorrect       = "密码不正确"
	MsgUserNotExist            = "用户不存在"
)

// 员工相关消息
const (
	MsgEmployeeCreateSuccess         = "员工创建成功"
	MsgEmployeeCreateFail            = "员工创建失败"
	MsgEmployeeLoginSuccess          = "员工登录成功"
	MsgEmployeeLoginFail             = "员工登录失败"
	MsgEmployeeLogoutSuccess         = "员工登出成功"
	MsgQueryEmployeeFail             = "查询员工失败"
	MsgQueryEmployeeSuccess          = "查询员工成功"
	MsgEmployeeStatusUpdateFail      = "员工状态更新失败"
	MsgEmployeeStatusUpdateSuccess   = "员工状态更新成功"
	MsgEmployeeChangePasswordFail    = "员工修改密码失败"
	MsgEmployeeChangePasswordSuccess = "员工修改密码成功"
	MsgEmployeeUpdateFail            = "员工信息更新失败"
	MsgEmployeeUpdateSuccess         = "员工信息更新成功"
	MsgPageQueryEmployeeFail         = "员工分页查询失败"
	MsgPageQueryEmployeeSuccess      = "员工分页查询成功"
)

// 分类相关消息
const (
	MsgCategoryCreateFail            = "分类创建失败"
	MsgCategoryCreateSuccess         = "分类创建成功"
	MsgCategoryStatusFail            = "更新分类状态失败"
	MsgCategoryStatusSuccess         = "更新分类状态成功"
	MsgCategoryUpdateFail            = "分类更新失败"
	MsgCategoryUpdateSuccess         = "分类更新成功"
	MsgExistAssociativeDishOrSetmeal = "存在关联菜品或套餐"
)

// 菜品相关消息
const (
	MsgDishOnSale                 = "菜品正在销售中，不能删除"
	MsgDishAssociativeWithSetmeal = "菜品与套餐关联，不能删除"
)

// 套餐相关消息
const (
	MsgSetmealOnSale                   = "套餐正在销售中，不能删除"
	MsgSetmealAssociativeDishHalfSales = "套餐关联的菜品才停售中，不能起售套餐"
)

// 员工相关消息
const (
	MsgUserLogoutSuccess = "用户退出成功"
	MsgUserLoginSuccess  = "用户登录成功"
	MsgUserLoginFail     = "用户登录失败"
)

// 购物车相关消息
const (
	MsgShoppingCartAddFail    = "购物车添加失败"
	MsgShoppingCartAddSuccess = "购物车添加成功"
)

const MsgNotExistDefaultAddress = "未设置默认地址"

// 订单相关消息
const (
	MsgAddressBookIsNull  = "用户地址为空"
	MsgShoppingCartIsNull = "购物车为空"
	MsgOrderSubmitFail    = "订单提交失败"
	MsgOrderSubmitSuccess = "订单提交成功"
	MsgOrderPaid          = "订单已支付"
	MsgPayFail            = "支付失败"
	MsgOrderNotFound      = "未查询到订单"
	MsgOrderStatusError   = "订单状态错误"
	MsgOrderCancelFail    = "订单取消失败"
	MsgOrderCancelSuccess = "订单取消成功"
)
