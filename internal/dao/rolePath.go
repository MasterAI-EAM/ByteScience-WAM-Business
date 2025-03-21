package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/pkg/db"
	"context"
)

// RolePathDao 角色路径数据访问对象
type RolePathDao struct{}

// NewRolePathDao 创建 RolePathDao 实例
func NewRolePathDao() *RolePathDao {
	return &RolePathDao{}
}

// GetByRoleID 根据角色ID获取路径列表
func (rpd *RolePathDao) GetByRoleID(ctx context.Context, roleID string) ([]*entity.Paths, error) {
	var paths []*entity.Paths
	err := db.Client.WithContext(ctx).
		Select("paths.*").
		Joins("JOIN role_paths ON role_paths.path_id = paths.id").
		Where("role_paths.role_id = ?", roleID).
		Where(entity.PathsColumns.DeletedAt + " IS NULL").
		Find(&paths).Error
	return paths, err
}

// GetByPathID 根据路径ID获取拥有该路径的角色列表
func (rpd *RolePathDao) GetByPathID(ctx context.Context, pathID string) ([]*entity.Roles, error) {
	var roles []*entity.Roles
	err := db.Client.WithContext(ctx).
		Select("roles.*").
		Joins("JOIN role_paths ON role_paths.role_id = roles.id").
		Where("role_paths.path_id = ?", pathID).
		Where(entity.RolesColumns.DeletedAt + " IS NULL").
		Find(&roles).Error
	return roles, err
}

// Query 分页查询角色路径中间表
func (rpd *RolePathDao) Query(ctx context.Context, page int, pageSize int, filters map[string]interface{}) ([]*entity.RolePaths, int64, error) {
	var (
		rolePaths []*entity.RolePaths
		total     int64
	)

	query := db.Client.WithContext(ctx).Model(&entity.RolePaths{})

	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(db.PageScope(page, pageSize)).
		Find(&rolePaths).Error; err != nil {
		return nil, 0, err
	}

	return rolePaths, total, nil
}
