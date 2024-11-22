package db

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/pkg/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Client 是数据库连接的全局变量
var Client *gorm.DB

// MysqlInit 初始化 MySQL 数据库连接
func MysqlInit() (err error) {
	config := conf.GlobalConf.Mysql

	// 从配置文件获取数据库配置信息
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Db)

	// 连接数据库
	Client, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.GormLogger, // 使用自定义的 logger
	})
	if err != nil {
		// 如果数据库连接失败，打印错误并终止
		logger.Logger.Fatalf("Failed to connect to the database: %v", err)
	}

	// 可选：自动迁移数据库结构
	err = Client.AutoMigrate()
	if err != nil {
		// 如果自动迁移失败，打印错误并终止
		logger.Logger.Fatalf("Failed to migrate database: %v", err)
	}

	logger.Logger.Info("=== Mysql initialization successful ===")
	return
}

// Close 关闭数据库连接
func Close() {
	sqlDB, err := Client.DB()
	if err != nil {
		logger.Logger.Fatalf("Failed to get database connection: %v", err)
	}
	sqlDB.Close()
}
