package entity

/******sql******
CREATE TABLE `user_roles` (
  `user_id` char(36) NOT NULL COMMENT '用户ID',
  `role_id` char(36) NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户与角色关联表'
******sql******/
// UserRoles 用户与角色关联表
type UserRoles struct {
	UserID string `gorm:"primaryKey;column:user_id;type:char(36);not null"`               // 用户ID
	RoleID string `gorm:"primaryKey;index:role_id;column:role_id;type:char(36);not null"` // 角色ID
}

// TableName get sql table name.获取数据库表名
func (m *UserRoles) TableName() string {
	return "user_roles"
}

// UserRolesColumns get sql column name.获取数据库列名
var UserRolesColumns = struct {
	UserID string
	RoleID string
}{
	UserID: "user_id",
	RoleID: "role_id",
}
