package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"gorm.io/gorm"
)

// RolePathDao 角色路径数据访问对象
type RolePathDao struct{}

// NewRolePathDao 创建 RolePathDao 实例
func NewRolePathDao() *RolePathDao {
	return &RolePathDao{}
}

// InsertBatchTx 在事务中批量插入角色路径关系
func (rpd *RolePathDao) InsertBatchTx(ctx context.Context, tx *gorm.DB, rolePaths []*entity.RolePaths) error {
	return tx.WithContext(ctx).Create(&rolePaths).Error
}

// Assign 分配路径给角色
func (rpd *RolePathDao) Assign(ctx context.Context, roleID, pathID string) error {
	rolePath := &entity.RolePaths{
		RoleID: roleID,
		PathID: pathID,
	}
	return db.Client.WithContext(ctx).Create(rolePath).Error
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

// Remove 移除角色的路径
func (rpd *RolePathDao) Remove(ctx context.Context, roleID, pathID string) error {
	return db.Client.WithContext(ctx).
		Delete(&entity.RolePaths{}, "role_id = ? AND path_id = ?", roleID, pathID).
		Error
}

// RemoveTx 在事务中移除角色的路径
func (rpd *RolePathDao) RemoveTx(ctx context.Context, tx *gorm.DB, roleID, pathID string) error {
	return tx.WithContext(ctx).
		Delete(&entity.RolePaths{}, "role_id = ? AND path_id = ?", roleID, pathID).
		Error
}

// RemoveByRoleIDTx 在事务中根据角色id移除角色的路径
func (rpd *RolePathDao) RemoveByRoleIDTx(ctx context.Context, tx *gorm.DB, roleID string) error {
	return tx.WithContext(ctx).
		Delete(&entity.RolePaths{}, "role_id = ?", roleID).
		Error
}

// GetByPathID 根据路径ID获取拥有该路径的角色列表
func (rpd *RolePathDao) GetByPathID(ctx context.Context, pathID string) ([]*entity.Roles, error) {
	var roles []*entity.Roles
	err := db.Client.WithContext(ctx).
		Select("roles.*").
		Joins("JOIN role_paths ON role_paths.role_id = roles.id").
		Where("role_paths.path_id = ?", pathID).
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
