package service

import (
	"ByteScience-WAM-Business/conf"
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
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ExperimentService struct {
	experimentDao     *dao.ExperimentDao
	experimentStepDao *dao.ExperimentStepDao
	materialDao       *dao.MaterialDao
	recipeDao         *dao.RecipeDao
}

// NewExperimentService 创建一个新的 ExperimentService 实例
func NewExperimentService() *ExperimentService {
	return &ExperimentService{
		experimentDao:     dao.NewExperimentDao(),
		experimentStepDao: dao.NewExperimentStepDao(),
		materialDao:       dao.NewMaterialDao(),
		recipeDao:         dao.NewRecipeDao(),
	}
}

// List 查询实验及其相关数据，包括实验步骤和材料组
func (ss *ExperimentService) List(ctx context.Context, req *data.ExperimentListRequest) (*data.ExperimentListResponse, error) {
	// 查询实验及其基础信息
	experiments, total, err := ss.experimentDao.Query(ctx, req.Page, req.PageSize, map[string]interface{}{
		entity.ExperimentColumns.ExperimentName: req.ExperimentName,
		entity.ExperimentColumns.Experimenter:   req.Experimenter,
	})
	if err != nil {
		logger.Logger.Errorf("[ExperimentService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 获取所有实验ID
	experimentIDs := make([]string, len(experiments))
	fileIDs := make([]string, 0)
	for i, experiment := range experiments {
		experimentIDs[i] = experiment.ID
		fileIDs = append(fileIDs, experiment.FileID)
	}

	var files []entity.ExperimentFiles
	fileMap := make(map[string]string, 0)
	fileIDs = utils.RemoveDuplicates(fileIDs)
	if err = db.Client.WithContext(ctx).Where("id IN (?)", fileIDs).Find(&files).Error; err != nil {
		logger.Logger.Errorf("[ExperimentService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	for _, file := range files {
		fileMap[file.ID] = file.FileName
	}

	// 一次性查询所有实验步骤
	var steps []entity.ExperimentSteps
	if err = db.Client.WithContext(ctx).Where("experiment_id IN (?)", experimentIDs).
		Order("step_order DESC").
		Find(&steps).Error; err != nil {
		logger.Logger.Errorf("[ExperimentService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 获取所有步骤ID
	stepIDs := make([]string, len(steps))
	for i, step := range steps {
		stepIDs[i] = step.ID
	}

	// 构建查询结果映射
	stepMap := make(map[string][]data.ExperimentStepData)

	// 填充 stepMap
	for _, step := range steps {
		stepMap[step.ExperimentID] = append(stepMap[step.ExperimentID], data.ExperimentStepData{
			StepID:              step.ID,
			StepName:            step.StepName,
			StepNameDescription: conf.StepNameData[step.StepName],
			ExperimentCondition: step.ExperimentCondition,
			ResultValue:         step.ResultValue,
			StepOrder:           step.StepOrder,
			RecipeID:            step.RecipeID,
		})
	}

	// 组装最终的实验数据
	var experimentDataList []data.ExperimentData
	for _, experiment := range experiments {
		stepsData := stepMap[experiment.ID]
		experimentDataList = append(experimentDataList, data.ExperimentData{
			ExperimentID:   experiment.ID,
			ExperimentName: experiment.ExperimentName,
			FileID:         experiment.FileID,
			Experimenter:   experiment.Experimenter,
			UserID:         experiment.UserID,
			EntryCategory:  experiment.EntryCategory,
			FileName:       fileMap[experiment.FileID],
			StartTime:      utils.FormatTime(experiment.StartTime),
			EndTime:        utils.FormatTime(experiment.EndTime),
			CreatedAt:      experiment.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Steps:          stepsData,
		})
	}

	// 返回最终结果
	return &data.ExperimentListResponse{
		Total: total,
		List:  experimentDataList,
	}, nil
}

// Delete 删除实验数据
func (ss *ExperimentService) Delete(ctx context.Context, req *data.ExperimentDeleteRequest) (*dto.Empty, error) {
	experiment, err := ss.experimentDao.GetByID(ctx, req.ExperimentID)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Delete] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	if experiment == nil || experiment.ID == "" {
		return nil, utils.NewBusinessError(utils.ExperimentDoesNotExistCode, "")
	}

	if err = db.Client.Transaction(func(tx *gorm.DB) error {
		experimentStepIdList := make([]string, 0)
		if err = tx.WithContext(ctx).Where(entity.ExperimentStepsColumns.ExperimentID+" = ?", req.ExperimentID).
			Pluck(entity.ExperimentStepsColumns.ID, &experimentStepIdList).Error; err != nil {
			return err
		}

		if err := ss.experimentStepDao.DeleteByExperimentIDTx(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] Delete experimentStep err: %v", err)
			return err
		}

		if err := ss.experimentDao.DeleteByIDTx(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] Delete experiment err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		logger.Logger.Errorf("[ExperimentService Delete] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	return nil, err
}

// Add 添加实验数据
func (ss *ExperimentService) Add(ctx context.Context, userId string, req *data.ExperimentAddRequest) (*dto.Empty, error) {
	experiment := entity.Experiment{
		ID:             uuid.NewString(),
		FileID:         "",
		ExperimentName: req.ExperimentName,
		EntryCategory:  2,
		Sort:           req.Sort,
		Experimenter:   req.Experimenter,
		UserID:         userId,
		StartTime:      nil,
		EndTime:        nil,
	}
	layout := "2006-01-02T15:04:05Z"
	if req.StartTime != "" {
		startTime, _ := time.Parse(layout, req.StartTime)
		experiment.StartTime = &startTime
	}
	if req.EndTime != "" {
		endTime, _ := time.Parse(layout, req.EndTime)
		experiment.EndTime = &endTime
	}

	experimentSteps := make([]entity.ExperimentSteps, 0)
	recipeIdList := make([]string, 0)
	for _, step := range req.Steps {
		if !utils.Contains(recipeIdList, step.RecipeID) {
			recipeIdList = append(recipeIdList, step.RecipeID)
		}
		experimentStep := entity.ExperimentSteps{
			ID:                  uuid.NewString(),
			ExperimentID:        experiment.ID,
			RecipeID:            step.RecipeID,
			StepOrder:           step.StepOrder,
			StepName:            step.StepName,
			ResultValue:         step.ResultValue,
			ExperimentCondition: step.ExperimentCondition,
		}
		experimentSteps = append(experimentSteps, experimentStep)
	}

	// 获取配方数据
	recipeMaterialGroups, materials, err := ss.recipeDao.GetMaterialByIdList(ctx, recipeIdList)
	if err != nil {
		return nil, err
	}

	// 生成试验签名
	experimentSignatureMap, _ := generateExperimentSignature([]entity.Experiment{experiment},
		experimentSteps, recipeMaterialGroups, materials)
	experiment.ExperimentSignature = experimentSignatureMap[experiment.ID]

	experimentList, err := ss.experimentDao.GetByExperimentSignature(ctx, experiment.ExperimentSignature)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Add] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	if len(experimentList) > 0 {
		return nil, utils.NewBusinessError(utils.DuplicateExperimentFormatCode, "", experimentList[0].ExperimentName)
	}

	// 使用事务闭包
	if err := db.Client.Transaction(func(tx *gorm.DB) error {
		maxNum := 500
		if err := tx.Create(&experiment).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Add] CreateInBatches experiment err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experimentSteps, maxNum).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Add] CreateInBatches experimentSteps err: %v", err)
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Errorf("[ExperimentService Add] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	return nil, nil
}

// Edit 修改实验数据
func (ss *ExperimentService) Edit(ctx context.Context, req *data.ExperimentUpdateRequest) (*dto.Empty, error) {
	experiment, err := ss.experimentDao.GetByID(ctx, req.ExperimentID)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Edit] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	if experiment == nil || experiment.ID == "" {
		return nil, utils.NewBusinessError(utils.ExperimentDoesNotExistCode, "")
	}

	experimentSteps := make([]entity.ExperimentSteps, 0)
	recipeIdList := make([]string, 0)
	for _, step := range req.Steps {
		if !utils.Contains(recipeIdList, step.RecipeID) {
			recipeIdList = append(recipeIdList, step.RecipeID)
		}
		experimentStep := entity.ExperimentSteps{
			ID:                  uuid.NewString(),
			ExperimentID:        experiment.ID,
			RecipeID:            step.RecipeID,
			StepOrder:           step.StepOrder,
			StepName:            step.StepName,
			ResultValue:         step.ResultValue,
			ExperimentCondition: step.ExperimentCondition,
		}
		experimentSteps = append(experimentSteps, experimentStep)
	}

	// 获取配方数据
	recipeMaterialGroups, materials, err := ss.recipeDao.GetMaterialByIdList(ctx, recipeIdList)
	if err != nil {
		return nil, err
	}

	// 生成试验签名
	experimentSignatureMap, _ := generateExperimentSignature([]entity.Experiment{*experiment},
		experimentSteps, recipeMaterialGroups, materials)
	experiment.ExperimentSignature = experimentSignatureMap[experiment.ID]

	experimentList, err := ss.experimentDao.GetByExperimentSignature(ctx, experiment.ExperimentSignature)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Add] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 标记是否为自身重复（允许修改的情况）
	isSelfDuplicate := len(experimentList) == 1 && experimentList[0].ID == req.ExperimentID

	// 检查查询到的实验列表情况
	if len(experimentList) > 0 && !isSelfDuplicate {
		// 存在重复实验（非自身重复），返回业务错误
		return nil, utils.NewBusinessError(utils.DuplicateExperimentFormatCode, "", experimentList[0].ExperimentName)
	}

	if err = db.Client.Transaction(func(tx *gorm.DB) error {
		// 修改实验
		update := map[string]interface{}{
			entity.ExperimentColumns.ExperimentName:      req.ExperimentName,
			entity.ExperimentColumns.Experimenter:        req.Experimenter,
			entity.ExperimentColumns.ExperimentSignature: experimentSignatureMap[experiment.ID],
			entity.ExperimentColumns.StartTime:           nil,
			entity.ExperimentColumns.EndTime:             nil,
		}
		if req.StartTime != "" {
			update[entity.ExperimentColumns.StartTime] = req.StartTime
		}
		if req.EndTime != "" {
			update[entity.ExperimentColumns.EndTime] = req.EndTime
		}
		if err = ss.experimentDao.Update(ctx, req.ExperimentID, update); err != nil {
			logger.Logger.Errorf("[ExperimentService Edit] Update experiment err: %v", err)
			return err
		}

		if err = ss.experimentStepDao.DeleteByExperimentIDTx(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[ExperimentService Edit] Delete experimentStep err: %v", err)
			return err
		}

		if err = tx.CreateInBatches(&experimentSteps, 500).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Edit] CreateInBatches experimentSteps err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		logger.Logger.Errorf("[ExperimentService Edit] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	return nil, nil
}

// Import 文件导入数据库
func (ss ExperimentService) Import(ctx context.Context, userId string, file *multipart.FileHeader) (*dto.Empty, error) {
	experimentFile := entity.ExperimentFiles{
		ID:        uuid.NewString(),
		FileName:  file.Filename,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 获取excel内容
	data, err := getXlsxContent(file)
	if err != nil {
		return nil, err
	}

	//  初始化 UUID 映射表
	uuidMap, experimentGroupIdList := initUUIDMaps(data)

	// 将excel数据转成mysql数据
	experiments, experimentSteps, recipes, recipeMaterialGroups, materials, err := getData(experimentFile.ID,
		userId, data, uuidMap, experimentGroupIdList)
	if err != nil {
		return nil, err
	}

	// 批量生成实验密钥
	experimentSignatureMap, experimentSignatureList := generateExperimentSignature(experiments, experimentSteps,
		recipeMaterialGroups,
		materials)
	for k, _ := range experiments {
		experiments[k].ExperimentSignature = experimentSignatureMap[experiments[k].ID]
	}

	// 获取已存在数据库的实验密钥
	experimentInSignatureList, err := ss.experimentDao.GetExperimentInSignatureList(ctx, experimentSignatureList)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Import] GetExperimentInSignatureList err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 批量生成配方密钥
	recipeSignatureMap, recipeSignatureList := GenerateBatchRecipeSignature(recipeMaterialGroups, materials)
	for k, _ := range recipes {
		recipes[k].RecipeSignature = recipeSignatureMap[recipes[k].ID]
	}

	// 获取已存在数据库的配方密钥
	recipeInSignatureMap, err := ss.recipeDao.GetRecipeInSignatureMap(ctx, recipeSignatureList)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Import] GetExperimentInSignatureList err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 去掉密钥重复的实验、配方数据
	experiments, experimentSteps, recipes, recipeMaterialGroups, materials = filterData(experiments,
		experimentSteps, recipes, recipeMaterialGroups, materials, experimentInSignatureList, recipeInSignatureMap)

	// 每条数据都导入过（文件重复导入）
	if len(experiments) == 0 {
		return nil, utils.NewBusinessError(utils.DuplicateFileImportCode, "")
	}

	// 将数据入库
	if err = WriteData(experimentFile, experiments, experimentSteps, recipes, recipeMaterialGroups, materials); err != nil {
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	return nil, nil
}

// 初始化 UUID 映射表
func initUUIDMaps(data [][]string) (uuidMap map[int]string, experimentGroupIdList map[int][]string) {
	maxCols := GetMaxCols(data)
	uuidMap = make(map[int]string, 0)
	experimentGroupIdList = make(map[int][]string, 0)
	for i := 1; i <= maxCols; i++ {
		uuidMap[i] = uuid.NewString()
		experimentGroupIdList[i] = []string{uuid.NewString(), uuid.NewString()}
	}
	return
}

// getXlsxContent 获取excel内容
func getXlsxContent(file *multipart.FileHeader) ([][]string, error) {
	// 打开文件以读取其内容
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	// 使用 excelize.OpenReader 直接解析文件内容
	f, err := excelize.OpenReader(fileContent)
	if err != nil {
		return nil, err
	}

	// 获取第一个工作表名称
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, err
	}

	// 获取工作表的所有行
	data, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// getData 将excel内容转换成入库数据
func getData(fileId, userId string, data [][]string, uuidMap map[int]string,
	experimentGroupIdList map[int][]string) ([]entity.Experiment, []entity.ExperimentSteps, []entity.Recipes,
	[]entity.RecipeMaterialGroups, []entity.Materials, error) {
	var index, index2, index3, index4 int
	var isIndexSet, isIndexSet2, isIndexSet3, isIndexSet4 bool // 添加一个布尔变量标志是否已设置 index
	var isInGroup bool
	var currentGroupId, groupName string
	var experiments []entity.Experiment
	var experimentSteps []entity.ExperimentSteps
	var recipes []entity.Recipes
	var recipeMaterialGroups []entity.RecipeMaterialGroups
	var materials []entity.Materials
	var err error
	var experimentNum, recipeNum int64
	if err = db.Client.Model(&entity.Experiment{}).Count(&experimentNum).Error; err != nil {
		logger.Logger.Errorf("[getData] Mysql err: %v", err)
		return nil, nil, nil, nil, nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	if err = db.Client.Model(&entity.Recipes{}).Count(&recipeNum).Error; err != nil {
		logger.Logger.Errorf("[getData] Mysql err: %v", err)
		return nil, nil, nil, nil, nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	experimentExistsMap := make(map[string]bool)
	for rowIndex, row := range data {
		if len(row) > 0 {
			if isGroupHeader(row[0]) {
				isInGroup = true
				groupName = row[0]
			}

			if isInGroup && isMaterialData(row[0]) {
				for key, val := range row {
					if isExperimentCondition(row[0], "A") && key > 0 {
						currentGroupId = experimentGroupIdList[key][0]
					}
					if isExperimentCondition(row[0], "B") && key > 0 {
						currentGroupId = experimentGroupIdList[key][1]
					}
					material := entity.Materials{
						ID:                        uuid.NewString(),
						MaterialName:              row[0],
						ExperimentMaterialGroupID: currentGroupId,
						MaterialGroupName:         groupName,
					}
					if key > 0 && val != "" {
						// 转换字符串为 float64
						floatValue, err := strconv.ParseFloat(val, 64)
						if err != nil {
							return nil, nil, nil, nil, nil, err
						}
						material.Percentage = floatValue
						materials = append(materials, material)
					}
				}
			}

			if isExperimentCondition(row[0], "E") {
				if !isIndexSet { // 仅在第一次满足条件时赋值
					index = rowIndex - 1
					isIndexSet = true
				}
			}

			if isExperimentCondition(row[0], "I") {
				if !isIndexSet2 {
					index2 = rowIndex - 3
					isIndexSet2 = true
				}

				if !isIndexSet4 {
					index4 = rowIndex - 2
					isIndexSet4 = true
				}
			}

			if isExperimentCondition(row[0], "D") {
				if !isIndexSet3 {
					index3 = rowIndex - 1
					isIndexSet3 = true
				}
			}

			if isExperiment(row[0]) {
				for key, value := range row {
					if key == 0 {
						continue
					}

					experimentStep := entity.ExperimentSteps{
						ID:                  uuid.NewString(),
						RecipeID:            uuid.NewString(),
						ExperimentID:        uuidMap[key],
						StepName:            row[0],
						StepOrder:           getStepOrder(row[0]),
						ResultValue:         value,
						ExperimentCondition: "",
					}

					// 没有实验步骤的实验要忽略入库
					experimentExistsMap[uuidMap[key]] = true

					recipeMaterialGroups1 := entity.RecipeMaterialGroups{
						RecipeID:                  experimentStep.RecipeID,
						ExperimentMaterialGroupID: experimentGroupIdList[key][0],
						Proportion:                100,
					}
					recipeMaterialGroups2 := entity.RecipeMaterialGroups{
						RecipeID:                  experimentStep.RecipeID,
						ExperimentMaterialGroupID: experimentGroupIdList[key][1],
						Proportion:                100,
					}

					if isExperimentCondition(row[0], "E") {
						experimentStep.ExperimentCondition = data[index][key]
						p1, p2, _ := convertRatioToPercentage(data[index3][key])
						recipeMaterialGroups1.Proportion = p1
						recipeMaterialGroups2.Proportion = p2
					}

					if isExperimentCondition(row[0], "I") {
						experimentStep.ExperimentCondition = data[index2][key]

						p1, p2, _ := convertRatioToPercentage(data[index4][key])
						recipeMaterialGroups1.Proportion = p1
						recipeMaterialGroups2.Proportion = p2
					}

					if isExperimentCondition(row[0], "D") {
						p1, p2, _ := convertRatioToPercentage(data[index3][key])
						recipeMaterialGroups1.Proportion = p1
						recipeMaterialGroups2.Proportion = p2
					}

					if value != "" {
						experimentSteps = append(experimentSteps, experimentStep)
						experimentName := ""
						if len(data) > 1 && len(data[1]) > key {
							experimentName = data[1][key]
						}
						recipes = append(recipes, entity.Recipes{
							ID:              experimentStep.RecipeID,
							RecipeName:      experimentName + experimentStep.StepName,
							RecipeSignature: "",
							Sort:            int(recipeNum) + key*20 + rowIndex,
						})
						switch row[0] {
						case "C1", "C2":
							recipeMaterialGroups = append(recipeMaterialGroups, recipeMaterialGroups1)
						case "C3", "C4":
							recipeMaterialGroups = append(recipeMaterialGroups, recipeMaterialGroups2)
						default:
							recipeMaterialGroups = append(recipeMaterialGroups, recipeMaterialGroups1)
							recipeMaterialGroups = append(recipeMaterialGroups, recipeMaterialGroups2)
						}
					}
				}
			}
		} else {
			// 如果 A 列为空
			fmt.Printf("第 %d 行: 空值\n", rowIndex+1)
			isInGroup = false
		}
	}

	// 初始化实验数据
	for k, id := range uuidMap {
		experimentName := fmt.Sprintf("G%d", int(experimentNum)+k)
		if len(data) > 1 && len(data[1]) > k {
			experimentName = data[1][k]
		}

		// 忽略没有实验步骤的实验
		if _, ok := experimentExistsMap[id]; !ok {
			continue
		}

		experiment := entity.Experiment{
			ID:             id,
			FileID:         fileId,
			ExperimentName: experimentName,
			EntryCategory:  1,
			UserID:         userId,
			Sort:           int(experimentNum) + k,
		}
		experiments = append(experiments, experiment)
	}

	return experiments, experimentSteps, recipes, recipeMaterialGroups, materials, err
}

// 过滤掉重复的数据
func filterData(experiments []entity.Experiment, experimentSteps []entity.ExperimentSteps, recipes []entity.Recipes,
	recipeMaterialGroups []entity.RecipeMaterialGroups, materials []entity.Materials, experimentInSignatureList []string,
	recipeInSignatureMap map[string]string) ([]entity.Experiment, []entity.ExperimentSteps, []entity.Recipes,
	[]entity.RecipeMaterialGroups, []entity.Materials) {
	// 创建一个集合来快速检查签名是否存在
	signatureSet := make(map[string]bool)
	for _, signature := range experimentInSignatureList {
		signatureSet[signature] = true
	}

	// 过滤 experiments
	var filteredExperiments []entity.Experiment
	validExperimentIDs := make(map[string]bool)
	experimentCountMap := make(map[string]bool)
	for _, exp := range experiments {
		if _, exists := signatureSet[exp.ExperimentSignature]; !exists {
			if _, exist := experimentCountMap[exp.ExperimentSignature]; !exist {
				validExperimentIDs[exp.ID] = true
				filteredExperiments = append(filteredExperiments, exp)
			}
		}
		experimentCountMap[exp.ExperimentSignature] = true
	}

	// 过滤 experimentSteps
	var filteredExperimentSteps []entity.ExperimentSteps
	validRecipeIDs := make(map[string]bool)
	for _, step := range experimentSteps {
		if validExperimentIDs[step.ExperimentID] {
			validRecipeIDs[step.RecipeID] = true
			filteredExperimentSteps = append(filteredExperimentSteps, step)
		}
	}

	// 过滤 Recipes
	recipeCountMap := make(map[string]string)
	stepReplaceMap := make(map[string]string)
	var filteredRecipes []entity.Recipes
	for _, recipe := range recipes {
		if _, exists := recipeCountMap[recipe.RecipeSignature]; exists {
			stepReplaceMap[recipe.ID] = recipeCountMap[recipe.RecipeSignature]
			continue
		}
		recipeCountMap[recipe.RecipeSignature] = recipe.ID

		if _, exists := validRecipeIDs[recipe.ID]; !exists {
			continue
		}
		if _, exists := recipeInSignatureMap[recipe.RecipeSignature]; exists {
			stepReplaceMap[recipe.ID] = recipeInSignatureMap[recipe.RecipeSignature]
			continue
		}
		filteredRecipes = append(filteredRecipes, recipe)
	}

	// 替换重复的配方id
	for k, step := range filteredExperimentSteps {
		if _, exists := stepReplaceMap[step.RecipeID]; !exists {
			continue
		}
		filteredExperimentSteps[k].RecipeID = stepReplaceMap[step.RecipeID]
	}

	// 过滤 experimentStepMaterials
	var filteredRecipeMaterialGroups []entity.RecipeMaterialGroups
	validMaterialGroupIDs := make(map[string]bool)
	for _, recipeMaterial := range recipeMaterialGroups {
		if validRecipeIDs[recipeMaterial.RecipeID] {
			validMaterialGroupIDs[recipeMaterial.ExperimentMaterialGroupID] = true
			filteredRecipeMaterialGroups = append(filteredRecipeMaterialGroups, recipeMaterial)
		}
	}

	// 过滤 materials
	var filteredMaterials []entity.Materials
	for _, material := range materials {
		if validMaterialGroupIDs[material.ExperimentMaterialGroupID] {
			filteredMaterials = append(filteredMaterials, material)
		}
	}

	return filteredExperiments, filteredExperimentSteps, filteredRecipes, filteredRecipeMaterialGroups,
		filteredMaterials
}

// generateExperimentSignature 批量生成实验密钥
func generateExperimentSignature(
	experiments []entity.Experiment,
	experimentSteps []entity.ExperimentSteps,
	recipeMaterialGroups []entity.RecipeMaterialGroups,
	materials []entity.Materials,
) (map[string]string, []string) {
	signatureMap, signatureList := make(map[string]string), make([]string, 0)

	// 建立实验步骤映射
	stepMap := make(map[string][]entity.ExperimentSteps)
	for _, step := range experimentSteps {
		stepMap[step.ExperimentID] = append(stepMap[step.ExperimentID], step)
	}

	// 建立步骤-材料组映射
	stepMaterialMap := make(map[string][]entity.RecipeMaterialGroups)
	for _, recipeMaterial := range recipeMaterialGroups {
		stepMaterialMap[recipeMaterial.RecipeID] = append(stepMaterialMap[recipeMaterial.RecipeID], recipeMaterial)
	}

	// 建立材料组-材料映射
	materialMap := make(map[string][]entity.Materials)
	for _, material := range materials {
		materialMap[material.ExperimentMaterialGroupID] = append(materialMap[material.ExperimentMaterialGroupID], material)
	}

	for _, experiment := range experiments {
		var signatureElements []string

		// 处理实验名称
		// signatureElements = append(signatureElements, experiment.ExperimentName)

		// 处理实验步骤
		steps := stepMap[experiment.ID]
		sort.Slice(steps, func(i, j int) bool {
			return steps[i].StepName < steps[j].StepName
		})
		for _, step := range steps {
			signatureElements = append(signatureElements, step.StepName, step.ResultValue, step.ExperimentCondition)

			// 处理实验步骤的材料组
			// 修正：使用 step.ExperimentStepID 作为键
			stepMaterials := stepMaterialMap[step.ID]
			sort.Slice(stepMaterials, func(i, j int) bool {
				return stepMaterials[i].Proportion < stepMaterials[j].Proportion
			})
			for _, stepMaterial := range stepMaterials {
				signatureElements = append(signatureElements, strconv.FormatFloat(stepMaterial.Proportion, 'f', 2, 64))

				// 处理材料
				materials := materialMap[stepMaterial.ExperimentMaterialGroupID]
				sort.Slice(materials, func(i, j int) bool {
					if materials[i].MaterialGroupName == materials[j].MaterialGroupName {
						return materials[i].MaterialName < materials[j].MaterialName
					}
					return materials[i].MaterialGroupName < materials[j].MaterialGroupName
				})
				for _, material := range materials {
					signatureElements = append(signatureElements, material.MaterialGroupName, material.MaterialName, strconv.FormatFloat(material.Percentage, 'f', 2, 64))
				}
			}
		}

		// 计算 MD5 哈希
		signature := md5Hash(strings.Join(signatureElements, "|"))
		signatureMap[experiment.ID] = signature
		signatureList = append(signatureList, signature)
	}

	return signatureMap, signatureList
}

// 计算MD5哈希
func md5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// WriteData 将excel内容入库
func WriteData(experimentFile entity.ExperimentFiles, experiments []entity.Experiment,
	experimentSteps []entity.ExperimentSteps, recipes []entity.Recipes,
	recipeMaterialGroups []entity.RecipeMaterialGroups, materials []entity.Materials) error {
	// 使用事务闭包
	if err := db.Client.Transaction(func(tx *gorm.DB) error {
		maxNum := 500
		if err := tx.Create(experimentFile).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Import] WriteData Create experimentFile err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(materials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Import] WriteData CreateInBatches materials err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experiments, maxNum).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Import] WriteData CreateInBatches experiments err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experimentSteps, maxNum).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Import] WriteData CreateInBatches experimentSteps err: %v", err)
			return err
		}
		if len(recipes) > 0 {
			if err := tx.CreateInBatches(recipes, maxNum).Error; err != nil {
				logger.Logger.Errorf("[ExperimentService Import] WriteData CreateInBatches recipes err: %v", err)
				return err
			}
		}
		if len(recipeMaterialGroups) > 0 {
			if err := tx.CreateInBatches(recipeMaterialGroups, maxNum).Error; err != nil {
				logger.Logger.Errorf("[ExperimentService Import] WriteData CreateInBatches recipeMaterialGroups err: %v", err)
				return err
			}
		}

		// 执行存储过程去重Materials表
		if err := tx.Exec("CALL remove_duplicate_material_groups()").Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Import] remove_duplicate_material_groups err: %v", err)
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Errorf("[WriteData] Mysql err: %v", err)
		return utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	return nil
}

// convertRatioToPercentage 将"100/30"格式的字符串转换为百分比小数
func convertRatioToPercentage(input string) (float64, float64, error) {
	if input == "" {
		return 0, 0, nil
	}

	// 分割字符串
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("输入格式错误: 需要形如 '100/30' 的字符串")
	}

	// 将数值部分转换为浮点数
	numerator, err1 := strconv.ParseFloat(parts[0], 64)
	denominator, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("数值转换失败: %v, %v", err1, err2)
	}

	// 检查分母是否为零
	if numerator+denominator == 0 {
		return 0, 0, fmt.Errorf("分子和分母之和不能为零")
	}

	// 计算百分比小数（分子占总和的比例）
	percentage1 := numerator / (numerator + denominator) * 100
	percentage2 := denominator / (numerator + denominator) * 100
	return percentage1, percentage2, nil
}

// isGroupHeader 检查是否是组标题
func isGroupHeader(cellValue string) bool {
	groupHeaders := []string{"树脂组", "固化剂组"} // 可扩展
	for _, header := range groupHeaders {
		if strings.Contains(cellValue, header) {
			return true
		}
	}
	return false
}

// isMaterialData 检查是否是材料数据
func isMaterialData(cellValue string) bool {
	return strings.HasPrefix(cellValue, "A") || strings.HasPrefix(cellValue, "B")
}

// isExperiment 检查是否是实验数据
func isExperiment(cellValue string) bool {
	return strings.HasPrefix(cellValue, "C") ||
		strings.HasPrefix(cellValue, "D") ||
		strings.HasPrefix(cellValue, "E")
	// || strings.HasPrefix(cellValue, "I")  //数据不确定先不录
}

func isExperimentCondition(cellValue, prefix string) bool {
	return strings.HasPrefix(cellValue, prefix)
}

// getStepOrder 获取步骤顺序
func getStepOrder(cellValue string) int {
	switch {
	case strings.HasPrefix(cellValue, "C"):
		return 400 - ExtractTrailingNumberRegex(cellValue)
	case strings.HasPrefix(cellValue, "D"):
		return 300 - ExtractTrailingNumberRegex(cellValue)
	case strings.HasPrefix(cellValue, "E"):
		return 200 - ExtractTrailingNumberRegex(cellValue)
	case strings.HasPrefix(cellValue, "I"):
		return 100 - ExtractTrailingNumberRegex(cellValue)
	default:
		return 0
	}
}

// ExtractTrailingNumberRegex 使用正则表达式提取字符串末尾的数字，若没有则返回 0
func ExtractTrailingNumberRegex(s string) int {
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(s)
	var num int
	if match != "" {
		fmt.Sscanf(match, "%d", &num)
	}
	return num
}

// GetMaxCols 获取数据最大列数
func GetMaxCols(data [][]string) int {
	maxCols := 0
	for _, row := range data {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	return maxCols - 1 // 第一行不算
}
