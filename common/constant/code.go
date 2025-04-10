package constant

// 错误码常量
const (
	// CodeSuccess 成功码
	CodeSuccess       = 1
	CodeBusinessError = -1 // 业务错误

	CodeJWTParseError = 2

	CodeBadRequest    = 400 // 无效请求
	CodeUnauthorized  = 401 // 未授权
	CodeForbidden     = 403 // 禁止访问
	CodeNotFound      = 404 // 资源不存在
	CodeServerError   = 500 // 服务器错误
	CodeUserNotExist  = 100 // 用户不存在
	CodePasswordError = 101 // 密码错误
	CodeUserDisabled  = 102 // 用户被禁用

	CodeDatabaseError = 900 // 数据库错误
	CodeCacheError    = 901 // 缓存错误
	CodeConfigError   = 902 // 配置错误
)

// 员工相关错误码 (1000 - 1099)
const (
	CodeEmployLoginFail       = 1000 // 员工登录失败
	CodeEmployeeCreateFail    = 1001 // 员工创建失败
	CodeEmployeeUpdateFail    = 1002 // 员工更新失败
	CodeEmployeePageQueryFail = 1003 // 员工分页查询失败
)

// 分类相关错误 (1100 - 1199)
const (
	CodeCategoryCreateFail = 1100 // 分类创建失败
	CodeDeleteCategoryFail = 1101 // 分类删除失败
)

const (
	CodeInternalError = 50 // 系统内部错误
)
