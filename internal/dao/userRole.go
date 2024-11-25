package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"fmt"
	"gorm.io/gorm"
)

// UserRoleDao 用户角色关联表数据访问对象
type UserRoleDao struct{}

// NewUserRoleDao 创建 UserRoleDao 实例
func NewUserRoleDao() *UserRoleDao {
	return &UserRoleDao{}
}

// InsertBatchTx 在事务中批量插入用户角色关联
func (urd *UserRoleDao) InsertBatchTx(ctx context.Context, tx *gorm.DB, userRoles []*entity.UserRoles) error {
	return tx.WithContext(ctx).CreateInBatches(&userRoles, 300).Error
}

// Assign 给用户分配角色
func (urd *UserRoleDao) Assign(ctx context.Context, userID, roleID string) error {
	userRole := &entity.UserRoles{
		UserID: userID,
		RoleID: roleID,
	}
	return db.Client.WithContext(ctx).Create(userRole).Error
}

// GetRolesByUserID 根据用户ID获取角色列表
func (urd *UserRoleDao) GetRolesByUserID(ctx context.Context, userID string) ([]*entity.Roles, error) {
	var roles []*entity.Roles
	err := db.Client.WithContext(ctx).
		Select("roles.*").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("roles.deleted_at" + " IS NULL").
		Find(&roles).Error
	return roles, err
}

// Remove 移除用户的角色
func (urd *UserRoleDao) Remove(ctx context.Context, userID, roleID string) error {
	return db.Client.WithContext(ctx).
		Delete(&entity.UserRoles{}, "user_id = ? AND role_id = ?", userID, roleID).
		Error
}

// RemoveByUserIDTx 在事务中根据用户ID移除所有关联角色
func (urd *UserRoleDao) RemoveByUserIDTx(ctx context.Context, tx *gorm.DB, userID string) error {
	return tx.WithContext(ctx).
		Delete(&entity.UserRoles{}, "user_id = ?", userID).
		Error
}

// RemoveByRoleIDTx 在事务中根据角色ID移除所有关联用户
func (urd *UserRoleDao) RemoveByRoleIDTx(ctx context.Context, tx *gorm.DB, roleID string) error {
	return tx.WithContext(ctx).
		Delete(&entity.UserRoles{}, "role_id = ?", roleID).
		Error
}

// GetUsersByRoleID 根据角色ID获取用户列表
func (urd *UserRoleDao) GetUsersByRoleID(ctx context.Context, roleID string) ([]*entity.Users, error) {
	var users []*entity.Users
	err := db.Client.WithContext(ctx).
		Select("users.*").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Where("users.deleted_at" + " IS NULL").
		Find(&users).Error
	return users, err
}

// Query 分页查询用户角色中间表
func (urd *UserRoleDao) Query(ctx context.Context, page int, pageSize int, filters map[string]interface{}) ([]*entity.UserRoles, int64, error) {
	var (
		userRoles []*entity.UserRoles
		total     int64
	)

	query := db.Client.WithContext(ctx).Model(&entity.UserRoles{})

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(db.PageScope(page, pageSize)).
		Find(&userRoles).Error; err != nil {
		return nil, 0, err
	}

	return userRoles, total, nil
}

// GetUserIDsByRoleIDTx 获取与指定角色关联的用户 ID 列表
func (urd *UserRoleDao) GetUserIDsByRoleIDTx(ctx context.Context, tx *gorm.DB, roleID string) ([]string, error) {
	var userIDs []string

	// 查询 user_roles 表，获取所有关联的 user_id
	err := tx.WithContext(ctx).
		Table("user_roles").
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user IDs for role %d: %w", roleID, err)
	}

	return userIDs, nil
}
