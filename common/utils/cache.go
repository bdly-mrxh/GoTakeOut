package utils

import (
	"context"
	"takeout/common/global"
)

// CleanCache 清理缓存
func CleanCache(pattern string) error {
	ctx := context.Background()
	keys, err := global.Redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	_, err = global.Redis.Del(ctx, keys...).Result()
	if err != nil {
		return err
	}
	return nil
}
