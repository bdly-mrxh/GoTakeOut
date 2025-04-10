package dao

import (
	"encoding/json"
	"gorm.io/gorm"
	"takeout/common/constant"
	"takeout/common/global"
	"takeout/common/utils"
	"takeout/model/entity"
	"time"
)

type UserDAO struct{}

// GetOpenId 获取openid
func (dao *UserDAO) GetOpenId(code string) (string, error) {
	// query参数
	values := map[string]string{
		"appid":      global.Config.Wechat.AppID,
		"secret":     global.Config.Wechat.AppSecretKey,
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	res, err := utils.DoGET(constant.WeChatLoginUrl, values) // body是一个json字符串
	if err != nil {
		return "", err
	}
	wxData := struct {
		OpenID string `json:"openid"`
	}{}
	if err = json.Unmarshal([]byte(res), &wxData); err != nil {
		return "", err
	}
	return wxData.OpenID, nil
}

// FindUserByOpenID 查询用户是否存在
func (dao *UserDAO) FindUserByOpenID(db *gorm.DB, openid string) (*entity.User, error) {
	var user entity.User
	result := db.Model(&entity.User{}).Where("openid = ?", openid).First(&user)
	return &user, result.Error
}

// Create 注册用户
func (dao *UserDAO) Create(db *gorm.DB, user *entity.User) error {
	return db.Model(&entity.User{}).Create(user).Error
}

// GetByID 根据ID获取用户信息
func (dao *UserDAO) GetByID(db *gorm.DB, id int) (*entity.User, error) {
	var user entity.User
	result := db.Model(&entity.User{}).Where("id = ?", id).First(&user)
	return &user, result.Error
}

// GetCount 获取用户人数
func (dao *UserDAO) GetCount(db *gorm.DB, begin *time.Time, end *time.Time) (int64, error) {
	var cnt int64
	query := db.Model(&entity.User{})
	if begin != nil {
		query = query.Where("create_time >= ?", begin)
	}
	if end != nil {
		query = query.Where("create_time <= ?", end)
	}
	err := query.Count(&cnt).Error
	return cnt, err
}
