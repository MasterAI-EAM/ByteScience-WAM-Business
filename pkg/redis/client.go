package redis

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/pkg/logger" // 引入已经封装好的日志包
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// Client 客户端实例
var Client *redis.Client

// RedisInit 初始化 Redis 客户端
func RedisInit() {
	// 获取全局 Logger 实例
	config := conf.GlobalConf.Redis

	// 创建 Redis 客户端配置
	redisOptions := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Db,
		PoolSize:     config.PoolSize,
		MaxRetries:   5,
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	// 创建 Redis 客户端
	client := redis.NewClient(redisOptions)

	// 如果启用 Redis 操作日志，添加钩子
	if config.LogEnabled {
		hook := &redisHook{log: logger.GormLogger} // 使用日志构造函数
		client.AddHook(hook)
	}

	// 检查 Redis 是否连接成功
	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Logger.Errorf("Redis connection failed: %v", err)
		panic(err)
	}

	Client = client
	logger.Logger.Infof("Redis connection established successfully")
}
