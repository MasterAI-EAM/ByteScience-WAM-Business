package dao

// UserPermissionDao 用户权限关联表数据访问对象
type UserPermissionDao struct{}

// NewUserPermissionDao 创建 UserPermissionDao 实例
func NewUserPermissionDao() *UserPermissionDao {
	return &UserPermissionDao{}
}
