package dao

import (
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/db"
	"ByteScience-WAM-Business/pkg/logger"
	"context"
	"errors"
	"time"

	"ByteScience-WAM-Business/internal/model/entity"
	"gorm.io/gorm"
)

// RecipeDao 数据访问对象，封装配方相关操作
type RecipeDao struct{}

// NewRecipeDao 创建一个新的 RecipeDao 实例
func NewRecipeDao() *RecipeDao {
	return &RecipeDao{}
}

// Insert 插入配方记录
func (ed *RecipeDao) Insert(ctx context.Context, recipe *entity.Recipes) error {
	return db.Client.WithContext(ctx).Create(recipe).Error
}

// Count 统计配方个数
func (ed *RecipeDao) Count(ctx context.Context) (int64, error) {
	var num int64
	err := db.Client.Model(&entity.Recipes{}).WithContext(ctx).Count(&num).Error

	return num, err
}

// GetByID 根据 ID 获取配方
func (ed *RecipeDao) GetByID(ctx context.Context, id string) (*entity.Recipes, error) {
	var recipe entity.Recipes
	err := db.Client.WithContext(ctx).
		Where(entity.RecipesColumns.ID+" = ?", id).
		First(&recipe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &recipe, err
}

// GetMaterialByIdList 根据配方ID列表获取材料信息
func (ed *RecipeDao) GetMaterialByIdList(ctx context.Context, recipeIdList []string) ([]entity.RecipeMaterialGroups,
	[]entity.Materials, error) {
	var recipeInIdList []string
	if err := db.Client.WithContext(ctx).Model(&entity.Recipes{}).Distinct(entity.RecipesColumns.ID).
		Where(entity.RecipesColumns.ID+" IN ?", recipeIdList).
		Pluck(entity.RecipesColumns.ID, &recipeInIdList).Error; err != nil {
		logger.Logger.Errorf("[dao.GetMaterialByIdList] get recipeInIdList Mysql err: %v", err)
		return nil, nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	if len(recipeInIdList) != len(recipeIdList) {
		return nil, nil, utils.NewBusinessError(utils.RecipeDoesNotExistCode)
	}

	recipeMaterialGroups := make([]entity.RecipeMaterialGroups, 0)
	if err := db.Client.WithContext(ctx).Where(entity.RecipeMaterialGroupsColumns.RecipeID+" IN ?",
		recipeIdList).Find(&recipeMaterialGroups).Error; err != nil {
		logger.Logger.Errorf("[dao.GetMaterialByIdList] get recipeMaterialGroups Mysql err: %v", err)
		return nil, nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	experimentMaterialGroupIdList := make([]string, 0)
	for _, recipeMaterialGroup := range recipeMaterialGroups {
		experimentMaterialGroupIdList = append(experimentMaterialGroupIdList, recipeMaterialGroup.ExperimentMaterialGroupID)
	}

	materials := make([]entity.Materials, 0)
	if err := db.Client.WithContext(ctx).Where(entity.MaterialsColumns.ExperimentMaterialGroupID+" IN ?",
		experimentMaterialGroupIdList).Find(&materials).Error; err != nil {
		logger.Logger.Errorf("[dao.GetMaterialByIdList] get materials Mysql err: %v", err)
		return nil, nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	return recipeMaterialGroups, materials, nil
}

// GetRecipeInSignatureMap 获取已存在的签名map
func (ed *RecipeDao) GetRecipeInSignatureMap(ctx context.Context,
	recipeSignatureList []string) (map[string]string, error) {
	var recipeList []entity.Recipes
	recipeMap := make(map[string]string)
	err := db.Client.WithContext(ctx).Model(&entity.Recipes{}).
		Where(entity.RecipesColumns.RecipeSignature+" in ?", recipeSignatureList).
		Find(&recipeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	for _, row := range recipeList {
		recipeMap[row.RecipeSignature] = row.ID
	}

	return recipeMap, err
}

// Update 更新配方信息
func (ed *RecipeDao) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return db.Client.WithContext(ctx).
		Model(&entity.Recipes{}).
		Where(entity.RecipesColumns.ID+" = ?", id).
		Updates(updates).
		Error
}

// DeleteByID 删除配方记录
func (ed *RecipeDao) DeleteByID(ctx context.Context, id string) error {
	return db.Client.WithContext(ctx).
		Where(entity.RecipesColumns.ID+" = ?", id).
		Delete(&entity.Recipes{}).Error
}

// DeleteByIDTx 删除配方记录(处理事务)
func (ed *RecipeDao) DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error {
	return tx.WithContext(ctx).
		Where(entity.RecipesColumns.ID+" = ?", id).
		Delete(&entity.Recipes{}).Error
}

// UpdateLastUpdatedTime 更新配方的最后更新时间
func (ed *RecipeDao) UpdateLastUpdatedTime(ctx context.Context, id string) error {
	return db.Client.WithContext(ctx).
		Model(&entity.Recipes{}).
		Where(entity.RecipesColumns.ID+" = ?", id).
		Update(entity.RecipesColumns.UpdatedAt, time.Now()).
		Error
}

// Query 分页查询管理员
func (ed *RecipeDao) Query(ctx context.Context, page int, pageSize int,
	filters map[string]interface{}) ([]*entity.Recipes, int64, error) {
	var (
		recipeList []*entity.Recipes
		total      int64
	)

	// 定义需要使用 LIKE 查询的字段
	likeFields := []string{
		entity.RecipesColumns.RecipeName,
	}

	query := db.Client.WithContext(ctx).Model(&entity.Recipes{})

	// 应用过滤条件
	for key, value := range filters {
		if value != nil && value != "" {
			if utils.Contains(likeFields, key) {
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
		Order(entity.RecipesColumns.Sort + " DESC").
		Find(&recipeList).Error; err != nil {
		return nil, 0, err
	}

	return recipeList, total, nil
}
