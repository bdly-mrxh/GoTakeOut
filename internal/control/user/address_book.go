package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"takeout/common/constant"
	"takeout/common/logger"
	"takeout/common/response"
	"takeout/internal/service"
	"takeout/model/dto"
)

type AddressBookController struct {
	addressBookService service.AddressBookService
}

func NewAddressBookController() *AddressBookController {
	return &AddressBookController{}
}

// List 查询当前用户的所有地址信息
func (c *AddressBookController) List(ctx *gin.Context) {
	list, err := c.addressBookService.List(ctx)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, list)
}

// Add 新增地址
func (c *AddressBookController) Add(ctx *gin.Context) {
	var addDTO dto.AddressBookDTO
	if err := ctx.ShouldBindJSON(&addDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.addressBookService.Add(ctx, &addDTO)
	if err != nil {
		logger.Error(constant.MsgCreateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgCreateSuccess, nil)
}

// GetByID 根据ID查询地址详细信息
func (c *AddressBookController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	address, err := c.addressBookService.GetByID(id)
	if err != nil {
		logger.Error(constant.MsgQueryFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, address)
}

// Update 根据ID修改地址
func (c *AddressBookController) Update(ctx *gin.Context) {
	var addressDTO dto.AddressBookDTO
	if err := ctx.ShouldBindJSON(&addressDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.addressBookService.Update(&addressDTO)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}

// SetDefault 设置默认地址
func (c *AddressBookController) SetDefault(ctx *gin.Context) {
	var addressDTO dto.AddressBookDTO
	if err := ctx.ShouldBindJSON(&addressDTO); err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err := c.addressBookService.SetDefault(ctx, &addressDTO)
	if err != nil {
		logger.Error(constant.MsgUpdateFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgUpdateSuccess, nil)
}

// DeleteByID 根据ID删除地址
func (c *AddressBookController) DeleteByID(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		logger.Error(constant.MsgMissingRequest)
		response.BadRequest(ctx, constant.MsgMissingRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(constant.MsgBadRequest, zap.Error(err))
		response.BadRequest(ctx, constant.MsgBadRequest)
		return
	}

	err = c.addressBookService.DeleteByID(id)
	if err != nil {
		logger.Error(constant.MsgDeleteFail, zap.Error(err))
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgDeleteSuccess, nil)
}

// GetDefault 查询默认地址
func (c *AddressBookController) GetDefault(ctx *gin.Context) {
	address, err := c.addressBookService.GetDefault(ctx)
	if err != nil {
		if strings.Contains(err.Error(), constant.MsgNotExistDefaultAddress) {
			logger.Info(constant.MsgQueryFail, zap.Error(err))
		} else {
			logger.Error(constant.MsgQueryFail, zap.Error(err))
		}
		response.ErrorResponse(ctx, err)
		return
	}
	response.Success(ctx, constant.MsgQuerySuccess, address)
}
