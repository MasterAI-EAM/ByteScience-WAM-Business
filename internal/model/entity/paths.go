package entity

import "time"

/******sql******
CREATE TABLE `paths` (
  `id` char(36) NOT NULL COMMENT '路径ID',
  `path` varchar(256) NOT NULL COMMENT '路由路径',
  `method` enum('GET','POST','PUT','DELETE') NOT NULL COMMENT 'HTTP 方法',
  `description` varchar(255) DEFAULT NULL COMMENT '路径描述',
  `menu_id` char(36) NOT NULL COMMENT '菜单ID，指向menus表的ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_path_method` (`path`,`method`,`deleted_at`),
  KEY `paths_ibfk_1` (`menu_id`),
  CONSTRAINT `paths_ibfk_1` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='接口表'
******sql******/
// Paths 接口表
type Paths struct {
	ID          string     `gorm:"primaryKey;column:id;type:char(36);not null"`                                                  // 路径ID
	Path        string     `gorm:"uniqueIndex:unique_path_method;column:path;type:varchar(256);not null"`                        // 路由路径
	Method      string     `gorm:"uniqueIndex:unique_path_method;column:method;type:enum('GET','POST','PUT','DELETE');not null"` // HTTP 方法
	Description *string    `gorm:"column:description;type:varchar(255)"`                                                         // 路径描述
	MenuID      string     `gorm:"index:paths_ibfk_1;column:menu_id;type:char(36);not null"`                                     // 菜单ID，指向menus表的ID
	CreatedAt   *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`                                   // 创建时间
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`                                   // 更新时间
	DeletedAt   *time.Time `gorm:"uniqueIndex:unique_path_method;column:deleted_at;type:timestamp"`                              // 软删除时间
}

// TableName get sql table name.获取数据库表名
func (m *Paths) TableName() string {
	return "paths"
}

// PathsColumns get sql column name.获取数据库列名
var PathsColumns = struct {
	ID          string
	Path        string
	Method      string
	Description string
	MenuID      string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}{
	ID:          "id",
	Path:        "path",
	Method:      "method",
	Description: "description",
	MenuID:      "menu_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}
