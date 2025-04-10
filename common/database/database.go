package database

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"

	"takeout/common/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormZapLogger 完全自定义的GORM日志实现
type GormZapLogger struct {
	Logger                    *zap.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewGormZapLogger 创建新日志实例
func NewGormZapLogger(zapLogger *zap.Logger) *GormZapLogger {
	return &GormZapLogger{
		Logger:                    zapLogger,
		LogLevel:                  logger.Info,
		SlowThreshold:             500 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// 实现logger.Interface接口的所有方法

// LogMode 方法
func (l *GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 方法
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Logger.Sugar().Infof(msg, data...)
	}
}

// Warn 方法
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Logger.Sugar().Warnf(msg, data...)
	}
}

// Error 方法
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Logger.Sugar().Errorf(msg, data...)
	}
}

// Trace 方法
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 创建字段
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	// 根据不同情况记录日志
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || err != gorm.ErrRecordNotFound):
		fields = append(fields, zap.Error(err))
		l.Logger.Error("GORM错误", fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.Logger.Warn("GORM慢查询", fields...)
	case l.LogLevel >= logger.Info:
		l.Logger.Info("GORM查询", fields...)
	}
}

// InitDB 初始化数据库连接
func InitDB() error {
	dbConfig := global.Config.Database

	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset)

	// fmt.Println(dsn)

	// 连接数据库
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewGormZapLogger(global.Logger),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 设置连接池
	sqlDB, err := global.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// 执行数据库迁移
	//if err = MigrateDB(); err != nil {
	//	return fmt.Errorf("database migration failed: %w", err)
	//}

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if global.DB != nil {
		sqlDB, err := global.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}
