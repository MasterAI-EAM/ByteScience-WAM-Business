package dao

import (
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/db"
	"context"
	"errors"
	"strings"
	"time"

	"ByteScience-WAM-Business/internal/model/entity"
	"gorm.io/gorm"
)

// UserDao 数据访问对象，封装角色相关操作
type UserDao struct{}

// NewUserDao 创建一个新的 UserDao 实例
func NewUserDao() *UserDao {
	return &UserDao{}
}

// Insert 插入用户记录
func (ud *UserDao) Insert(ctx context.Context, tx *gorm.DB, user *entity.Users) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).Create(user).Error
}

// UpdateStatus 更新用户状态
func (ud *UserDao) UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status int) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Users{}).
		Where(entity.UsersColumns.ID+" = ?", id).
		Update(entity.UsersColumns.Status, status).
		Error
}

// Update 更新用户信息
func (ud *UserDao) Update(ctx context.Context, tx *gorm.DB, id string, updates map[string]interface{}) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Users{}).
		Where(entity.UsersColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// SoftDeleteByID 软删除用户记录
func (ud *UserDao) SoftDeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Users{}).
		Where(entity.UsersColumns.ID+" = ?", id).
		Update(entity.UsersColumns.DeletedAt, time.Now()).
		Error
}

// UpdateLastLoginTime 更新管理员的最后登录时间
func (ud *UserDao) UpdateLastLoginTime(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = db.Client
	}
	return tx.WithContext(ctx).
		Model(&entity.Admins{}).
		Where(entity.AdminsColumns.ID+" = ?", id).
		Update(entity.AdminsColumns.LastLoginAt, time.Now()).
		Error
}

// GetByID 根据 ID 获取用户
func (ud *UserDao) GetByID(ctx context.Context, id string) (*entity.Users, error) {
	var user entity.Users
	err := db.Client.WithContext(ctx).
		Where(entity.UsersColumns.ID+" = ?", id).
		Where(entity.UsersColumns.DeletedAt + " IS NULL").
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByIDs 根据 ID 获取用户
func (ud *UserDao) GetByIDs(ctx context.Context, ids []string) ([]entity.Users, error) {
	var users []entity.Users
	err := db.Client.WithContext(ctx).
		Where(entity.UsersColumns.ID+" in ?", ids).
		Where(entity.UsersColumns.DeletedAt + " IS NULL").
		First(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return users, err
}

// GetByFields 根据字段（用户名、邮箱、手机号）获取用户
func (ud *UserDao) GetByFields(ctx context.Context, username, email, phone string) (*entity.Users, error) {
	// 构建查询条件
	query := db.Client.WithContext(ctx).Model(&entity.Users{}).
		Where(entity.UsersColumns.DeletedAt + " IS NULL")

	conditions := []string{}
	params := []interface{}{}

	if username != "" {
		conditions = append(conditions, entity.UsersColumns.Username+" = ?")
		params = append(params, username)
	}
	if email != "" {
		conditions = append(conditions, entity.UsersColumns.Email+" = ?")
		params = append(params, email)
	}
	if phone != "" {
		conditions = append(conditions, entity.UsersColumns.Phone+" = ?")
		params = append(params, phone)
	}

	if len(conditions) > 0 {
		query = query.Where("("+strings.Join(conditions, " OR ")+")", params...)
	}

	var user entity.Users
	err := query.First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Query 分页查询用户
func (ud *UserDao) Query(ctx context.Context, page int, pageSize int,
	filters map[string]interface{}) ([]*entity.Users, int64, error) {
	var (
		users []*entity.Users
		total int64
	)

	// 定义需要使用 LIKE 查询的字段
	likeFields := []string{
		entity.UsersColumns.Username,
		entity.UsersColumns.Email,
		entity.UsersColumns.Phone,
		entity.UsersColumns.Nickname,
	}

	query := db.Client.WithContext(ctx).Model(&entity.Users{}).Where(entity.UsersColumns.DeletedAt + " IS NULL")

	for key, value := range filters {
		if value != nil && value != "" {
			if utils.Contains(likeFields, key) {
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
		Order(entity.UsersColumns.CreatedAt + " DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
