package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"
	"gorm.io/gorm"
)

// RoleDao 数据访问对象，封装角色相关操作
type RoleDao struct{}

// NewRoleDao 创建一个新的 RoleDao 实例
func NewRoleDao() *RoleDao {
	return &RoleDao{}
}

// GetByID 根据 ID 获取角色
func (rd *RoleDao) GetByID(ctx context.Context, id string) (*entity.Roles, error) {
	var role entity.Roles
	err := db.Client.WithContext(ctx).
		Where(entity.RolesColumns.ID+" = ?", id).
		Where(entity.RolesColumns.DeletedAt + " IS NULL").
		First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// GetByName 根据名称获取角色
func (rd *RoleDao) GetByName(ctx context.Context, name string) (*entity.Roles, error) {
	var role entity.Roles
	err := db.Client.WithContext(ctx).
		Where(entity.RolesColumns.Name+" = ?", name).
		Where(entity.RolesColumns.DeletedAt + " IS NULL").
		First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// Query 分页查询角色
func (rd *RoleDao) Query(ctx context.Context, page int, pageSize int, filters map[string]interface{}) ([]*entity.
	Roles, int64, error) {
	var (
		roles []*entity.Roles
		total int64
	)

	query := db.Client.WithContext(ctx).Model(&entity.Roles{}).Where(entity.RolesColumns.DeletedAt + " IS NULL")

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			if key == entity.RolesColumns.Name {
				query = query.Where(key+" LIKE ?", value.(string)+"%")
			} else {
				query = query.Where(key+" = ?", value)
			}
		}
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Scopes(db.PageScope(page, pageSize)).
		Order(entity.RolesColumns.CreatedAt + " DESC").
		Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
