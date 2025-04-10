package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"takeout/common/constant"
	"takeout/common/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// --- 配置常量 ---
const (
	// 调用者信息列的固定宽度 (例如: "gorm@v1.25.12/callbacks.go:134")
	callerWidth = 40 // 请根据你项目中路径的最大预期长度调整
)

// 2. 自定义带固定宽度和左对齐的Caller编码器
func paddedCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// caller.TrimmedPath() 返回格式如: "gorm@v1.25.12/callbacks.go:134"
	callerStr := caller.TrimmedPath()
	// 左对齐填充到指定宽度
	paddedCaller := fmt.Sprintf("%-*s", callerWidth, callerStr)
	enc.AppendString(paddedCaller)
}

// InitLogger 初始化日志
func InitLogger() error {
	// 创建日志目录
	logDir := filepath.Dir(global.Config.Log.Filename)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	// 设置日志级别
	var level zapcore.Level
	switch global.Config.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 控制台输出的编码器配置
	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeCaller = paddedCallerEncoder

	// 采用lumberjack来切割日志
	lumberJackLogger := &lumberjack.Logger{
		Filename:   global.Config.Log.Filename,
		MaxSize:    global.Config.Log.MaxSize,
		MaxAge:     global.Config.Log.MaxAge,
		MaxBackups: global.Config.Log.MaxBackups,
		Compress:   global.Config.Log.Compress,
	}

	// 同时输出到控制台和文件
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(lumberJackLogger), zapcore.WarnLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoderConfig), zapcore.AddSync(os.Stdout), level),
	)

	// 创建Logger
	global.Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	global.SugarLogger = global.Logger.Sugar()

	return nil
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	global.Logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	global.Logger.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	global.Logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	global.Logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	global.Logger.Fatal(msg, fields...)
}

// CustomTimeEncoder 自定义时间编码器，格式为 YYYY-MM-DD HH:MM:SS
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(constant.DefaultTimeFormat))
}

// Debugf 格式化调试日志
func Debugf(format string, args ...any) {
	global.SugarLogger.Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...any) {
	global.SugarLogger.Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...any) {
	global.SugarLogger.Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...any) {
	global.SugarLogger.Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...any) {
	global.SugarLogger.Fatalf(format, args...)
}
