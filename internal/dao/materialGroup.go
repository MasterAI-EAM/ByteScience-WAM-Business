package dao

import (
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"

	"ByteScience-WAM-Business/internal/model/entity"
	"gorm.io/gorm"
)

// MaterialGroupDao 数据访问对象，封装材料表相关操作
type MaterialGroupDao struct{}

// NewMaterialGroupDao 创建一个新的 MaterialGroupDao 实例
func NewMaterialGroupDao() *MaterialGroupDao {
	return &MaterialGroupDao{}
}

// Insert 插入材料记录
func (md *MaterialGroupDao) Insert(ctx context.Context, tx *gorm.DB, Materials *entity.MaterialGroups) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).Create(Materials).Error
}

// Update 更新材料信息
func (md *MaterialGroupDao) Update(ctx context.Context, tx *gorm.DB, id string, updates map[string]interface{}) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.MaterialGroups{}).
		Where(entity.MaterialGroupsColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// DeleteByID 删除材料记录
func (md *MaterialGroupDao) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Where(entity.MaterialGroupsColumns.ID+" = ?", id).
		Delete(&entity.MaterialGroups{}).Error
}

// DeleteByExperimentID 删除实验下的彩料组
func (md *MaterialGroupDao) DeleteByExperimentID(ctx context.Context, tx *gorm.DB, experimentID string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Where(entity.MaterialGroupsColumns.ExperimentID, experimentID).
		Delete(&entity.MaterialGroups{}).Error
}

// GetByID 根据 ID 获取材料
func (md *MaterialGroupDao) GetByID(ctx context.Context, id string) (*entity.MaterialGroups, error) {
	var materialGroup entity.MaterialGroups
	err := db.Client.WithContext(ctx).
		Where(entity.MaterialGroupsColumns.ID+" = ?", id).
		First(&materialGroup).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &materialGroup, err
}

// GetByRecipeID 根据配方ID获取材料
func (md *MaterialGroupDao) GetByRecipeID(ctx context.Context, experimentID string) ([]*entity.MaterialGroups, error) {
	var materialGroups []*entity.MaterialGroups
	err := db.Client.WithContext(ctx).
		Where(entity.MaterialGroupsColumns.ExperimentID+" = ?", experimentID).
		Order(entity.MaterialGroupsColumns.Sort + " ASC").
		Find(&materialGroups).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return materialGroups, err
}

// GetGroupIdListByRecipeID 根据组id列表
func (md *MaterialGroupDao) GetGroupIdListByRecipeID(ctx context.Context, experimentID string) ([]string, error) {
	var groupIdList []string
	err := db.Client.WithContext(ctx).
		Where(entity.MaterialGroupsColumns.ExperimentID+" = ?", experimentID).
		Find(entity.MaterialGroupsColumns.ID, &groupIdList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return groupIdList, err
}
