package dao

import (
	"context"
	"strconv"
	"takeout/common/constant"
	"takeout/common/global"
)

// ShopDAO 店铺数据访问对象
type ShopDAO struct{}

// SetStatus 设置店铺状态
func (dao *ShopDAO) SetStatus(status int) error {
	ctx := context.Background()
	return global.Redis.Set(ctx, constant.RedisKeyShopStatus, strconv.Itoa(status), 0).Err()
}

// GetStatus 获取店铺状态
func (dao *ShopDAO) GetStatus() (int, error) {
	ctx := context.Background()
	statusStr, err := global.Redis.Get(ctx, constant.RedisKeyShopStatus).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(statusStr)
}
