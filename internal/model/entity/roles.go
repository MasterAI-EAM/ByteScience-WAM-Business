package entity

import (
	"time"
)

/******sql******
CREATE TABLE `roles` (
  `id` char(36) NOT NULL COMMENT '角色ID',
  `name` varchar(128) NOT NULL COMMENT '角色名称',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint DEFAULT '1' COMMENT '状态: 1=启用, 0=禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name_deleted` (`name`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色表'
******sql******/
// Roles 角色表
type Roles struct {
	ID          string    `gorm:"primaryKey;column:id;type:char(36);not null" json:"id"`                                          // 角色ID
	Name        string    `gorm:"uniqueIndex:unique_name_deleted;column:name;type:varchar(128);not null" json:"name"`             // 角色名称
	Description string    `gorm:"column:description;type:varchar(255);default:null" json:"description"`                           // 角色描述
	Status      int8      `gorm:"column:status;type:tinyint;default:null;default:1" json:"status"`                                // 状态: 1=启用, 0=禁用
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:null;default:CURRENT_TIMESTAMP" json:"createdAt"`       // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;default:null;default:CURRENT_TIMESTAMP" json:"updatedAt"`       // 更新时间
	DeletedAt   time.Time `gorm:"uniqueIndex:unique_name_deleted;column:deleted_at;type:timestamp;default:null" json:"deletedAt"` // 软删除时间
}

// TableName get sql table name.获取数据库表名
func (m *Roles) TableName() string {
	return "roles"
}

// RolesColumns get sql column name.获取数据库列名
var RolesColumns = struct {
	ID          string
	Name        string
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}
