package entity

import (
	"time"
)

/******sql******
CREATE TABLE `users` (
  `id` varchar(36) NOT NULL COMMENT '唯一标识',
  `username` varchar(128) NOT NULL COMMENT '用户名',
  `nickname` varchar(128) DEFAULT NULL COMMENT '昵称',
  `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '加密后的密码',
  `email` varchar(256) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(32) DEFAULT NULL COMMENT '手机号码',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态(1: 启用, 0: 禁用)',
  `remark` varchar(256) DEFAULT NULL COMMENT '备注',
  `last_login_at` datetime DEFAULT NULL COMMENT '上次登录时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`,`deleted_at`),
  UNIQUE KEY `email_deleted_at` (`email`,`deleted_at`),
  UNIQUE KEY `phone_deleted_at` (`phone`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表'
******sql******/
// Users 用户表
type Users struct {
	ID          string     `gorm:"primaryKey;column:id;type:varchar(36);not null"`                                                                  // 唯一标识
	Username    string     `gorm:"uniqueIndex:username;column:username;type:varchar(128);not null"`                                                 // 用户名
	Nickname    *string    `gorm:"column:nickname;type:varchar(128)"`                                                                               // 昵称
	Password    string     `gorm:"column:password;type:varchar(64);not null"`                                                                       // 加密后的密码
	Email       *string    `gorm:"uniqueIndex:email_deleted_at;column:email;type:varchar(256)"`                                                     // 邮箱
	Phone       *string    `gorm:"uniqueIndex:phone_deleted_at;column:phone;type:varchar(32)"`                                                      // 手机号码
	Status      int8       `gorm:"column:status;type:tinyint;not null;default:1"`                                                                   // 状态(1: 启用, 0: 禁用)
	Remark      *string    `gorm:"column:remark;type:varchar(256)"`                                                                                 // 备注
	LastLoginAt *time.Time `gorm:"column:last_login_at;type:datetime"`                                                                              // 上次登录时间
	CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`                                              // 创建时间
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`                                              // 更新时间
	DeletedAt   *time.Time `gorm:"uniqueIndex:username;uniqueIndex:email_deleted_at;uniqueIndex:phone_deleted_at;column:deleted_at;type:timestamp"` // 软删除时间
}

// TableName get sql table name.获取数据库表名
func (m *Users) TableName() string {
	return "users"
}

// UsersColumns get sql column name.获取数据库列名
var UsersColumns = struct {
	ID          string
	Username    string
	Nickname    string
	Password    string
	Email       string
	Phone       string
	Status      string
	Remark      string
	LastLoginAt string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}{
	ID:          "id",
	Username:    "username",
	Nickname:    "nickname",
	Password:    "password",
	Email:       "email",
	Phone:       "phone",
	Status:      "status",
	Remark:      "remark",
	LastLoginAt: "last_login_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}
