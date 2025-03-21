package dao

import (
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"

	"ByteScience-WAM-Business/internal/model/entity"
	"gorm.io/gorm"
)

// MaterialDao 数据访问对象，封装材料表相关操作
type MaterialDao struct{}

// NewMaterialDao 创建一个新的 MaterialDao 实例
func NewMaterialDao() *MaterialDao {
	return &MaterialDao{}
}

// Insert 插入材料记录
func (md *MaterialDao) Insert(ctx context.Context, Materials *entity.Materials) error {
	return db.Client.WithContext(ctx).Create(Materials).Error
}

// GetByID 根据 ID 获取材料
func (md *MaterialDao) GetByID(ctx context.Context, id string) (*entity.Materials, error) {
	var Materials entity.Materials
	err := db.Client.WithContext(ctx).
		Where(entity.MaterialsColumns.ID+" = ?", id).
		First(&Materials).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &Materials, err
}

// GetByGroupID 根据实验材料组 ID 获取材料
func (md *MaterialDao) GetByGroupID(ctx context.Context, groupID string) ([]*entity.Materials, error) {
	var Materialss []*entity.Materials
	err := db.Client.WithContext(ctx).
		Where(entity.MaterialsColumns.MaterialGroupID+" = ?", groupID).
		Order(entity.MaterialsColumns.MaterialName + " ASC").
		Find(&Materialss).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return Materialss, err
}

// Update 更新材料信息
func (md *MaterialDao) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return db.Client.WithContext(ctx).
		Model(&entity.Materials{}).
		Where(entity.MaterialsColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// DeleteByID 删除材料记录
func (md *MaterialDao) DeleteByID(ctx context.Context, id string) error {
	return db.Client.WithContext(ctx).
		Where(entity.MaterialsColumns.ID+" = ?", id).
		Delete(&entity.Materials{}).Error
}

// DeleteByGroupID 删除某个材料组内的所有材料
func (md *MaterialDao) DeleteByGroupID(ctx context.Context, groupID string) error {
	return db.Client.WithContext(ctx).
		Where(entity.MaterialsColumns.MaterialGroupID+" = ?", groupID).
		Delete(&entity.Materials{}).Error
}

// DeleteByGroupIdListTx 删除某个材料组内的所有材料
func (md *MaterialDao) DeleteByGroupIdListTx(ctx context.Context, tx *gorm.DB, groupIdList []string) error {
	return tx.WithContext(ctx).
		Where(entity.MaterialsColumns.MaterialGroupID+" in ?", groupIdList).
		Delete(&entity.Materials{}).Error
}

// Query 分页查询材料
func (md *MaterialDao) Query(ctx context.Context, page int, pageSize int,
	filters map[string]interface{}) ([]*entity.Materials, int64, error) {
	var (
		Materialss []*entity.Materials
		total      int64
	)

	query := db.Client.WithContext(ctx).Model(&entity.Materials{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			if utils.Contains([]string{entity.MaterialsColumns.MaterialName}, key) {
				query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
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
		Order(entity.MaterialsColumns.Sort + " DESC").
		Find(&Materialss).Error; err != nil {
		return nil, 0, err
	}

	return Materialss, total, nil
}
