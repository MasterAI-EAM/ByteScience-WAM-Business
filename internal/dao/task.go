package dao

import (
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"
	"gorm.io/gorm/logger"
	"time"

	"ByteScience-WAM-Business/internal/model/entity"
	"gorm.io/gorm"
)

// TaskDao 数据访问对象
type TaskDao struct{}

// NewTaskDao 创建一个新的 TaskDao 实例
func NewTaskDao() *TaskDao {
	return &TaskDao{}
}

const (
	insertNum            = 500
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusSuccess    = "success"
	TaskStatusFailure    = "failure"
)

// Insert 插入记录
func (ed *TaskDao) Insert(ctx context.Context, tx *gorm.DB, task *entity.Task) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).Create(task).Error
}

// BatchInsert 插入记录
func (ed *TaskDao) BatchInsert(ctx context.Context, tx *gorm.DB, tasks []*entity.Task) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).CreateInBatches(tasks, insertNum).Error
}

// Count 统计个数
func (ed *TaskDao) Count(ctx context.Context) (int64, error) {
	var num int64
	err := db.Client.Model(&entity.Task{}).WithContext(ctx).Count(&num).Error

	return num, err
}

// GetByID 根据 ID 获取数据
func (ed *TaskDao) GetByID(ctx context.Context, id string) (*entity.Task, error) {
	var info entity.Task
	err := db.Client.WithContext(ctx).
		Where(entity.TaskColumns.ID+" = ?", id).
		First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &info, err
}

// GetByIDList 根据 ID 获取实验
func (ed *TaskDao) GetByIDList(ctx context.Context, idList []string) ([]entity.Task, error) {
	var list []entity.Task
	err := db.Client.WithContext(ctx).
		Where(entity.TaskColumns.ID+" in ?", idList).
		Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return list, err
}

// GetPendingData 获取一条最旧的未审批数据（不记录日志）
func (ed *TaskDao) GetPendingData(ctx context.Context) (*entity.Task, error) {
	var info entity.Task
	err := db.Client.WithContext(ctx).
		Session(&gorm.Session{Logger: logger.Discard}). // 关闭日志
		Where(entity.TaskColumns.Status, TaskStatusPending).
		Order(entity.TaskColumns.CreatedAt).
		Limit(1).
		First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &info, err
}

// Update 更新实验信息
func (ed *TaskDao) Update(ctx context.Context, id string, updates map[string]interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Task{}).
		Where(entity.TaskColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// DeleteByID 删除实验记录
func (ed *TaskDao) DeleteByID(ctx context.Context, id string, tx *gorm.DB) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Where(entity.TaskColumns.ID+" = ?", id).
		Delete(&entity.Task{}).Error
}

// UpdateLastUpdatedTime 更新实验的最后更新时间
func (ed *TaskDao) UpdateLastUpdatedTime(ctx context.Context, id string, tx *gorm.DB) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Task{}).
		Where(entity.TaskColumns.ID+" = ?", id).
		Update(entity.TaskColumns.UpdatedAt, time.Now()).
		Error
}

// Query 分页查询管理员
func (ed *TaskDao) Query(ctx context.Context, page int, pageSize int,
	filters map[string]interface{}) ([]*entity.Task, int64, error) {
	var (
		experiments []*entity.Task
		total       int64
	)

	// 定义需要使用 LIKE 查询的字段
	likeFields := []string{
		entity.TaskColumns.FileName,
	}

	query := db.Client.WithContext(ctx).Model(&entity.Task{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			if utils.Contains(likeFields, key) {
				query = query.Where(key+" LIKE ?", value.(string)+"%")
			} else if key == "startTime" {
				query = query.Where(entity.TaskColumns.CreatedAt+" >= ?", value.(string))
			} else if key == "endTime" {
				query = query.Where(entity.TaskColumns.CreatedAt+" <= ?", value.(string))
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
		Order(entity.TaskColumns.CreatedAt + " DESC").
		Find(&experiments).Error; err != nil {
		return nil, 0, err
	}

	return experiments, total, nil
}
