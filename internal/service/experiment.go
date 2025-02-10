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
	"fmt"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ExperimentService struct {
	experimentDao     *dao.ExperimentDao
	experimentStepDao *dao.ExperimentStepDao
	materialDao       *dao.MaterialDao
}

// NewExperimentService 创建一个新的 ExperimentService 实例
func NewExperimentService() *ExperimentService {
	return &ExperimentService{
		experimentDao:     dao.NewExperimentDao(),
		experimentStepDao: dao.NewExperimentStepDao(),
		materialDao:       dao.NewMaterialDao(),
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
		logger.Logger.Errorf("[List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
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
		logger.Logger.Errorf("[List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}
	for _, file := range files {
		fileMap[file.ID] = file.FileName
	}

	// 一次性查询所有实验步骤
	var steps []entity.ExperimentSteps
	if err := db.Client.WithContext(ctx).Where("experiment_id IN (?)", experimentIDs).
		Order("step_order DESC").
		Find(&steps).Error; err != nil {
		logger.Logger.Errorf("[List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}

	// 获取所有步骤ID
	stepIDs := make([]string, len(steps))
	for i, step := range steps {
		stepIDs[i] = step.ID
	}

	// 一次性查询所有材料组
	var materialGroups []entity.ExperimentStepMaterial
	if err := db.Client.WithContext(ctx).
		Where("experiment_step_id IN (?)", stepIDs).
		Find(&materialGroups).Error; err != nil {
		logger.Logger.Errorf("[List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}

	// 获取所有材料ID
	materialGroupIDs := make([]string, len(materialGroups))
	for i, group := range materialGroups {
		materialGroupIDs[i] = group.ExperimentMaterialGroupID
	}

	// 一次性查询所有材料
	var materials []entity.Materials
	if err := db.Client.WithContext(ctx).
		Where("experiment_material_group_id IN (?)", materialGroupIDs).
		Find(&materials).Error; err != nil {
		logger.Logger.Errorf("[List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}

	// 构建查询结果映射
	stepMap := make(map[string][]data.ExperimentStepData)
	materialGroupMap := make(map[string][]data.MaterialGroupData)

	// 填充 stepMap
	for _, step := range steps {
		stepMap[step.ExperimentID] = append(stepMap[step.ExperimentID], data.ExperimentStepData{
			StepID:              step.ID,
			StepName:            step.StepName,
			StepNameDescription: conf.StepNameData[step.StepName],
			ExperimentCondition: step.ExperimentCondition,
			ResultValue:         step.ResultValue,
			StepOrder:           step.StepOrder,
		})
	}

	// 填充 materialGroupMap
	groupMap := make(map[string]string)
	for _, material := range materials {
		groupMap[material.ExperimentMaterialGroupID] = material.MaterialGroupName
	}

	for _, group := range materialGroups {
		materialGroupMap[group.ExperimentStepID] = append(materialGroupMap[group.ExperimentStepID], data.MaterialGroupData{
			MaterialGroupID:   group.ExperimentMaterialGroupID,
			MaterialGroupName: groupMap[group.ExperimentMaterialGroupID],
			Proportion:        group.Proportion,
		})
	}

	// 填充每个材料组的材料信息
	materialMap := make(map[string][]data.MaterialData)
	for _, material := range materials {
		materialMap[material.ExperimentMaterialGroupID] = append(materialMap[material.ExperimentMaterialGroupID], data.MaterialData{
			MaterialID:   material.ID,
			MaterialName: material.MaterialName,
			Percentage:   material.Percentage,
		})
	}

	// 组装最终的实验数据
	var experimentDataList []data.ExperimentData
	for _, experiment := range experiments {
		stepsData := stepMap[experiment.ID]

		// 填充每个实验步骤的材料组和材料
		for i := range stepsData {
			stepData := &stepsData[i]
			materialGroups := materialGroupMap[stepData.StepID]

			// 填充每个材料组的材料信息
			for j := range materialGroups {
				materialGroupData := &materialGroups[j]
				materialGroupData.Materials = materialMap[materialGroupData.MaterialGroupID]
			}

			stepData.MaterialGroups = materialGroups
		}

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
		logger.Logger.Errorf("[Delete] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}

	if experiment == nil || experiment.ID == "" {
		return nil, utils.NewBusinessError(utils.ExperimentDoesNotExistCode)
	}

	if err = db.Client.Transaction(func(tx *gorm.DB) error {
		experimentStepIdList := make([]string, 0)
		if err = tx.WithContext(ctx).Where(entity.ExperimentStepsColumns.ExperimentID+" = ?", req.ExperimentID).
			Pluck(entity.ExperimentStepsColumns.ID, &experimentStepIdList).Error; err != nil {
			return err
		}

		// 删除实验关联数据
		if err = tx.WithContext(ctx).
			Where(entity.ExperimentStepMaterialColumns.ExperimentStepID+" in ?", experimentStepIdList).
			Delete(&entity.ExperimentStepMaterial{}).Error; err != nil {
			logger.Logger.Errorf("[Edit] delete ExperimentStepMaterial err: %v", err)
			return err
		}

		if err := ss.experimentStepDao.DeleteByExperimentIDTx(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[Delete] Delete experimentStep err: %v", err)
			return err
		}

		if err := ss.experimentDao.DeleteByIDTx(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[Delete] Delete experiment err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		logger.Logger.Errorf("[Delete] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}

	return nil, err
}

// Add 添加实验数据
func (ss *ExperimentService) Add(ctx context.Context, userId string, req *data.ExperimentAddRequest) (*dto.Empty, error) {
	materials, experiment, experimentSteps, experimentStepMaterials := getInsertData(userId, req)
	// 使用事务闭包
	if err := db.Client.Transaction(func(tx *gorm.DB) error {
		maxNum := 500
		if err := tx.Create(experiment).Error; err != nil {
			logger.Logger.Errorf("[Add] CreateInBatches experiment err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experimentSteps, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Add] CreateInBatches experimentSteps err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experimentStepMaterials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Add] CreateInBatches experimentStepMaterials err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(materials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Add] CreateInBatches materials err: %v", err)
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Errorf("[Add] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}

	return nil, nil
}

// Edit 修改实验数据
func (ss *ExperimentService) Edit(ctx context.Context, req *data.ExperimentUpdateRequest) (*dto.Empty, error) {
	experiment, err := ss.experimentDao.GetByID(ctx, req.ExperimentID)
	if err != nil {
		logger.Logger.Errorf("[Edit] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
	}

	if experiment == nil || experiment.ID == "" {
		return nil, utils.NewBusinessError(utils.ExperimentDoesNotExistCode)
	}

	materials, experimentSteps, experimentStepMaterials := getUpdateData(experiment.CreatedAt, req)
	if err = db.Client.Transaction(func(tx *gorm.DB) error {
		// 修改实验
		update := map[string]interface{}{
			entity.ExperimentColumns.ExperimentName: req.ExperimentName,
			entity.ExperimentColumns.Experimenter:   req.Experimenter,
			entity.ExperimentColumns.StartTime:      nil,
			entity.ExperimentColumns.EndTime:        nil,
		}
		if req.StartTime != "" {
			update[entity.ExperimentColumns.StartTime] = req.StartTime
		}
		if req.EndTime != "" {
			update[entity.ExperimentColumns.EndTime] = req.EndTime
		}
		if err = ss.experimentDao.Update(ctx, req.ExperimentID, update); err != nil {
			logger.Logger.Errorf("[Edit] Update experiment err: %v", err)
			return err
		}

		experimentStepIdList := make([]string, 0)
		if err = tx.WithContext(ctx).Model(&entity.ExperimentSteps{}).Where(entity.ExperimentStepsColumns.ExperimentID+" = ?",
			req.ExperimentID).
			Pluck(entity.ExperimentStepsColumns.ID, &experimentStepIdList).Error; err != nil {
			return err
		}

		// 删除实验关联数据
		if err = tx.WithContext(ctx).
			Where(entity.ExperimentStepMaterialColumns.ExperimentStepID+" in ?", experimentStepIdList).
			Delete(&entity.ExperimentStepMaterial{}).Error; err != nil {
			logger.Logger.Errorf("[Edit] delete ExperimentStepMaterial err: %v", err)
			return err
		}

		if err = ss.experimentStepDao.DeleteByExperimentIDTx(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[Edit] Delete experimentStep err: %v", err)
			return err
		}

		// 插入新的实验关联数据
		maxNum := 500
		if err = tx.CreateInBatches(experimentSteps, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Edit] CreateInBatches experimentSteps err: %v", err)
			return err
		}
		if err = tx.CreateInBatches(experimentStepMaterials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Edit] CreateInBatches experimentStepMaterials err: %v", err)
			return err
		}
		if err = tx.CreateInBatches(materials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Edit] CreateInBatches materials err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		logger.Logger.Errorf("[Edit] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.InternalError)
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
	materials, experiments, experimentSteps, experimentStepMaterial, err := getData(experimentFile.ID, userId, data,
		uuidMap, experimentGroupIdList)

	// 将数据入库
	if err = WriteData(experimentFile, materials, experiments, experimentSteps, experimentStepMaterial); err != nil {
		return nil, err
	}

	return nil, nil
}

// getInsertData 将插入数据内容转换成入库数据
func getInsertData(userId string, data *data.ExperimentAddRequest) ([]entity.Materials, entity.Experiment, []entity.ExperimentSteps,
	[]entity.ExperimentStepMaterial) {
	var materials []entity.Materials
	var experimentSteps []entity.ExperimentSteps
	var experimentStepMaterials []entity.ExperimentStepMaterial
	now := time.Now()
	layout := "2006-01-02T15:04:05Z"
	experiment := entity.Experiment{
		ID:             uuid.NewString(),
		FileID:         "",
		ExperimentName: "",
		EntryCategory:  2,
		Sort:           data.Sort,
		Experimenter:   data.Experimenter,
		UserID:         userId,
		StartTime:      nil,
		EndTime:        nil,
		CreatedAt:      now,
	}
	if data.StartTime != "" {
		startTime, err := time.Parse(layout, data.StartTime)
		if err != nil {
			logger.Logger.Errorf("[getInsertData] time.Parse Error: %v", err)
		} else {
			experiment.StartTime = &startTime
		}
	}
	if data.EndTime != "" {
		endTime, err := time.Parse(layout, data.EndTime)
		if err != nil {
			logger.Logger.Errorf("[getInsertData] time.Parse Error: %v", err)
		} else {
			experiment.EndTime = &endTime
		}
	}
	for _, step := range data.Steps {
		experimentStep := entity.ExperimentSteps{
			ID:                  uuid.NewString(),
			ExperimentID:        experiment.ID,
			StepOrder:           step.StepOrder,
			StepName:            step.StepName,
			ResultValue:         step.ResultValue,
			ExperimentCondition: step.ExperimentCondition,
			CreatedAt:           now,
		}
		experimentSteps = append(experimentSteps, experimentStep)

		for _, materialGroup := range step.MaterialGroups {
			ExperimentStepMaterial := entity.ExperimentStepMaterial{
				ExperimentStepID:          experimentStep.ID,
				ExperimentMaterialGroupID: uuid.NewString(),
				Proportion:                materialGroup.Proportion,
			}
			experimentStepMaterials = append(experimentStepMaterials, ExperimentStepMaterial)
			for _, material := range materialGroup.Materials {
				materials = append(materials, entity.Materials{
					ID:                        uuid.NewString(),
					MaterialName:              material.MaterialName,
					ExperimentMaterialGroupID: ExperimentStepMaterial.ExperimentMaterialGroupID,
					MaterialGroupName:         materialGroup.MaterialGroupName,
					Percentage:                material.Percentage,
					CreatedAt:                 now,
				})
			}
		}
	}

	materials, experimentStepMaterials = DeduplicateMaterials(materials, experimentStepMaterials)

	return materials, experiment, experimentSteps, experimentStepMaterials
}

// getUpdateData 将修改的数据内容转换成入库数据
func getUpdateData(createdAt time.Time, data *data.ExperimentUpdateRequest) ([]entity.Materials, []entity.ExperimentSteps,
	[]entity.ExperimentStepMaterial) {
	var materials []entity.Materials
	var experimentSteps []entity.ExperimentSteps
	var experimentStepMaterials []entity.ExperimentStepMaterial
	for _, step := range data.Steps {
		experimentStep := entity.ExperimentSteps{
			ID:                  uuid.NewString(),
			ExperimentID:        data.ExperimentID,
			StepOrder:           step.StepOrder,
			StepName:            step.StepName,
			ResultValue:         step.ResultValue,
			ExperimentCondition: step.ExperimentCondition,
			CreatedAt:           createdAt,
		}
		experimentSteps = append(experimentSteps, experimentStep)

		for _, materialGroup := range step.MaterialGroups {
			ExperimentStepMaterial := entity.ExperimentStepMaterial{
				ExperimentStepID:          experimentStep.ID,
				ExperimentMaterialGroupID: uuid.NewString(),
				Proportion:                materialGroup.Proportion,
			}
			experimentStepMaterials = append(experimentStepMaterials, ExperimentStepMaterial)
			for _, material := range materialGroup.Materials {
				materials = append(materials, entity.Materials{
					ID:                        uuid.NewString(),
					MaterialName:              material.MaterialName,
					ExperimentMaterialGroupID: ExperimentStepMaterial.ExperimentMaterialGroupID,
					MaterialGroupName:         materialGroup.MaterialGroupName,
					Percentage:                material.Percentage,
					CreatedAt:                 createdAt,
				})
			}
		}
	}

	materials, experimentStepMaterials = DeduplicateMaterials(materials, experimentStepMaterials)

	return materials, experimentSteps, experimentStepMaterials
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

// DeduplicateMaterials 去重 materials 并更新 experimentStepMaterials
func DeduplicateMaterials(materials []entity.Materials, experimentStepMaterials []entity.ExperimentStepMaterial) ([]entity.Materials, []entity.ExperimentStepMaterial) {
	// 用于存储唯一的材料组，key 是材料组合签名，value 是 original_group
	materialMap := make(map[string]string)
	// 记录去重后的材料数据
	var deduplicatedMaterials []entity.Materials
	// 记录 group_id 映射关系
	groupMapping := make(map[string]string)

	// 遍历 materials，找出重复项
	for _, material := range materials {
		// 生成唯一的材料签名（确保 percentage 精度一致）
		signature := fmt.Sprintf("%s:%s:%.2f", material.MaterialName, material.MaterialGroupName, material.Percentage)

		// 检查是否已有该材料组合
		if originalGroup, exists := materialMap[signature]; exists {
			// 记录重复项的 group_id 替换关系
			groupMapping[material.ExperimentMaterialGroupID] = originalGroup
		} else {
			// 该材料组合首次出现，存入 map 并加入去重后的列表
			materialMap[signature] = material.ExperimentMaterialGroupID
			deduplicatedMaterials = append(deduplicatedMaterials, material)
		}
	}

	// 遍历 experimentStepMaterials，替换 group ID
	for i, stepMaterial := range experimentStepMaterials {
		if newGroupID, exists := groupMapping[stepMaterial.ExperimentMaterialGroupID]; exists {
			experimentStepMaterials[i].ExperimentMaterialGroupID = newGroupID
		}
	}

	logger.Logger.Infof("[Deduplication] Removed %d duplicate materials", len(materials)-len(deduplicatedMaterials))
	return deduplicatedMaterials, experimentStepMaterials
}

// getData 将excel内容转换成入库数据
func getData(fileId, userId string, data [][]string, uuidMap map[int]string,
	experimentGroupIdList map[int][]string) ([]entity.
	Materials,
	[]entity.Experiment, []entity.ExperimentSteps,
	[]entity.ExperimentStepMaterial, error) {
	var index, index2, index3, index4 int
	var isIndexSet, isIndexSet2, isIndexSet3, isIndexSet4 bool // 添加一个布尔变量标志是否已设置 index
	var isInGroup bool
	var currentGroupId, groupName string
	var materials []entity.Materials
	var experiments []entity.Experiment
	var experimentSteps []entity.ExperimentSteps
	var experimentStepMaterials []entity.ExperimentStepMaterial
	var err error
	now := time.Now()

	var num int64
	if err = db.Client.Model(&entity.Experiment{}).Count(&num).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	// 初始化实验数据
	for k, id := range uuidMap {
		experiment := entity.Experiment{
			ID:             id,
			FileID:         fileId,
			ExperimentName: fmt.Sprintf("G%d", num+int64(k)),
			EntryCategory:  1,
			Experimenter:   "",
			UserID:         userId,
			Sort:           int(num) + k,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		experiments = append(experiments, experiment)
	}

	for rowIndex, row := range data {
		// 检查是否存在 A 列数据
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
						CreatedAt:                 now,
						UpdatedAt:                 now,
					}
					if key > 0 && val != "" {
						// 转换字符串为 float64
						floatValue, err := strconv.ParseFloat(val, 64)
						if err != nil {
							return nil, nil, nil, nil, err
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
						ID:           uuid.NewString(),
						ExperimentID: uuidMap[key],
						StepName:     row[0],
						StepOrder:    getOrder(row[0]),
						ResultValue:  value,
						CreatedAt:    now,
						UpdatedAt:    now,
					}

					experimentStepMaterial1 := entity.ExperimentStepMaterial{
						ExperimentStepID:          experimentStep.ID,
						ExperimentMaterialGroupID: experimentGroupIdList[key][0],
						Proportion:                100,
					}
					experimentStepMaterial2 := entity.ExperimentStepMaterial{
						ExperimentStepID:          experimentStep.ID,
						ExperimentMaterialGroupID: experimentGroupIdList[key][1],
						Proportion:                100,
					}

					if isExperimentCondition(row[0], "E") {
						experimentStep.ExperimentCondition = data[index][key]
						p1, p2, _ := convertRatioToPercentage(data[index3][key])
						experimentStepMaterial1.Proportion = p1
						experimentStepMaterial2.Proportion = p2
					}

					if isExperimentCondition(row[0], "I") {
						experimentStep.ExperimentCondition = data[index2][key]

						p1, p2, _ := convertRatioToPercentage(data[index4][key])
						experimentStepMaterial1.Proportion = p1
						experimentStepMaterial2.Proportion = p2
					}

					if isExperimentCondition(row[0], "D") {
						p1, p2, _ := convertRatioToPercentage(data[index3][key])
						experimentStepMaterial1.Proportion = p1
						experimentStepMaterial2.Proportion = p2
					}

					if value != "" {
						experimentSteps = append(experimentSteps, experimentStep)
						switch row[0] {
						case "C1", "C2":
							experimentStepMaterials = append(experimentStepMaterials, experimentStepMaterial1)
						case "C3", "C4":
							experimentStepMaterials = append(experimentStepMaterials, experimentStepMaterial2)
						default:
							experimentStepMaterials = append(experimentStepMaterials, experimentStepMaterial1)
							experimentStepMaterials = append(experimentStepMaterials, experimentStepMaterial2)
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

	// 去重
	materials, experimentStepMaterials = DeduplicateMaterials(materials, experimentStepMaterials)

	return materials, experiments, experimentSteps, experimentStepMaterials, err
}

// WriteData 将excel内容入库
func WriteData(experimentFile entity.ExperimentFiles, materials []entity.Materials, experiments []entity.Experiment,
	experimentSteps []entity.ExperimentSteps,
	experimentStepMaterials []entity.ExperimentStepMaterial) error {
	// 使用事务闭包
	if err := db.Client.Transaction(func(tx *gorm.DB) error {
		maxNum := 500
		if err := tx.Create(experimentFile).Error; err != nil {
			logger.Logger.Errorf("[Import] WriteData Create experimentFile err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(materials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Import] WriteData CreateInBatches materials err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experiments, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Import] WriteData CreateInBatches experiments err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experimentSteps, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Import] WriteData CreateInBatches experimentSteps err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experimentStepMaterials, maxNum).Error; err != nil {
			logger.Logger.Errorf("[Import] WriteData CreateInBatches experimentStepMaterials err: %v", err)
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Errorf("[WriteData] Mysql err: %v", err)
		return utils.NewBusinessError(utils.InternalError)
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

// getOrder 获取步骤顺序
func getOrder(cellValue string) int {
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
