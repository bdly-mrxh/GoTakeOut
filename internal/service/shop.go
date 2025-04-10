package service

import (
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/internal/dao"
)

// ShopService 店铺服务
type ShopService struct {
	shopDAO dao.ShopDAO
}

// SetStatus 设置店铺状态
func (s *ShopService) SetStatus(status int) error {
	err := s.shopDAO.SetStatus(status)
	if err != nil {
		return errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
	}
	return nil
}

// GetStatus 获取店铺状态
func (s *ShopService) GetStatus() (int, error) {
	status, err := s.shopDAO.GetStatus()
	if err != nil {
		return 0, errs.Wrap(err, constant.CodeInternalError, constant.MsgQueryFail)
	}
	return status, nil
}
