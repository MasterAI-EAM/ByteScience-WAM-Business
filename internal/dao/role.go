package dao

import (
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"
	"time"

	"ByteScience-WAM-Business/internal/model/entity"
	"gorm.io/gorm"
)

// RoleDao 数据访问对象，封装角色相关操作
type RoleDao struct{}

// NewRoleDao 创建一个新的 RoleDao 实例
func NewRoleDao() *RoleDao {
	return &RoleDao{}
}

// Insert 插入角色记录
func (rd *RoleDao) Insert(ctx context.Context, role *entity.Roles) error {
	return db.Client.WithContext(ctx).Create(role).Error
}

// InsertTx 在事务中插入角色
func (rd *RoleDao) InsertTx(ctx context.Context, tx *gorm.DB, role *entity.Roles) error {
	return tx.WithContext(ctx).Create(role).Error
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

// Update 更新角色信息
func (rd *RoleDao) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return db.Client.WithContext(ctx).
		Model(&entity.Roles{}).
		Where(entity.RolesColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// UpdateTx 在事务中更新角色信息
func (rd *RoleDao) UpdateTx(ctx context.Context, tx *gorm.DB, id string, updates map[string]interface{}) error {
	return tx.WithContext(ctx).
		Model(&entity.Roles{}).
		Where(entity.RolesColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// SoftDeleteByID 软删除角色记录
func (rd *RoleDao) SoftDeleteByID(ctx context.Context, id string) error {
	return db.Client.WithContext(ctx).
		Model(&entity.Roles{}).
		Where(entity.RolesColumns.ID+" = ?", id).
		Update(entity.RolesColumns.DeletedAt, time.Now()).
		Error
}

// SoftDeleteByIDTx 软删除角色记录
func (rd *RoleDao) SoftDeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error {
	return tx.WithContext(ctx).
		Model(&entity.Roles{}).
		Where(entity.RolesColumns.ID+" = ?", id).
		Update(entity.RolesColumns.DeletedAt, time.Now()).
		Error
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

// UpdateStatus 更新角色的状态
func (rd *RoleDao) UpdateStatus(ctx context.Context, id string, status int) error {
	return db.Client.WithContext(ctx).
		Model(&entity.Roles{}).
		Where(entity.RolesColumns.ID+" = ?", id).
		Update(entity.RolesColumns.Status, status).
		Error
}
