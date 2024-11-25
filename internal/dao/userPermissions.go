package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"context"
	"fmt"
	"gorm.io/gorm"
)

// UserPermissionDao 用户权限关联表数据访问对象
type UserPermissionDao struct{}

// NewUserPermissionDao 创建 UserPermissionDao 实例
func NewUserPermissionDao() *UserPermissionDao {
	return &UserPermissionDao{}
}

// RemoveByUserIDsTx 删除指定用户的权限记录
func (dao *UserPermissionDao) RemoveByUserIDsTx(ctx context.Context, tx *gorm.DB, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	// 使用结构体中的列名来替代硬编码的字段名
	if err := tx.WithContext(ctx).
		Where(entity.UserPermissionsColumns.UserID+" IN ?", userIDs).
		Delete(&entity.UserPermissions{}).Error; err != nil {
		return fmt.Errorf("failed to remove user permissions: %w", err)
	}

	return nil
}

// UpdateUserPermissionsTx 更新用户权限表
func (dao *UserPermissionDao) UpdateUserPermissionsTx(ctx context.Context, tx *gorm.DB, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	// 删除 user_permissions 中对应用户的所有权限
	if err := tx.WithContext(ctx).
		Where(entity.UserPermissionsColumns.UserID+" IN ?", userIDs).
		Delete(&entity.UserPermissions{}).Error; err != nil {
		return fmt.Errorf("failed to delete user permissions: %w", err)
	}

	// 查询 user_roles 中所有关联的角色 ID
	var userRoles []entity.UserRoles
	if err := tx.WithContext(ctx).
		Where(entity.UserRolesColumns.UserID+" IN ?", userIDs).
		Find(&userRoles).Error; err != nil {
		return fmt.Errorf("failed to fetch user roles: %w", err)
	}

	// 提取所有 roleID
	roleIDSet := make(map[string]struct{})
	for _, ur := range userRoles {
		roleIDSet[ur.RoleID] = struct{}{}
	}
	var roleIDs []string
	for roleID := range roleIDSet {
		roleIDs = append(roleIDs, roleID)
	}

	// 查询 role_paths 中所有路径 ID
	var rolePaths []entity.RolePaths
	if err := tx.WithContext(ctx).
		Where(entity.RolePathsColumns.RoleID+" IN ?", roleIDs).
		Find(&rolePaths).Error; err != nil {
		return fmt.Errorf("failed to fetch role paths: %w", err)
	}

	// 构造角色与路径的映射
	rolePathMap := make(map[string][]string)
	for _, rp := range rolePaths {
		rolePathMap[rp.RoleID] = append(rolePathMap[rp.RoleID], rp.PathID)
	}

	// 构造用户与路径的映射
	userPathMap := make(map[string]map[string]struct{})
	for _, ur := range userRoles {
		if paths, exists := rolePathMap[ur.RoleID]; exists {
			if _, ok := userPathMap[ur.UserID]; !ok {
				userPathMap[ur.UserID] = make(map[string]struct{})
			}
			for _, pathID := range paths {
				userPathMap[ur.UserID][pathID] = struct{}{}
			}
		}
	}

	// 准备批量插入的数据
	var userPermissions []entity.UserPermissions
	for userID, pathSet := range userPathMap {
		for pathID := range pathSet {
			userPermissions = append(userPermissions, entity.UserPermissions{
				UserID: userID,
				PathID: pathID,
			})
		}
	}

	// 批量插入新的权限记录
	if len(userPermissions) > 0 {
		if err := tx.WithContext(ctx).CreateInBatches(&userPermissions, 300).Error; err != nil {
			return fmt.Errorf("failed to insert user permissions: %w", err)
		}
	}

	return nil
}
