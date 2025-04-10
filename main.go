package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"takeout/common/config"
	"takeout/common/global"
	"takeout/common/logger"
	"takeout/router"

	"go.uber.org/zap"
)

// @title takeout API
// @version 1.0
// @description take_out
// @BasePath /
func main() {
	// 初始化配置、日志和数据库
	if err := config.Init(""); err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		os.Exit(1)
	}
	defer config.Close()
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(global.Logger)

	// 初始化路由
	r := router.InitRouter()

	// 启动HTTP服务器
	serverAddr := fmt.Sprintf(":%d", global.Config.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		logger.Info("Server is running", zap.String("address", serverAddr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 创建上下文用于通知服务器结束
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
