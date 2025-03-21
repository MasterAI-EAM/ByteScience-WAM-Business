package entity

/******sql******
CREATE TABLE `user_permissions` (
  `user_id` char(36) NOT NULL COMMENT '用户ID',
  `path_id` char(36) NOT NULL COMMENT '路径ID',
  PRIMARY KEY (`user_id`,`path_id`),
  KEY `path_id` (`path_id`),
  CONSTRAINT `user_permissions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `user_permissions_ibfk_2` FOREIGN KEY (`path_id`) REFERENCES `paths` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户权限预计算表'
******sql******/
// UserPermissions 用户权限预计算表
type UserPermissions struct {
	UserID string `gorm:"primaryKey;column:user_id;type:char(36);not null"`               // 用户ID
	PathID string `gorm:"primaryKey;index:path_id;column:path_id;type:char(36);not null"` // 路径ID
}

// TableName get sql table name.获取数据库表名
func (m *UserPermissions) TableName() string {
	return "user_permissions"
}

// UserPermissionsColumns get sql column name.获取数据库列名
var UserPermissionsColumns = struct {
	UserID string
	PathID string
}{
	UserID: "user_id",
	PathID: "path_id",
}
