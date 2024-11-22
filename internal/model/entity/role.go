package entity

import (
	"time"
)

/******sql******
CREATE TABLE `role` (
  `id` varchar(36) NOT NULL COMMENT '角色唯一标识，使用UUID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色表'
******sql******/
// Role 角色表
type Role struct {
	ID          string    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                            // 角色唯一标识，使用UUID
	Name        string    `gorm:"unique;column:name;type:varchar(50);not null" json:"name"`                            // 角色名称
	Description string    `gorm:"column:description;type:varchar(255);default:null" json:"description"`                // 角色描述
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"` // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"` // 更新时间
}

// TableName get sql table name.获取数据库表名
func (m *Role) TableName() string {
	return "role"
}

// RoleColumns get sql column name.获取数据库列名
var RoleColumns = struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}
