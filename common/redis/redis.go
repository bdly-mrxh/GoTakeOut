package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"takeout/common/global"
	"time"
)

func InitRedis() error {
	global.Redis = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", global.Config.Redis.Host, global.Config.Redis.Port),
		Password:        global.Config.Redis.Password,
		DB:              global.Config.Redis.Database,
		MinIdleConns:    global.Config.Redis.MinIdleConns,
		ConnMaxIdleTime: time.Duration(global.Config.Redis.IdleTimeout) * time.Second,
		DialTimeout:     time.Duration(global.Config.Redis.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(global.Config.Redis.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(global.Config.Redis.WriteTimeout) * time.Second,
	})
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := global.Redis.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

func Close() error {
	if global.Redis != nil {
		return global.Redis.Close()
	}
	return nil
}
