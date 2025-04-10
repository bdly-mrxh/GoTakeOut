package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
	"takeout/common/constant"
	"takeout/common/errs"
)

// AutoFill 自动填充装饰器，实现公共字段自动填充，可进行事务操作
func AutoFill(next func(*gin.Context, *gorm.DB, any, string) error) func(*gin.Context, *gorm.DB, any, string) error {
	return func(ctx *gin.Context, db *gorm.DB, entity any, op string) error {
		// 传入的必须是 entity 指针类型
		val := reflect.ValueOf(entity) // 返回一个新的值
		if val.Kind() != reflect.Ptr {
			return errs.New(constant.CodeInternalError, "entity 必须是指针")
		}
		elem := val.Elem()
		// 获取id
		id, err := GetId(ctx)
		if err != nil {
			return err
		}

		// 为字段赋值
		switch op {
		case constant.Create:
			for i := 0; i < elem.NumField(); i++ {
				field := elem.Type().Field(i)
				fieldValue := elem.Field(i)
				if fieldValue.CanSet() {
					if field.Name == constant.CreateUser || field.Name == constant.UpdateUser {
						fieldValue.SetInt(int64(id))
						// fieldValue.Set(reflect.ValueOf(id)) 性能差
					}
				}
			}
		case constant.Update:
			for i := 0; i < elem.NumField(); i++ {
				field := elem.Type().Field(i)
				fieldValue := elem.Field(i)
				if fieldValue.CanSet() {
					if field.Name == constant.UpdateUser {
						fieldValue.SetInt(int64(id))
						//fieldValue.Set(reflect.ValueOf(id))
					}
				}
			}
		}

		return next(ctx, db, entity, op)
	}
}
