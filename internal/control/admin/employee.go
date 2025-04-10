package admin

import (
	"strconv"
	"takeout/internal/service"

	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/model/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EmployeeController 员工控制器
type EmployeeController struct {
	employeeService service.EmployeeService
}

// NewEmployeeController 创建员工控制器
func NewEmployeeController() *EmployeeController {
	return &EmployeeController{}
}

// Logout 员工登出
// @Summary 员工登出
// @Description 员工登出系统
// @Tags 员工管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "登出成功"
// @Router /admin/employee/logout [post]
func (c *EmployeeController) Logout(ctx *gin.Context) {
	// 登出成功
	logger.Info(constant.MsgEmployeeLogoutSuccess)
	response.Success(ctx, constant.MsgEmployeeLogoutSuccess, nil)
}

// Login 员工登录
// @Summary 员工登录
// @Description 员工登录系统
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param login body dto.EmployeeLoginDTO true "登录信息"
// @Success 200 {object} response.Response{data=vo.EmployeeLoginVO} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Router /admin/employee/login [post]
func (c *EmployeeController) Login(ctx *gin.Context) {
	// 绑定请求参数
	var loginDTO dto.EmployeeLoginDTO
	if err := ctx.ShouldBindJSON(&loginDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层处理登录
	loginVO, err := c.employeeService.Login(&loginDTO)
	if err != nil {
		logger.Error(constant.MsgEmployeeLoginFail, zap.Error(err), zap.String("username", loginDTO.Username))
		response.ErrorResponse(ctx, err)
		return
	}

	// 登录成功
	logger.Info(constant.MsgEmployeeLoginSuccess, zap.String("username", loginDTO.Username), zap.String("token", loginVO.Token))
	response.Success(ctx, constant.MsgEmployeeLoginSuccess, loginVO)
}

// Create 创建新员工
func (c *EmployeeController) Create(ctx *gin.Context) {
	// 绑定请求参数
	var createDTO dto.EmployeeCreateDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层创建员工
	err := c.employeeService.Create(ctx, &createDTO)
	if err != nil {
		logger.Error(constant.MsgEmployeeCreateFail, zap.Error(err), zap.String("username", createDTO.Username))
		response.ErrorResponse(ctx, err)
		return
	}

	// 创建成功
	logger.Info(constant.MsgEmployeeCreateSuccess, zap.String("username", createDTO.Username))
	response.Success(ctx, constant.MsgEmployeeCreateSuccess, nil)
}

// GetById 根据ID查询员工
func (c *EmployeeController) GetById(ctx *gin.Context) {
	// 获取路径参数中的ID
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err), zap.String("id", idStr))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层查询员工
	employeeVO, err := c.employeeService.GetById(id)
	if err != nil {
		logger.Error(constant.MsgQueryEmployeeFail, zap.Error(err), zap.Int("id", id))
		response.ErrorResponse(ctx, err)
		return
	}

	// 查询成功
	logger.Info(constant.MsgQueryEmployeeSuccess, zap.Int("id", id))
	response.Success(ctx, constant.MsgQueryEmployeeSuccess, employeeVO)
}

// UpdateStatus 更新员工状态
func (c *EmployeeController) UpdateStatus(ctx *gin.Context) {
	// 获取路径参数中的状态值
	statusStr := ctx.Param("status")
	if statusStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err), zap.String("status", statusStr))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 验证状态值是否有效
	if status != constant.EmployeeStatusEnable && status != constant.EmployeeStatusDisable {
		logger.Error(constant.MsgBadRequest, zap.Int("status", status))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 获取查询参数中的员工ID
	idStr := ctx.Query("id")
	if idStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err), zap.String("id", idStr))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层更新员工状态
	err = c.employeeService.UpdateStatusById(status, id)
	if err != nil {
		logger.Error(constant.MsgEmployeeStatusUpdateFail, zap.Error(err), zap.Int("status", status), zap.Int("id", id))
		response.ErrorResponse(ctx, err)
		return
	}

	// 更新成功
	statusText := "启用"
	if status == constant.EmployeeStatusDisable {
		statusText = "禁用"
	}
	logger.Info(constant.MsgEmployeeStatusUpdateSuccess, zap.Int("status", status))
	response.Success(ctx, "员工"+statusText+"成功", nil)
}

// UpdatePassword 修改员工密码
func (c *EmployeeController) UpdatePassword(ctx *gin.Context) {
	// 绑定请求参数
	var passwordDTO dto.EmployeePasswordDTO
	if err := ctx.ShouldBindJSON(&passwordDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层修改密码
	err := c.employeeService.UpdatePassword(ctx, &passwordDTO)
	if err != nil {
		logger.Error(constant.MsgEmployeeChangePasswordFail, zap.Error(err), zap.Int("id", passwordDTO.EmpId))
		response.ErrorResponse(ctx, err)
		return
	}

	// 修改成功
	logger.Info(constant.MsgEmployeeChangePasswordSuccess, zap.Int("id", passwordDTO.EmpId))
	response.Success(ctx, constant.MsgEmployeeChangePasswordSuccess, nil)
}

// Update 更新员工信息
func (c *EmployeeController) Update(ctx *gin.Context) {
	// 绑定请求参数
	var updateDTO dto.EmployeeUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 调用服务层更新员工信息
	err := c.employeeService.Update(ctx, &updateDTO)
	if err != nil {
		logger.Error(constant.MsgEmployeeUpdateFail, zap.Error(err), zap.Int("id", updateDTO.ID))
		response.ErrorResponse(ctx, err)
		return
	}

	// 更新成功
	logger.Info(constant.MsgEmployeeUpdateSuccess, zap.Int("id", updateDTO.ID))
	response.Success(ctx, constant.MsgEmployeeUpdateSuccess, nil)
}

// Page 分页查询员工信息
func (c *EmployeeController) Page(ctx *gin.Context) {
	// 绑定查询参数
	var pageDTO dto.EmployeePageDTO
	if err := ctx.ShouldBindQuery(&pageDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	// 设置默认值
	if pageDTO.Page <= 0 {
		pageDTO.Page = constant.DefaultPageNum
	}
	if pageDTO.PageSize <= 0 {
		pageDTO.PageSize = constant.DefaultPageSize
	}

	// 调用服务层分页查询
	pageResult, err := c.employeeService.PageQuery(&pageDTO)
	if err != nil {
		logger.Error(constant.MsgPageQueryEmployeeFail, zap.Error(err), zap.Any("pageDTO", pageDTO))
		response.ErrorResponse(ctx, err)
		return
	}

	// 查询成功
	logger.Info(constant.MsgPageQueryEmployeeSuccess, zap.Int("page", pageDTO.Page), zap.Int("pageSize", pageDTO.PageSize))
	response.Success(ctx, constant.MsgPageQueryEmployeeSuccess, pageResult)
}
