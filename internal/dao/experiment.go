package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

// ExperimentDao 数据访问对象，封装实验相关操作
type ExperimentDao struct{}

// NewExperimentDao 创建一个新的 ExperimentDao 实例
func NewExperimentDao() *ExperimentDao {
	return &ExperimentDao{}
}

// Insert 插入实验记录
func (ed *ExperimentDao) Insert(ctx context.Context, tx *gorm.DB, experiment *entity.Experiment) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).Create(experiment).Error
}

// Update 更新实验信息
func (ed *ExperimentDao) Update(ctx context.Context, tx *gorm.DB, id string, updates map[string]interface{}) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Experiment{}).
		Where(entity.ExperimentColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// DeleteByID 删除实验记录
func (ed *ExperimentDao) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Where(entity.ExperimentColumns.ID+" = ?", id).
		Delete(&entity.Experiment{}).Error
}

// UpdateLastUpdatedTime 更新实验的最后更新时间
func (ed *ExperimentDao) UpdateLastUpdatedTime(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Experiment{}).
		Where(entity.ExperimentColumns.ID+" = ?", id).
		Update(entity.ExperimentColumns.UpdatedAt, time.Now()).
		Error
}

// Count 统计实验个数
func (ed *ExperimentDao) Count(ctx context.Context) (int64, error) {
	var num int64
	err := db.Client.Model(&entity.Experiment{}).WithContext(ctx).Count(&num).Error

	return num, err
}

// GetByID 根据 ID 获取实验
func (ed *ExperimentDao) GetByID(ctx context.Context, id string) (*entity.Experiment, error) {
	var experiment entity.Experiment
	err := db.Client.WithContext(ctx).
		Where(entity.ExperimentColumns.ID+" = ?", id).
		First(&experiment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &experiment, err
}

// GetByIDList 根据 ID 获取实验
func (ed *ExperimentDao) GetByIDList(ctx context.Context, idList []string) ([]entity.Experiment, error) {
	var experimentList []entity.Experiment
	err := db.Client.WithContext(ctx).
		Where(entity.ExperimentColumns.ID+" in ?", idList).
		Find(&experimentList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return experimentList, err
}

// GetByExperimentSignature 根据 experimentSignature 获取实验
func (ed *ExperimentDao) GetByExperimentSignature(ctx context.Context, experimentSignature string) ([]*entity.Experiment, error) {
	var experiments []*entity.Experiment
	err := db.Client.WithContext(ctx).
		Where(entity.ExperimentColumns.Signature+" = ?", experimentSignature).
		Find(&experiments).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return experiments, err
}

// GetExperimentInSignatureList 获取已存在的签名列表
func (ed *ExperimentDao) GetExperimentInSignatureList(ctx context.Context,
	experimentSignatureList []string) ([]string, error) {
	var experimentInSignatureList []string
	err := db.Client.WithContext(ctx).Model(&entity.Experiment{}).
		Where(entity.ExperimentColumns.Signature+" in ?", experimentSignatureList).
		Pluck(entity.ExperimentColumns.Signature, &experimentInSignatureList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return experimentInSignatureList, err
}

// Query 分页查询管理员
func (ed *ExperimentDao) Query(ctx context.Context, page int, pageSize int,
	filters map[string]interface{}) ([]*entity.Experiment, int64, error) {
	var (
		experiments []*entity.Experiment
		total       int64
	)

	// 定义需要使用 LIKE 查询的字段
	likeFields := []string{
		entity.ExperimentColumns.ExperimentName,
		entity.ExperimentColumns.Experimenter,
	}

	query := db.Client.WithContext(ctx).Model(&entity.Experiment{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			if utils.Contains(likeFields, key) {
				query = query.Where(key+" LIKE ?", value.(string)+"%")
			} else if key == "startTime" {
				query = query.Where(entity.ExperimentColumns.CreatedAt+" >= ?", value.(string))
			} else if key == "endTime" {
				query = query.Where(entity.ExperimentColumns.CreatedAt+" <= ?", value.(string))
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
		Order(entity.ExperimentColumns.Sort + " DESC").
		Find(&experiments).Error; err != nil {
		return nil, 0, err
	}

	return experiments, total, nil
}
