package entity

import (
	"time"
)

/******sql******
CREATE TABLE `menus` (
  `id` char(36) NOT NULL COMMENT '菜单ID',
  `parent_id` char(36) DEFAULT NULL COMMENT '父菜单ID，指向上一级菜单',
  `name` varchar(128) NOT NULL COMMENT '菜单名称',
  `sort` int DEFAULT '0' COMMENT '排序字段',
  `status` tinyint DEFAULT '1' COMMENT '状态: 1=启用, 0=禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `parent_id` (`parent_id`),
  CONSTRAINT `menus_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜单表'
******sql******/
// Menus 菜单表
type Menus struct {
	ID        string    `gorm:"primaryKey;column:id;type:char(36);not null" json:"id"`                                    // 菜单ID
	ParentID  string    `gorm:"index:parent_id;column:parent_id;type:char(36);default:null" json:"parentId"`              // 父菜单ID，指向上一级菜单
	Name      string    `gorm:"column:name;type:varchar(128);not null" json:"name"`                                       // 菜单名称
	Sort      int       `gorm:"column:sort;type:int;default:null;default:0" json:"sort"`                                  // 排序字段
	Status    int8      `gorm:"column:status;type:tinyint;default:null;default:1" json:"status"`                          // 状态: 1=启用, 0=禁用
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:null;default:CURRENT_TIMESTAMP" json:"createdAt"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:null;default:CURRENT_TIMESTAMP" json:"updatedAt"` // 更新时间
	DeletedAt time.Time `gorm:"column:deleted_at;type:timestamp;default:null" json:"deletedAt"`                           // 软删除时间
}

// TableName get sql table name.获取数据库表名
func (m *Menus) TableName() string {
	return "menus"
}

// MenusColumns get sql column name.获取数据库列名
var MenusColumns = struct {
	ID        string
	ParentID  string
	Name      string
	Sort      string
	Status    string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}{
	ID:        "id",
	ParentID:  "parent_id",
	Name:      "name",
	Sort:      "sort",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}
