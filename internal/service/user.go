package service

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/internal/dao"
	"takeout/model/dto"
	"takeout/model/entity"
	"takeout/model/vo"
)

// UserService 用户服务
type UserService struct {
	userDAO dao.UserDAO
}

// Login 微信用户登录
func (s *UserService) Login(userLoginDTO *dto.UserLoginDTO) (*vo.UserLoginVO, error) {
	// 获取用户openid
	openid, err := s.userDAO.GetOpenId(userLoginDTO.Code)
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgUserLoginFail)
	}
	if openid == "" {
		return nil, errs.New(constant.CodeBusinessError, constant.MsgUserLoginFail)
	}
	// 查看用户是否注册
	user, err := s.userDAO.FindUserByOpenID(global.DB, openid)
	if err != nil {
		// 未注册自动注册
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &entity.User{OpenID: openid}
			err = s.userDAO.Create(global.DB, user)
			if err != nil {
				return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
			}
		} else {
			return nil, errs.Wrap(err, constant.CodeDatabaseError, constant.MsgDatabaseError)
		}
	}
	// 生成token
	token, err := utils.GenerateToken(constant.UserID, strconv.Itoa(user.ID))
	if err != nil {
		return nil, errs.Wrap(err, constant.CodeInternalError, constant.MsgServerError)
	}
	// 返回参数
	return &vo.UserLoginVO{
		ID:     user.ID,
		OpenID: openid,
		Token:  token,
	}, nil
}
