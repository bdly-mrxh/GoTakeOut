package entity

import "takeout/model/wrap"

type User struct {
	ID         int            `json:"id" gorm:"column:id;primary_key"`
	OpenID     string         `json:"openid" gorm:"column:openid"`
	Name       string         `json:"name"`
	Phone      string         `json:"phone"`
	Sex        string         `json:"sex"`
	IdNumber   string         `json:"idNumber" gorm:"column:id_number"`
	Avatar     string         `json:"avatar" `
	CreateTime wrap.LocalTime `json:"createTime" gorm:"column:create_time;autoCreateTime"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
