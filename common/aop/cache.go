package aop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"takeout/common/constant"
	"takeout/common/errs"
	"takeout/common/global"
	"time"
)

// CacheOptions 定义缓存选项
type CacheOptions struct {
	CacheName  string
	Key        any
	Expiration time.Duration
	AllEntries bool
	//KeyGenerator func(method string, args ...any) string
	//Condition    func(args ...any) bool
}

// Cacheable 实现类似Spring的@Cacheable注解
// 用于包装函数，提供缓存功能
// 函数必须返回结果和error
func Cacheable[T any, F func() (T, error)](fn F, opts *CacheOptions) F {
	return func() (T, error) {
		var result T
		ctx := context.Background()
		// 拼接 key
		cacheKey := opts.CacheName + fmt.Sprint(opts.Key)
		// 尝试从缓存获取
		cached, err := global.Redis.Get(ctx, cacheKey).Result()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				return result, errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
			}
		} else {
			err = json.Unmarshal([]byte(cached), &result)
			if err != nil {
				return result, errs.Wrap(err, constant.CodeInternalError, constant.MsgUnmarshalFail)
			}
			return result, nil
		}

		// 缓存未命中，调用原始函数
		result, err = fn()
		if err != nil {
			return result, err
		}
		// 序列化数据
		resJson, err := json.Marshal(result)
		if err != nil {
			return result, errs.Wrap(err, constant.CodeInternalError, constant.MsgMarshalFail)
		}

		// 将结果存入缓存
		err = global.Redis.Set(ctx, cacheKey, string(resJson), opts.Expiration).Err()
		if err != nil {
			return result, errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
		}

		return result, nil
	}
}

// CacheEvict 实现类似Spring的@CacheEvict注解
// 在函数执行后删除缓存
func CacheEvict(fn func() error, opts *CacheOptions) func() error {
	return func() error {
		var (
			cacheKeys []string
			err       error
		)

		err = fn()
		if err != nil {
			return err
		}
		// 删除缓存
		ctx := context.Background()
		if opts.AllEntries {
			cacheKeys, err = global.Redis.Keys(ctx, opts.CacheName+"*").Result()
			if err != nil {
				return errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
			}
		} else {
			cacheKeys = append(cacheKeys, opts.CacheName+fmt.Sprint(opts.Key))
		}
		if len(cacheKeys) != 0 {
			err = global.Redis.Del(ctx, cacheKeys...).Err()
			if err != nil {
				return errs.Wrap(err, constant.CodeCacheError, constant.MsgCacheError)
			}
		}
		return nil
	}
}
