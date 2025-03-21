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
	"gorm.io/gorm"
	"time"
)

type ExperimentService struct {
	experimentDao     *dao.ExperimentDao
	experimentStepDao *dao.ExperimentStepDao
	materialDao       *dao.MaterialDao
	materialGroupDao  *dao.MaterialGroupDao
	recordDao         *dao.RecordDao
	userDao           *dao.UserDao
}

// NewExperimentService 创建一个新的 ExperimentService 实例
func NewExperimentService() *ExperimentService {
	return &ExperimentService{
		experimentDao:     dao.NewExperimentDao(),
		experimentStepDao: dao.NewExperimentStepDao(),
		materialDao:       dao.NewMaterialDao(),
		materialGroupDao:  dao.NewMaterialGroupDao(),
		recordDao:         dao.NewRecordDao(),
		userDao:           dao.NewUserDao(),
	}
}

// List 查询实验及其相关数据，包括实验步骤和材料组
func (ss *ExperimentService) List(ctx context.Context, req *data.ExperimentListRequest) (*data.ExperimentListResponse, error) {
	res := &data.ExperimentListResponse{
		Total: 0,
		List:  make([]data.ExperimentData, 0),
	}
	// 查询实验及其基础信息
	experiments, total, err := ss.experimentDao.Query(ctx, req.Page, req.PageSize, map[string]interface{}{
		entity.ExperimentColumns.ExperimentName: req.ExperimentName,
		entity.ExperimentColumns.Experimenter:   req.Experimenter,
		entity.ExperimentColumns.TaskID:         req.TaskId,
		entity.ExperimentColumns.Status:         req.Status,
		"startTime":                             req.StartTime,
		"endTIme":                               req.EndTime,
	})
	if err != nil {
		logger.Logger.Errorf("[ExperimentService List] Mysql err: %v", err)
		return res, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 获取所有实验ID
	userIds := make([]string, 0)
	for _, experiment := range experiments {
		userIds = append(userIds, experiment.UserID)
	}

	var users []entity.Users
	userMap := make(map[string]string, 0)
	userIds = utils.RemoveDuplicates(userIds)
	if err = db.Client.WithContext(ctx).
		Where(entity.UsersColumns.ID+" IN (?)", userIds).Find(&users).Error; err != nil {
		logger.Logger.Errorf("[ExperimentService List] Mysql err: %v", err)
		return res, utils.NewBusinessError(utils.InternalError, "")
	}
	for _, user := range users {
		userMap[user.ID] = user.Username
	}

	for _, experiment := range experiments {
		username := ""
		if _, ok := userMap[experiment.UserID]; ok {
			username = userMap[experiment.UserID]
		}
		experimentData := data.ExperimentData{
			ExperimentID:   experiment.ID,
			ExperimentName: experiment.ExperimentName,
			EntryCategory:  experiment.EntryCategory,
			Experimenter:   experiment.Experimenter,
			UserID:         experiment.UserID,
			Username:       username,
			Status:         experiment.Status,
			CreatedAt:      experiment.CreatedAt.Format(time.DateTime),
		}
		if experiment.StartTime != nil {
			experimentData.StartTime = experiment.StartTime.Format(time.DateTime)
		}
		if experiment.EndTime != nil {
			experimentData.EndTime = experiment.EndTime.Format(time.DateTime)
		}
		res.List = append(res.List, experimentData)
	}
	res.Total = total

	// 返回最终结果
	return res, nil
}

// Delete 删除实验数据
func (ss *ExperimentService) Delete(ctx context.Context, userId string, req *data.ExperimentDeleteRequest) (*dto.Empty,
	error) {
	experiment, err := ss.experimentDao.GetByID(ctx, req.ExperimentID)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Delete] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	if experiment == nil || experiment.ID == "" {
		return nil, utils.NewBusinessError(utils.ExperimentDoesNotExistCode, "")
	}

	if err = db.Client.Transaction(func(tx *gorm.DB) error {
		materialGroupIdList := make([]string, 0)
		if err = tx.WithContext(ctx).Where(entity.MaterialGroupsColumns.ExperimentID+" = ?", req.ExperimentID).
			Pluck(entity.MaterialGroupsColumns.ID, &materialGroupIdList).Error; err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] Get materialGroupIdList err: %v", err)
			return err
		}

		if err = ss.experimentDao.DeleteByID(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] Delete experiment err: %v", err)
			return err
		}

		if err = ss.experimentStepDao.DeleteByExperimentID(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] Delete experimentStep err: %v", err)
			return err
		}

		if err = ss.materialGroupDao.DeleteByExperimentID(ctx, tx, req.ExperimentID); err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] Delete materialGroup err: %v", err)
			return err
		}

		if err = ss.materialDao.DeleteByGroupIdListTx(ctx, tx, materialGroupIdList); err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] Delete material err: %v", err)
			return err
		}

		// 记录操作
		if err = ss.recordDao.Insert(ctx, tx, &entity.OperationRecord{
			ID:     uuid.NewString(),
			OpType: conf.DeleteExperiment,
			UserID: userId,
			Desc:   fmt.Sprintf("删除实验: %s", experiment.ExperimentName),
		}); err != nil {
			logger.Logger.Errorf("[ExperimentService Delete] create record: %v", err)
			return err
		}

		return nil
	}); err != nil {
		logger.Logger.Errorf("[ExperimentService Delete] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	return nil, err

	return nil, nil
}

// Info 获取实验数据详情
func (ss *ExperimentService) Info(ctx context.Context, req *data.ExperimentInfoRequest) (*data.ExperimentInfoResponse, error) {
	experiment, err := ss.experimentDao.GetByID(ctx, req.ExperimentID)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Delete] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	if experiment == nil || experiment.ID == "" {
		return nil, utils.NewBusinessError(utils.ExperimentDoesNotExistCode, "")
	}

	userInfo, err := ss.userDao.GetByID(ctx, experiment.UserID)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Info] userDao.GetByID Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	res := &data.ExperimentInfoResponse{
		ExperimentID:   experiment.ID,
		ExperimentName: experiment.ExperimentName,
		EntryCategory:  experiment.EntryCategory,
		Experimenter:   experiment.Experimenter,
		UserID:         experiment.UserID,
		Status:         experiment.Status,
		CreatedAt:      experiment.CreatedAt.Format(time.DateTime),
		StepInfo:       make([]data.ExperimentStepInfo, 0),
		MaterialGroups: make([]data.MaterialGroupInfo, 0),
	}
	if experiment.StartTime != nil {
		res.StartTime = experiment.StartTime.Format(time.DateTime)
	}
	if experiment.EndTime != nil {
		res.EndTime = experiment.EndTime.Format(time.DateTime)
	}
	if userInfo != nil {
		res.Username = userInfo.Username
	}

	experimentSteps, err := ss.experimentStepDao.GetByExperimentID(ctx, req.ExperimentID)
	if err != nil {
		logger.Logger.Errorf("[ExperimentService Info] experimentStepDao.GetByExperimentID Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	for _, step := range experimentSteps {
		res.StepInfo = append(res.StepInfo, data.ExperimentStepInfo{
			StepID:        step.ID,
			StepName:      step.StepName,
			StepCondition: *step.StepCondition,
			StepCategory:  step.StepCategory,
			ResultValue:   step.ResultValue.String(),
		})
	}

	// 查询该配方的所有材料组
	var materialGroups []entity.MaterialGroups
	if err = db.Client.WithContext(ctx).
		Where(entity.MaterialGroupsColumns.ExperimentID+" = ?", req.ExperimentID).
		Find(&materialGroups).Error; err != nil {
		logger.Logger.Errorf("[ExperimentService Info]  Get MaterialGroups Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 组织材料组数据
	materialGroupIDs := make([]string, len(materialGroups))
	materialGroupMap := make(map[string]*data.MaterialGroupInfo)
	for _, group := range materialGroups {
		materialGroupIDs = append(materialGroupIDs, group.ID)
		materialGroup := &data.MaterialGroupInfo{
			MaterialGroupID:       group.ID,
			Proportion:            group.Proportion,
			MaterialGroupName:     group.MaterialGroupName,
			MaterialGroupCategory: group.MaterialGroupCategory,
			Materials:             []data.MaterialInfo{},
		}
		materialGroupMap[group.ID] = materialGroup
	}

	// 查询所有材料
	var materials []entity.Materials
	if err = db.Client.WithContext(ctx).
		Where(entity.MaterialsColumns.MaterialGroupID+" IN (?)", materialGroupIDs).
		Find(&materials).Error; err != nil {
		logger.Logger.Errorf("[ExperimentService Info] Get materials Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 填充材料数据到材料组中
	for _, material := range materials {
		if materialGroup, exists := materialGroupMap[material.MaterialGroupID]; exists {
			materialGroup.Materials = append(materialGroup.Materials, data.MaterialInfo{
				MaterialID: material.ID,
				MaterialData: data.MaterialData{
					MaterialName: material.MaterialName,
					Percentage:   material.Proportion,
				},
			})
		}
	}

	// 将材料组添加到对应的配方中
	for _, group := range materialGroups {
		if materialGroup, exists := materialGroupMap[group.ID]; exists {
			res.MaterialGroups = append(res.MaterialGroups, *materialGroup)
		}
	}

	return res, nil
}

//
// // Add 添加实验数据
// func (ss *ExperimentService) Add(ctx context.Context, userId string, req *data.ExperimentAddRequest) (*dto.Empty, error) {
// 	// experiment := entity.Experiment{
// 	// 	ID:             uuid.NewString(),
// 	// 	RecipeID:       req.RecipeID,
// 	// 	FileID:         "",
// 	// 	ExperimentName: req.ExperimentName,
// 	// 	EntryCategory:  2,
// 	// 	Sort:           req.Sort,
// 	// 	Experimenter:   req.Experimenter,
// 	// 	UserID:         userId,
// 	// }
// 	// layout := "2006-01-02T15:04:05Z"
// 	// if req.StartTime != "" {
// 	// 	startTime, _ := time.Parse(layout, req.StartTime)
// 	// 	experiment.StartTime = &startTime
// 	// }
// 	// if req.EndTime != "" {
// 	// 	endTime, _ := time.Parse(layout, req.EndTime)
// 	// 	experiment.EndTime = &endTime
// 	// }
// 	//
// 	// experimentSteps := make([]entity.ExperimentSteps, 0)
// 	// for _, step := range req.Steps {
// 	// 	experimentStep := entity.ExperimentSteps{
// 	// 		ID:            uuid.NewString(),
// 	// 		ExperimentID:  experiment.ID,
// 	// 		StepOrder:     step.StepOrder,
// 	// 		StepName:      step.StepName,
// 	// 		ResultValue:   step.ResultValue,
// 	// 		StepCondition: step.StepCondition,
// 	// 	}
// 	// 	experimentSteps = append(experimentSteps, experimentStep)
// 	// }
// 	//
// 	// // 获取配方数据
// 	// materialGroups, materials, err := ss.recipeDao.GetMaterialById(ctx, req.RecipeID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// // 生成试验签名
// 	// experimentSignatureMap, _ := generateExperimentSignature([]entity.Experiment{experiment},
// 	// 	experimentSteps, materialGroups, materials)
// 	// experiment.ExperimentSignature = experimentSignatureMap[experiment.ID]
// 	//
// 	// experimentList, err := ss.experimentDao.GetByExperimentSignature(ctx, experiment.ExperimentSignature)
// 	// if err != nil {
// 	// 	logger.Logger.Errorf("[ExperimentService Add] Mysql err: %v", err)
// 	// 	return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
// 	// }
// 	//
// 	// if len(experimentList) > 0 {
// 	// 	return nil, utils.NewBusinessError(utils.DuplicateExperimentFormatCode, "", experimentList[0].ExperimentName)
// 	// }
// 	//
// 	// // 使用事务闭包
// 	// if err := db.Client.Transaction(func(tx *gorm.DB) error {
// 	// 	maxNum := 500
// 	// 	if err := tx.Create(&experiment).Error; err != nil {
// 	// 		logger.Logger.Errorf("[ExperimentService Add] CreateInBatches experiment err: %v", err)
// 	// 		return err
// 	// 	}
// 	// 	if err := tx.CreateInBatches(experimentSteps, maxNum).Error; err != nil {
// 	// 		logger.Logger.Errorf("[ExperimentService Add] CreateInBatches experimentSteps err: %v", err)
// 	// 		return err
// 	// 	}
// 	// 	// 记录操作
// 	// 	if err = ss.recordDao.Insert(ctx, tx, &entity.OperationRecord{
// 	// 		ID:     uuid.NewString(),
// 	// 		OpType: conf.AddExperiment,
// 	// 		UserID: userId,
// 	// 		Desc:   fmt.Sprintf("添加实验: %s", experiment.ExperimentName),
// 	// 	}); err != nil {
// 	// 		logger.Logger.Errorf("[ExperimentService Add] Error create record: %v", err)
// 	// 		return err
// 	// 	}
// 	// 	return nil
// 	// }); err != nil {
// 	// 	logger.Logger.Errorf("[ExperimentService Add] Mysql err: %v", err)
// 	// 	return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
// 	// }
//
// 	return nil, nil
// }
//
// // Edit 修改实验数据
// func (ss *ExperimentService) Edit(ctx context.Context, userId string, req *data.ExperimentUpdateRequest) (*dto.Empty, error) {
// 	// experiment, err := ss.experimentDao.GetByID(ctx, req.ExperimentID)
// 	// if err != nil {
// 	// 	logger.Logger.Errorf("[ExperimentService Edit] Mysql err: %v", err)
// 	// 	return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
// 	// }
// 	//
// 	// if experiment == nil || experiment.ID == "" {
// 	// 	return nil, utils.NewBusinessError(utils.ExperimentDoesNotExistCode, "")
// 	// }
// 	//
// 	// experimentSteps := make([]entity.ExperimentSteps, 0)
// 	// for _, step := range req.Steps {
// 	// 	experimentStep := entity.ExperimentSteps{
// 	// 		ID:            uuid.NewString(),
// 	// 		ExperimentID:  experiment.ID,
// 	// 		StepOrder:     step.StepOrder,
// 	// 		StepName:      step.StepName,
// 	// 		ResultValue:   step.ResultValue,
// 	// 		StepCondition: step.StepCondition,
// 	// 	}
// 	// 	experimentSteps = append(experimentSteps, experimentStep)
// 	// }
// 	//
// 	// // 获取配方数据
// 	// materialGroups, materials, err := ss.recipeDao.GetMaterialById(ctx, req.RecipeID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// // 生成试验签名
// 	// experimentSignatureMap, _ := generateExperimentSignature([]entity.Experiment{*experiment},
// 	// 	experimentSteps, materialGroups, materials)
// 	// experiment.ExperimentSignature = experimentSignatureMap[experiment.ID]
// 	//
// 	// experimentList, err := ss.experimentDao.GetByExperimentSignature(ctx, experiment.ExperimentSignature)
// 	// if err != nil {
// 	// 	logger.Logger.Errorf("[ExperimentService Add] Mysql err: %v", err)
// 	// 	return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
// 	// }
// 	//
// 	// // 标记是否为自身重复（允许修改的情况）
// 	// isSelfDuplicate := len(experimentList) == 1 && experimentList[0].ID == req.ExperimentID
// 	//
// 	// // 检查查询到的实验列表情况
// 	// if len(experimentList) > 0 && !isSelfDuplicate {
// 	// 	// 存在重复实验（非自身重复），返回业务错误
// 	// 	return nil, utils.NewBusinessError(utils.DuplicateExperimentFormatCode, "", experimentList[0].ExperimentName)
// 	// }
// 	//
// 	// if err = db.Client.Transaction(func(tx *gorm.DB) error {
// 	// 	// 修改实验
// 	// 	update := map[string]interface{}{
// 	// 		entity.ExperimentColumns.ExperimentName:      req.ExperimentName,
// 	// 		entity.ExperimentColumns.Experimenter:        req.Experimenter,
// 	// 		entity.ExperimentColumns.ExperimentSignature: experimentSignatureMap[experiment.ID],
// 	// 		entity.ExperimentColumns.StartTime:           nil,
// 	// 		entity.ExperimentColumns.EndTime:             nil,
// 	// 	}
// 	// 	if req.StartTime != "" {
// 	// 		update[entity.ExperimentColumns.StartTime] = req.StartTime
// 	// 	}
// 	// 	if req.EndTime != "" {
// 	// 		update[entity.ExperimentColumns.EndTime] = req.EndTime
// 	// 	}
// 	// 	if err = ss.experimentDao.Update(ctx, req.ExperimentID, update); err != nil {
// 	// 		logger.Logger.Errorf("[ExperimentService Edit] Update experiment err: %v", err)
// 	// 		return err
// 	// 	}
// 	//
// 	// 	if err = ss.experimentStepDao.DeleteByExperimentIDTx(ctx, tx, req.ExperimentID); err != nil {
// 	// 		logger.Logger.Errorf("[ExperimentService Edit] Delete experimentStep err: %v", err)
// 	// 		return err
// 	// 	}
// 	//
// 	// 	if err = tx.CreateInBatches(&experimentSteps, 500).Error; err != nil {
// 	// 		logger.Logger.Errorf("[ExperimentService Edit] CreateInBatches experimentSteps err: %v", err)
// 	// 		return err
// 	// 	}
// 	//
// 	// 	// 记录操作
// 	// 	if err = ss.recordDao.Insert(ctx, tx, &entity.OperationRecord{
// 	// 		ID:     uuid.NewString(),
// 	// 		OpType: conf.EditExperiment,
// 	// 		UserID: userId,
// 	// 		Desc:   fmt.Sprintf("修改实验: %s", experiment.ExperimentName),
// 	// 	}); err != nil {
// 	// 		logger.Logger.Errorf("[ExperimentService Edit] create record: %v", err)
// 	// 		return err
// 	// 	}
// 	//
// 	// 	return nil
// 	// }); err != nil {
// 	// 	logger.Logger.Errorf("[ExperimentService Edit] Mysql err: %v", err)
// 	// 	return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
// 	// }
//
// 	return nil, nil
// }
