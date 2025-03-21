package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"
	"gorm.io/gorm"
)

// MenuDao 菜单数据访问对象
type MenuDao struct{}

// NewMenuDao 创建 MenuDao 实例
func NewMenuDao() *MenuDao {
	return &MenuDao{}
}

// GetByID 根据 ID 获取菜单
func (md *MenuDao) GetByID(ctx context.Context, id string) (*entity.Menus, error) {
	var menu entity.Menus
	err := db.Client.WithContext(ctx).
		Where(entity.MenusColumns.ID+" = ?", id).
		Where(entity.MenusColumns.DeletedAt + " IS NULL").
		First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &menu, err
}

// GetByParentID 根据父菜单 ID 获取子菜单列表
func (md *MenuDao) GetByParentID(ctx context.Context, parentID string) ([]*entity.Menus, error) {
	var menus []*entity.Menus
	err := db.Client.WithContext(ctx).
		Where(entity.MenusColumns.ParentID+" = ?", parentID).
		Where(entity.MenusColumns.DeletedAt + " IS NULL").
		Order(entity.MenusColumns.Sort + " ASC").
		Find(&menus).Error
	return menus, err
}

// Query 分页查询菜单
func (md *MenuDao) Query(ctx context.Context, page int, pageSize int, filters map[string]interface{}) ([]*entity.Menus, int64, error) {
	var (
		menus []*entity.Menus
		total int64
	)

	query := db.Client.WithContext(ctx).Model(&entity.Menus{}).Where(entity.MenusColumns.DeletedAt + " IS NULL")

	for key, value := range filters {
		if value != nil && value != "" {
			if key == entity.MenusColumns.Name {
				query = query.Where(key+" LIKE ?", value.(string)+"%")
			} else {
				query = query.Where(key+" = ?", value)
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(db.PageScope(page, pageSize)).
		Order(entity.MenusColumns.Sort + " ASC").
		Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

// GetAll 获取所有菜单
func (md *MenuDao) GetAll(ctx context.Context) ([]*entity.Menus, error) {
	var menus []*entity.Menus
	err := db.Client.WithContext(ctx).
		Where(entity.MenusColumns.DeletedAt + " IS NULL").
		Order(entity.MenusColumns.Sort + " ASC").
		Find(&menus).Error
	return menus, err
}
