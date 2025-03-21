package dao

import (
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"
	"gorm.io/gorm"
)

// RecordDao 数据访问对象，封装实验相关操作
type RecordDao struct{}

// NewRecordDao 创建一个新的 RecordDao 实例
func NewRecordDao() *RecordDao {
	return &RecordDao{}
}

// Insert 插入实验记录
func (ed *RecordDao) Insert(ctx context.Context, tx *gorm.DB, experiment *entity.OperationRecord) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).Create(experiment).Error
}

// Update 更新实验信息
func (ed *RecordDao) Update(ctx context.Context, tx *gorm.DB, id string, updates map[string]interface{}) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.OperationRecord{}).
		Where(entity.OperationRecordColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// DeleteByID 删除实验记录
func (ed *RecordDao) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Where(entity.OperationRecordColumns.ID+" = ?", id).
		Delete(&entity.OperationRecord{}).Error
}

// Count 统计实验个数
func (ed *RecordDao) Count(ctx context.Context) (int64, error) {
	var num int64
	err := db.Client.Model(&entity.OperationRecord{}).WithContext(ctx).Count(&num).Error

	return num, err
}

// GetByID 根据 ID 获取实验
func (ed *RecordDao) GetByID(ctx context.Context, id string) (*entity.OperationRecord, error) {
	var operationRecord entity.OperationRecord
	err := db.Client.WithContext(ctx).
		Where(entity.OperationRecordColumns.ID+" = ?", id).
		First(&operationRecord).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &operationRecord, err
}

// GetByIDList 根据 ID 获取数据
func (ed *RecordDao) GetByIDList(ctx context.Context, idList []string) ([]entity.OperationRecord, error) {
	var operationRecordList []entity.OperationRecord
	err := db.Client.WithContext(ctx).
		Where(entity.OperationRecordColumns.ID+" in ?", idList).
		Find(&operationRecordList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return operationRecordList, err
}

// GetByFileID 根据 userId 获取数据
func (ed *RecordDao) GetByFileID(ctx context.Context, userId string) ([]*entity.OperationRecord, error) {
	var operationRecordList []*entity.OperationRecord
	err := db.Client.WithContext(ctx).
		Where(entity.OperationRecordColumns.UserID+" = ?", userId).
		Find(&operationRecordList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return operationRecordList, err
}

// Query 分页查询管理员
func (ed *RecordDao) Query(ctx context.Context, page int, pageSize int,
	filters map[string]interface{}) ([]*entity.OperationRecord, int64, error) {
	var (
		experiments []*entity.OperationRecord
		total       int64
	)

	query := db.Client.WithContext(ctx).Model(&entity.OperationRecord{})

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
		Order(entity.OperationRecordColumns.CreatedAt + " DESC").
		Find(&experiments).Error; err != nil {
		return nil, 0, err
	}

	return experiments, total, nil
}
