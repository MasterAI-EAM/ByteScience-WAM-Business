package entity

/******sql******
CREATE TABLE `role_paths` (
  `role_id` char(36) NOT NULL COMMENT '角色ID',
  `path_id` char(36) NOT NULL COMMENT '路径ID',
  PRIMARY KEY (`role_id`,`path_id`),
  KEY `path_id` (`path_id`),
  CONSTRAINT `role_paths_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `role_paths_ibfk_2` FOREIGN KEY (`path_id`) REFERENCES `paths` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色接口中间表'
******sql******/
// RolePaths 角色接口中间表
type RolePaths struct {
	RoleID string `gorm:"primaryKey;column:role_id;type:char(36);not null" json:"roleId"`               // 角色ID
	PathID string `gorm:"primaryKey;index:path_id;column:path_id;type:char(36);not null" json:"pathId"` // 路径ID
}

// TableName get sql table name.获取数据库表名
func (m *RolePaths) TableName() string {
	return "role_paths"
}

// RolePathsColumns get sql column name.获取数据库列名
var RolePathsColumns = struct {
	RoleID string
	PathID string
}{
	RoleID: "role_id",
	PathID: "path_id",
}
