package database

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"takeout/common/global"
	"takeout/model/entity"
)

// MigrateDB 数据库迁移
func MigrateDB() error {
	// 自动迁移表结构
	err := global.DB.AutoMigrate(
		&entity.Employee{},
		&entity.Category{},
		&entity.Dish{},
		&entity.DishFlavor{},
		&entity.Setmeal{},
		&entity.SetmealDish{},
		&entity.User{},
		&entity.ShoppingCart{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 初始化管理员账号
	return initAdminAccount()
}

// 初始化管理员账号
func initAdminAccount() error {
	// 检查是否已存在管理员账号
	var count int64
	global.DB.Model(&entity.Employee{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 创建默认管理员账号
	md5Password := md5.Sum([]byte("123456"))
	md5PasswordStr := hex.EncodeToString(md5Password[:])

	admin := entity.Employee{
		Username: "admin",
		Password: md5PasswordStr,
		Name:     "管理员",
		Phone:    "13800000000",
		Sex:      "1",
		IdNumber: "110101199001010001",
		Status:   1,
	}

	result := global.DB.Create(&admin)
	if result.Error != nil {
		return fmt.Errorf("创建管理员账号失败: %w", result.Error)
	}

	return nil
}
