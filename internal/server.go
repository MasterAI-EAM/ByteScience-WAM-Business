package internal

import (
	"ByteScience-WAM-Business/cron"
	"ByteScience-WAM-Business/internal/routers"
	"ByteScience-WAM-Business/middleware"
	"ByteScience-WAM-Business/pkg/db"
	"ByteScience-WAM-Business/pkg/gpt"
	"ByteScience-WAM-Business/pkg/logger"
	"ByteScience-WAM-Business/pkg/redis"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ByteScience-WAM-Business/conf"
	"github.com/gin-gonic/gin"
)

// ServerStart 服务启动
func ServerStart(eng *gin.Engine, mode string) {
	// 加载配置文件并设置全局配置常量
	conf.LoadConf(mode)
	eng.Use(gin.Recovery())

	// 配置跨域
	if conf.GlobalConf.System.Security.Cors.Enabled {
		eng.Use(middleware.CorsMiddleware(
			conf.GlobalConf.System.Security.Cors.AllowOrigins,
			conf.GlobalConf.System.Security.Cors.AllowMethods,
		))
	}

	// 创建日志实例
	err := logger.NewLogger()
	if err != nil {
		// 如果日志初始化失败，打印错误并终止
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	routers.Register(eng) // 注册路由
	db.MysqlInit()        // 初始化MySQL连接（日志也是在这里初始化）
	redis.RedisInit()     // 初始化Redis连接
	cron.StartForever()   // 执行消费者任务

	// 创建gpt实例
	_, err = gpt.NewGptClient(nil)
	if err != nil {
		log.Fatalf("Failed to initialize gpt: %v", err)
	}

	server := &http.Server{
		Addr:         ":" + conf.GlobalConf.System.Addr,
		Handler:      eng,
		ReadTimeout:  conf.GlobalConf.System.Http.ReadTimeout,
		WriteTimeout: conf.GlobalConf.System.Http.WriteTimeout,
		IdleTimeout:  conf.GlobalConf.System.Http.IdleTimeout,
	}

	// 映射静态文件
	eng.Static("/files", conf.GlobalConf.File.TaskPath)

	// 启动服务
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Logger.Info("server start err", err)
		}
	}()
	logger.Logger.Infof("=== %s(%s) is starting ===", conf.GlobalConf.System.Name,
		conf.GlobalConf.System.Version)

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Logger.Println("Shutdown Server ...")

	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Info("Server Shutdown:", err)
	}
}

// ServerExit 服务退出
func ServerExit(server *gin.Engine) {
	logger.Logger.Infof("=== %s(%s) is exit ===", conf.GlobalConf.System.Name,
		conf.GlobalConf.System.Version)
}
