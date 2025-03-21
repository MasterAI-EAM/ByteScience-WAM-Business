package dao

import (
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"

	"ByteScience-WAM-Business/internal/model/entity"
	"gorm.io/gorm"
)

// ExperimentStepDao 数据访问对象，封装实验步骤相关操作
type ExperimentStepDao struct{}

// NewExperimentStepDao 创建一个新的 ExperimentStepDao 实例
func NewExperimentStepDao() *ExperimentStepDao {
	return &ExperimentStepDao{}
}

// Insert 插入实验步骤记录
func (esd *ExperimentStepDao) Insert(ctx context.Context, tx *gorm.DB, step *entity.ExperimentSteps) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).Create(step).Error
}

// UpdateResultValue 更新步骤的结果值
func (esd *ExperimentStepDao) UpdateResultValue(ctx context.Context, tx *gorm.DB, id string, resultValue string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.ExperimentSteps{}).
		Where(entity.ExperimentStepsColumns.ID+" = ?", id).
		Update(entity.ExperimentStepsColumns.ResultValue, resultValue).
		Error
}

// Update 更新实验步骤
func (esd *ExperimentStepDao) Update(ctx context.Context, tx *gorm.DB, id string, updates map[string]interface{}) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.ExperimentSteps{}).
		Where(entity.ExperimentStepsColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// DeleteByID 删除实验步骤记录
func (esd *ExperimentStepDao) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Where(entity.ExperimentStepsColumns.ID+" = ?", id).
		Delete(&entity.ExperimentSteps{}).Error
}

// DeleteByExperimentID 删除某实验的所有步骤
func (esd *ExperimentStepDao) DeleteByExperimentID(ctx context.Context, tx *gorm.DB, experimentID string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Where(entity.ExperimentStepsColumns.ExperimentID+" = ?", experimentID).
		Delete(&entity.ExperimentSteps{}).Error
}

// GetByID 根据 ID 获取实验步骤
func (esd *ExperimentStepDao) GetByID(ctx context.Context, id string) (*entity.ExperimentSteps, error) {
	var step entity.ExperimentSteps
	err := db.Client.WithContext(ctx).
		Where(entity.ExperimentStepsColumns.ID+" = ?", id).
		First(&step).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &step, err
}

// GetByExperimentID 根据实验 ID 获取所有步骤
func (esd *ExperimentStepDao) GetByExperimentID(ctx context.Context, experimentID string) ([]*entity.ExperimentSteps, error) {
	var steps []*entity.ExperimentSteps
	err := db.Client.WithContext(ctx).
		Where(entity.ExperimentStepsColumns.ExperimentID+" = ?", experimentID).
		Order(entity.ExperimentStepsColumns.Sort + " ASC").
		Find(&steps).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return steps, err
}

// Query 分页查询实验步骤
func (esd *ExperimentStepDao) Query(ctx context.Context, page int, pageSize int,
	filters map[string]interface{}) ([]*entity.ExperimentSteps, int64, error) {
	var (
		steps []*entity.ExperimentSteps
		total int64
	)

	query := db.Client.WithContext(ctx).Model(&entity.ExperimentSteps{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key+" = ?", value)
		}
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Scopes(db.PageScope(page, pageSize)).
		Order(entity.ExperimentStepsColumns.Sort + " ASC").
		Find(&steps).Error; err != nil {
		return nil, 0, err
	}

	return steps, total, nil
}
