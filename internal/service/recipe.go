package service

import (
	"ByteScience-WAM-Business/internal/dao"
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/data"
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/db"
	"ByteScience-WAM-Business/pkg/logger"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

type RecipeService struct {
	recipeDao *dao.RecipeDao
}

// NewRecipeService 创建一个新的 RecipeService 实例
func NewRecipeService() *RecipeService {
	return &RecipeService{
		recipeDao: dao.NewRecipeDao(),
	}
}

// List 查询配方及其相关数据，包括材料组及材料信息
func (rs *RecipeService) List(ctx context.Context, req *data.RecipeListRequest) (*data.RecipeListResponse, error) {
	// 查询配方及其基础信息
	recipeList, total, err := rs.recipeDao.Query(ctx, req.Page, req.PageSize, map[string]interface{}{
		entity.RecipesColumns.RecipeName: req.RecipeName,
	})
	if err != nil {
		logger.Logger.Errorf("[RecipeService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	if len(recipeList) == 0 {
		return &data.RecipeListResponse{Total: 0, List: []data.RecipeData{}}, nil
	}

	// 获取所有配方ID
	recipeIDs := make([]string, len(recipeList))
	recipeMap := make(map[string]*data.RecipeData) // 用于存储配方数据（无序）
	for i, recipe := range recipeList {
		recipeIDs[i] = recipe.ID
		recipeMap[recipe.ID] = &data.RecipeData{
			RecipeId:       recipe.ID,
			RecipeName:     recipe.RecipeName,
			CreatedAt:      recipe.CreatedAt.Format(time.RFC3339),
			MaterialGroups: []data.MaterialGroupInfo{},
		}
	}

	// 查询所有材料组
	var materialGroups []entity.RecipeMaterialGroups
	if err = db.Client.WithContext(ctx).
		Where(entity.RecipeMaterialGroupsColumns.RecipeID+" IN (?)", recipeIDs).
		Find(&materialGroups).Error; err != nil {
		logger.Logger.Errorf("[RecipeService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	// 组织材料组数据
	materialGroupIDs := make([]string, len(materialGroups))
	materialGroupMap := make(map[string]*data.MaterialGroupInfo)
	for i, group := range materialGroups {
		materialGroupIDs[i] = group.ExperimentMaterialGroupID
		materialGroup := &data.MaterialGroupInfo{
			MaterialGroupID: group.ExperimentMaterialGroupID,
			Proportion:      group.Proportion,
			Materials:       []data.MaterialInfo{},
		}
		materialGroupMap[group.ExperimentMaterialGroupID] = materialGroup
	}

	// 查询所有材料
	var materials []entity.Materials
	if err = db.Client.WithContext(ctx).
		Where(entity.MaterialsColumns.ExperimentMaterialGroupID+" IN (?)", materialGroupIDs).
		Find(&materials).Error; err != nil {
		logger.Logger.Errorf("[RecipeService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	// 填充材料数据到材料组中
	for _, material := range materials {
		if materialGroup, exists := materialGroupMap[material.ExperimentMaterialGroupID]; exists {
			materialGroup.Materials = append(materialGroup.Materials, data.MaterialInfo{
				MaterialID: material.ID,
				MaterialData: data.MaterialData{
					MaterialName: material.MaterialName,
					Percentage:   material.Percentage,
				},
			})
			materialGroup.MaterialGroupName = material.MaterialGroupName
		}
	}

	// 将材料组添加到对应的配方中
	for _, group := range materialGroups {
		if recipe, exists := recipeMap[group.RecipeID]; exists {
			if materialGroup, exists := materialGroupMap[group.ExperimentMaterialGroupID]; exists {
				recipe.MaterialGroups = append(recipe.MaterialGroups, *materialGroup)
			}
		}
	}

	// **按 recipeList 的顺序填充 recipeListResponse**
	recipeListResponse := make([]data.RecipeData, 0, len(recipeList))
	for _, recipe := range recipeList {
		if recipeData, exists := recipeMap[recipe.ID]; exists {
			recipeListResponse = append(recipeListResponse, *recipeData)
		}
	}

	return &data.RecipeListResponse{
		Total: total,
		List:  recipeListResponse,
	}, nil
}

// Add 添加配方数据
func (rs *RecipeService) Add(ctx context.Context, req *data.RecipeAddRequest) (*dto.Empty, error) {
	recipe, recipeMaterialGroups, materials := getInsertRecipeData(req)

	var foundRecipes []entity.Recipes
	if err := db.Client.WithContext(ctx).Model(entity.Recipes{}).
		Where(entity.RecipesColumns.RecipeSignature, recipe.RecipeSignature).Find(&foundRecipes).Error; err != nil {
		logger.Logger.Errorf("[RecipeService Add] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	if len(foundRecipes) > 0 {
		return nil, utils.NewBusinessError(utils.DuplicateRecipeFormatCode, foundRecipes[0].RecipeName)
	}

	// 使用事务闭包
	if err := db.Client.Transaction(func(tx *gorm.DB) error {
		maxNum := 500
		if err := tx.Create(&recipe).Error; err != nil {
			logger.Logger.Errorf("[RecipeService Add] Create recipe err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(recipeMaterialGroups, maxNum).Error; err != nil {
			logger.Logger.Errorf("[RecipeService Add] CreateInBatches recipeMaterialGroups err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(materials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[RecipeService Add] CreateInBatches materials err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		logger.Logger.Errorf("[RecipeService Add] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}
	return nil, nil
}

// Info 查询单个配方及其相关数据，包括材料组及材料信息
func (rs *RecipeService) Info(ctx context.Context, req *data.RecipeInfoRequest) (*data.RecipeInfoResponse, error) {
	// 查询单个配方及其基础信息
	recipe, err := rs.recipeDao.GetByID(ctx, req.RecipeId)
	if err != nil {
		logger.Logger.Errorf("[RecipeService Info] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}
	if recipe == nil {
		return nil, nil
	}

	recipeData := &data.RecipeInfoResponse{
		RecipeId:              recipe.ID,
		RecipeName:            recipe.RecipeName,
		CreatedAt:             recipe.CreatedAt.Format(time.RFC3339),
		MaterialGroups:        make([]data.MaterialGroupInfo, 0),
		RecipeBasedExperiment: make([]data.ExperimentInfo, 0),
	}

	var experimentIdList []string
	if err = db.Client.WithContext(ctx).Model(&entity.ExperimentSteps{}).Distinct(entity.ExperimentStepsColumns.
		ExperimentID).Where(entity.ExperimentStepsColumns.RecipeID+" = ?", req.RecipeId).
		Pluck(entity.ExperimentStepsColumns.ExperimentID, &experimentIdList).Error; err != nil {
		logger.Logger.Errorf("[RecipeService Info] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	if err = db.Client.WithContext(ctx).Model(&entity.Experiment{}).
		Select(entity.ExperimentColumns.ID+" as id", entity.ExperimentColumns.ExperimentName+" as name").
		Where(entity.ExperimentColumns.ID+" IN ?", experimentIdList).
		Find(&recipeData.RecipeBasedExperiment).Error; err != nil {
		logger.Logger.Errorf("[RecipeService Info] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	// 查询该配方的所有材料组
	var materialGroups []entity.RecipeMaterialGroups
	if err = db.Client.WithContext(ctx).
		Where(entity.RecipeMaterialGroupsColumns.RecipeID+" = ?", req.RecipeId).
		Find(&materialGroups).Error; err != nil {
		logger.Logger.Errorf("[RecipeService Info] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	// 组织材料组数据
	materialGroupIDs := make([]string, len(materialGroups))
	materialGroupMap := make(map[string]*data.MaterialGroupInfo)
	for i, group := range materialGroups {
		materialGroupIDs[i] = group.ExperimentMaterialGroupID
		materialGroup := &data.MaterialGroupInfo{
			MaterialGroupID: group.ExperimentMaterialGroupID,
			Proportion:      group.Proportion,
			Materials:       []data.MaterialInfo{},
		}
		materialGroupMap[group.ExperimentMaterialGroupID] = materialGroup
	}

	// 查询所有材料
	var materials []entity.Materials
	if err = db.Client.WithContext(ctx).
		Where(entity.MaterialsColumns.ExperimentMaterialGroupID+" IN (?)", materialGroupIDs).
		Find(&materials).Error; err != nil {
		logger.Logger.Errorf("[RecipeService Info] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	// 填充材料数据到材料组中
	for _, material := range materials {
		if materialGroup, exists := materialGroupMap[material.ExperimentMaterialGroupID]; exists {
			materialGroup.Materials = append(materialGroup.Materials, data.MaterialInfo{
				MaterialID: material.ID,
				MaterialData: data.MaterialData{
					MaterialName: material.MaterialName,
					Percentage:   material.Percentage,
				},
			})
			materialGroup.MaterialGroupName = material.MaterialGroupName
		}
	}

	// 将材料组添加到对应的配方中
	for _, group := range materialGroups {
		if materialGroup, exists := materialGroupMap[group.ExperimentMaterialGroupID]; exists {
			recipeData.MaterialGroups = append(recipeData.MaterialGroups, *materialGroup)
		}
	}

	return recipeData, nil
}

func (rs *RecipeService) FormList(ctx context.Context, req *data.RecipeFormListRequest) (*data.RecipeFormListResponse, error) {
	// 查询配方及其基础信息
	recipeList, total, err := rs.recipeDao.Query(ctx, req.Page, req.PageSize, map[string]interface{}{
		entity.RecipesColumns.RecipeName: req.RecipeName,
	})
	if err != nil {
		logger.Logger.Errorf("[RecipeService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode)
	}

	res := &data.RecipeFormListResponse{List: make([]data.RecipeInfo, 0), Total: total}
	for _, recipe := range recipeList {
		res.List = append(res.List, data.RecipeInfo{
			Id:   recipe.ID,
			Name: recipe.RecipeName,
		})
	}

	return res, err
}

// getInsertRecipeData 将插入数据内容转换成入库数据
func getInsertRecipeData(data *data.RecipeAddRequest) (entity.Recipes, []entity.RecipeMaterialGroups, []entity.Materials) {
	var materials []entity.Materials
	var recipeMaterialGroups []entity.RecipeMaterialGroups
	recipe := entity.Recipes{
		ID:         uuid.NewString(),
		RecipeName: data.RecipeName,
		Sort:       data.Sort,
	}
	for _, materialGroup := range data.MaterialGroups {
		recipeMaterialGroup := entity.RecipeMaterialGroups{
			RecipeID:                  recipe.ID,
			ExperimentMaterialGroupID: uuid.NewString(),
			Proportion:                materialGroup.Proportion,
		}
		recipeMaterialGroups = append(recipeMaterialGroups, recipeMaterialGroup)
		for _, material := range materialGroup.Materials {
			materials = append(materials, entity.Materials{
				ID:                        uuid.NewString(),
				MaterialName:              material.MaterialName,
				ExperimentMaterialGroupID: recipeMaterialGroup.ExperimentMaterialGroupID,
				MaterialGroupName:         materialGroup.MaterialGroupName,
				Percentage:                material.Percentage,
			})
		}
	}

	// 生成配方签名
	recipe.RecipeSignature = GenerateRecipeSignature(recipeMaterialGroups, materials)

	return recipe, recipeMaterialGroups, materials
}

// GenerateRecipeSignature 生成配方唯一密钥(传入的数据只有一个配方的时候用)
func GenerateRecipeSignature(recipeMaterialGroups []entity.RecipeMaterialGroups, materials []entity.Materials) string {
	// 创建一个映射，用于快速查找每个材料组的比例
	groupProportionMap := make(map[string]float64)
	for _, group := range recipeMaterialGroups {
		groupProportionMap[group.ExperimentMaterialGroupID] = group.Proportion
	}

	// 拉平所有材料，并补充比例信息
	allMaterials := make([]struct {
		MaterialGroupName string
		Proportion        float64
		MaterialName      string
		Percentage        float64
	}, 0, len(materials))
	for _, material := range materials {
		proportion, exists := groupProportionMap[material.ExperimentMaterialGroupID]
		if exists {
			allMaterials = append(allMaterials, struct {
				MaterialGroupName string
				Proportion        float64
				MaterialName      string
				Percentage        float64
			}{
				MaterialGroupName: material.MaterialGroupName,
				Proportion:        proportion,
				MaterialName:      material.MaterialName,
				Percentage:        material.Percentage,
			})
		}
	}

	// 按照材料组名称和材料名称排序
	sort.SliceStable(allMaterials, func(i, j int) bool {
		materialGroupNameI := strings.ToLower(strings.TrimSpace(allMaterials[i].MaterialGroupName))
		materialGroupNameJ := strings.ToLower(strings.TrimSpace(allMaterials[j].MaterialGroupName))

		materialNameI := strings.ToLower(strings.TrimSpace(allMaterials[i].MaterialName))
		materialNameJ := strings.ToLower(strings.TrimSpace(allMaterials[j].MaterialName))

		if materialGroupNameI == materialGroupNameJ {
			return materialNameI < materialNameJ
		}
		return materialGroupNameI < materialGroupNameJ
	})

	// 拼接所有材料信息
	var sb []string
	for _, material := range allMaterials {
		sb = append(sb, fmt.Sprintf("%s,%.2f,%s,%.2f", material.MaterialGroupName,
			material.Proportion,
			material.MaterialName, material.Percentage))
	}

	// 将拼接的字符串使用 "," 连接
	joined := strings.Join(sb, ",")

	// 使用 MD5 生成最终的签名
	hash := md5.New()
	hash.Write([]byte(joined))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateBatchRecipeSignature 批量生成配方唯一密钥
func GenerateBatchRecipeSignature(recipeMaterialGroups []entity.RecipeMaterialGroups, materials []entity.Materials,
) (map[string]string, []string) {
	// 创建映射：recipe_id -> 材料组比例
	recipeGroupMap := make(map[string]map[string]float64)
	for _, group := range recipeMaterialGroups {
		if _, exists := recipeGroupMap[group.RecipeID]; !exists {
			recipeGroupMap[group.RecipeID] = make(map[string]float64)
		}
		recipeGroupMap[group.RecipeID][group.ExperimentMaterialGroupID] = group.Proportion
	}

	// 创建映射：recipe_id -> 所有材料信息
	recipeMaterialsMap := make(map[string][]struct {
		MaterialGroupName string
		Proportion        float64
		MaterialName      string
		Percentage        float64
	})

	for _, material := range materials {
		for recipeID, groupProportions := range recipeGroupMap {
			if proportion, exists := groupProportions[material.ExperimentMaterialGroupID]; exists {
				recipeMaterialsMap[recipeID] = append(recipeMaterialsMap[recipeID], struct {
					MaterialGroupName string
					Proportion        float64
					MaterialName      string
					Percentage        float64
				}{
					MaterialGroupName: material.MaterialGroupName,
					Proportion:        proportion,
					MaterialName:      material.MaterialName,
					Percentage:        material.Percentage,
				})
			}
		}
	}

	// 生成签名
	recipeSignatureMap := make(map[string]string)
	recipeSignatureList := make([]string, len(recipeMaterialsMap))
	for recipeID, materials := range recipeMaterialsMap {
		// 排序
		sort.SliceStable(materials, func(i, j int) bool {
			groupNameI := strings.ToLower(strings.TrimSpace(materials[i].MaterialGroupName))
			groupNameJ := strings.ToLower(strings.TrimSpace(materials[j].MaterialGroupName))
			materialNameI := strings.ToLower(strings.TrimSpace(materials[i].MaterialName))
			materialNameJ := strings.ToLower(strings.TrimSpace(materials[j].MaterialName))

			if groupNameI == groupNameJ {
				return materialNameI < materialNameJ
			}
			return groupNameI < groupNameJ
		})

		// 拼接字符串
		var sb []string
		for _, material := range materials {
			sb = append(sb, fmt.Sprintf("%s,%.2f,%s,%.2f",
				material.MaterialGroupName,
				material.Proportion,
				material.MaterialName,
				material.Percentage))
		}

		// 生成 MD5 签名
		joined := strings.Join(sb, ",")
		hash := md5.New()
		hash.Write([]byte(joined))
		signature := hex.EncodeToString(hash.Sum(nil))

		recipeSignatureMap[recipeID] = signature
		recipeSignatureList = append(recipeSignatureList, signature)
	}

	return recipeSignatureMap, recipeSignatureList
}
