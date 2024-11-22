package entity

import (
	"time"
)

/******sql******
CREATE TABLE `admins` (
  `id` varchar(36) NOT NULL COMMENT '唯一标识',
  `username` varchar(128) NOT NULL COMMENT '用户名',
  `nickname` varchar(128) DEFAULT NULL COMMENT '昵称',
  `password` varchar(64) NOT NULL COMMENT '加密后的密码',
  `email` varchar(256) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(32) DEFAULT NULL COMMENT '手机号码',
  `remark` varchar(256) DEFAULT NULL COMMENT '备注',
  `last_login_at` datetime DEFAULT NULL COMMENT '上次登录时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`,`deleted_at`),
  UNIQUE KEY `email_deleted_at` (`email`,`deleted_at`),
  UNIQUE KEY `phone_deleted_at` (`phone`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员表'
******sql******/
// Admins 管理员表
type Admins struct {
	ID          string    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                                                                                      // 唯一标识
	Username    string    `gorm:"uniqueIndex:username;column:username;type:varchar(128);not null" json:"username"`                                                               // 用户名
	Nickname    string    `gorm:"column:nickname;type:varchar(128);default:null" json:"nickname"`                                                                                // 昵称
	Password    string    `gorm:"column:password;type:varchar(64);not null" json:"password"`                                                                                     // 加密后的密码
	Email       string    `gorm:"uniqueIndex:email_deleted_at;column:email;type:varchar(256);default:null" json:"email"`                                                         // 邮箱
	Phone       string    `gorm:"uniqueIndex:phone_deleted_at;column:phone;type:varchar(32);default:null" json:"phone"`                                                          // 手机号码
	Remark      string    `gorm:"column:remark;type:varchar(256);default:null" json:"remark"`                                                                                    // 备注
	LastLoginAt time.Time `gorm:"column:last_login_at;type:datetime;default:null" json:"lastLoginAt"`                                                                            // 上次登录时间
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`                                                           // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`                                                           // 更新时间
	DeletedAt   time.Time `gorm:"uniqueIndex:username;uniqueIndex:email_deleted_at;uniqueIndex:phone_deleted_at;column:deleted_at;type:timestamp;default:null" json:"deletedAt"` // 软删除时间
}

// TableName get sql table name.获取数据库表名
func (m *Admins) TableName() string {
	return "admins"
}

// AdminsColumns get sql column name.获取数据库列名
var AdminsColumns = struct {
	ID          string
	Username    string
	Nickname    string
	Password    string
	Email       string
	Phone       string
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
	Remark:      "remark",
	LastLoginAt: "last_login_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}
