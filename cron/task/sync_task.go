package task

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/dao"
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/db"
	"ByteScience-WAM-Business/pkg/gpt"
	"ByteScience-WAM-Business/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ParsedExperimentDataForever 消费文件解析入库任务
func ParsedExperimentDataForever() {
	ctx := context.Background()
	task := dao.NewTaskDao()

	// 初始休眠时间
	sleepDuration := 1 * time.Second
	maxSleepDuration := 2 * time.Minute

	for {
		info, err := task.GetPendingData(ctx)
		if err != nil {
			logger.Logger.Error("[ParsedExperimentDataForever] GetPendingData err:", err)
			sleepDuration = ExponentialBackoffSleep(sleepDuration, maxSleepDuration)
			continue
		}

		if info == nil {
			sleepDuration = ExponentialBackoffSleep(sleepDuration, maxSleepDuration)
			continue
		}

		// **有任务时，重置休眠时间**
		sleepDuration = 1 * time.Second

		// **执行业务逻辑**
		err = processExperimentData(ctx, info)
		if err != nil {
			logger.Logger.Errorf("[ParsedExperimentDataForever] 任务 ID: %s 处理失败: %v", info.ID, err)
		}
	}
}

// processExperimentData 业务逻辑
func processExperimentData(ctx context.Context, info *entity.Task) (err error) {
	task := dao.NewTaskDao()

	// 确保最终更新任务状态
	defer func() {
		status := dao.TaskStatusSuccess
		remark := ""
		if err != nil {
			status = dao.TaskStatusFailure
			remark = err.Error()
		}

		updateErr := task.Update(ctx, info.ID, map[string]interface{}{
			entity.TaskColumns.Status: status,
			entity.TaskColumns.Remark: remark,
		}, nil)
		if updateErr != nil {
			logger.Logger.Error("[processExperimentData] task.Update failure err:", updateErr)
		}
	}()

	// 开始处理任务
	err = task.Update(ctx, info.ID, map[string]interface{}{
		entity.TaskColumns.Status: dao.TaskStatusProcessing,
	}, nil)
	if err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	openPath := filepath.Join(conf.GlobalConf.File.TaskPath, info.FilePath)
	fileName := filepath.Base(info.FilePath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	beforePath := "json/before_" + fileNameWithoutExt + ".json"
	beforeSavePath := filepath.Join(conf.GlobalConf.File.TaskPath, "/", beforePath)

	// Excel 解析
	beforeData, err := utils.ExcelToJson(openPath, beforeSavePath)
	if err != nil {
		return fmt.Errorf("excel 解析失败: %v", err)
	}

	// GPT 处理
	prefix := "作为一个专业的材料实验数据分析助手，你的任务是从提供的JSON格式实验报告中提取关键信息。本实验涉及树脂和固化剂的制备与测试。实验过程包括：首先单独制备并测试不同组分的树脂和固化剂，然后按特定比例混合这些树脂(A剂)和固化剂(B剂)，形成混合物，并对其进行进一步的测试。您需要提取：\n基本信息：树脂编号，树脂配方，固化剂编号，固化剂配方，试验组的配比。\n树脂和固化剂的基本测试：环氧当量，胺值，不同温度下的粘度。\n混合物的基本测试：不同温度下的混合粘度，不同温度下的可使用时间\n混合物的力学性能测试：拉伸强度，拉伸模量，断裂延伸率，弯曲强度，弯曲模量，压缩强度，冲击韧性，Tg。\n\n提示：\n1. 出于涉密保护，数据的成分或者测试的名称可能被编号替代，提取后的JSON数据值也可以直接使用编号。\n2. 每一条JSON数据为一组试验数据，值列表长度一致。\n3. 每个值列表中且相同索引处的值在源文件中位于同一行。\n4. JSON文件由'.xlsx'表格文件转换得到。\n5. 当某个配方的所有组份都是空时，有可能是这个编号的配方已经出现过，但是也存在配方组份未记录的情况。\n6. JSON数据包含表头数据和正文数据，表头数据内不可能包含需要被提取的数据。\n7. 表头数据内出现的纯数字是干扰项，请忽略。\n\n处理要求：\n1. 合并重复出现的相同内容，表格中跨单元格的合并内容需完整提取。\n2. 配方数据要保留原始编号格式。\n3. 遇到缺失数据标注为\"null\"。\n4. 仅提取原数据中存在的数据，不要编造数据，不要发生数据串行的情况。\n5. 以JSON格式输出，输出仅为一段JSON数据，不包含任何其他内容，比如分析的过程。\n6. 输出格式和示例如下，除了树脂配方和固化剂配方下的值字典中的键名可变，其他所有的键名都不可变：\n{\"树脂编号\":\"M5\",\"树脂配方\":{\"YYYA1\":0,\"YYYA2\":88.4,\"YYYA3\":0,\"YYYA4\":10,\"YYYA5\":0,\"YYYA6\":0.1,\"YYYA7\":1.5},\"树脂混合物测试项\":{\"环氧当量\":0.5547609148,\"树脂粘度\":{\"温度\":27,\"粘度\":0.5547609148}},\"固化剂编号\":\"B-M5\",\"固化剂配方\":{\"YYYB1\":65,\"YYYB2\":0,\"YYYB3\":35,\"YYYB4\":0},\"固化剂混合物测试项\":{\"胺值\":1.938657407,\"固化剂粘度\":{\"温度\":27,\"粘度\":0.5547609148}},\"配比\":\"100/30\",\"树脂/固化剂混合物测试项\":{\"可用时间\":94,\"混合粘度\":320,\"温度\":\"27℃\"},\"固化温度与时长\":\"70℃*1h+100℃*1h\",\"力学性能测试\":{\"拉伸强度\":71,\"拉伸模量\":3222,\"断裂延伸率\":7.8,\"弯曲强度\":\"\",\"弯曲模量\":\"\",\"压缩强度\":\"\",\"冲击韧性\":\"\",\"Tg\":90}}\n现在请分析以下实验数据："
	afterData, err := gpt.ExtractInformationWithGPT4(beforeData, prefix, 2)
	if err != nil {
		return fmt.Errorf("GPT 解析失败: %v", err)
	}

	// 结果存储
	afterPath := "json/after_" + fileNameWithoutExt + ".json"
	afterSavePath := filepath.Join(conf.GlobalConf.File.TaskPath, "/", afterPath)
	err = utils.WriteJSONToFile(afterData, afterSavePath)
	if err != nil {
		return fmt.Errorf("写入 JSON 失败: %v", err)
	}

	// 正确时更新 JSON 路径
	err = task.Update(ctx, info.ID, map[string]interface{}{
		entity.TaskColumns.JSONFilePath: beforePath,
		entity.TaskColumns.AiFilePath:   afterPath,
	}, nil)
	if err != nil {
		return fmt.Errorf("更新 JSON 路径失败: %v", err)
	}

	// 解析json
	experiments, experimentSteps, materialGroups, materials, err := ExperimentDataToSqlData(afterData, info)
	if err != nil {
		return fmt.Errorf("JSON 转换成 sql 失败: %v", err)
	}

	for i, _ := range experiments {
		experiments[i].ExperimentName = fmt.Sprintf("%s_%d", strings.Split(fileNameWithoutExt, "_")[0], i+1)
	}

	if err := db.Client.Transaction(func(tx *gorm.DB) error {
		insertNum := 500
		if err := tx.Omit(entity.ExperimentColumns.Sort).CreateInBatches(experiments, insertNum).Error; err != nil {
			logger.Logger.Errorf("[processExperimentData] CreateInBatches experiments err: %v", err)
			return err
		}
		if err := tx.CreateInBatches(experimentSteps, insertNum).Error; err != nil {
			logger.Logger.Errorf("[processExperimentData] CreateInBatches experiments err: %v", err)
			return err
		}
		if len(materialGroups) > 0 {
			if err := tx.CreateInBatches(materialGroups, insertNum).Error; err != nil {
				logger.Logger.Errorf("[processExperimentData] CreateInBatches recipeMaterialGroups err: %v", err)
				return err
			}
		}
		if err := tx.CreateInBatches(materials, insertNum).Error; err != nil {
			logger.Logger.Errorf("[processExperimentData] CreateInBatches materials err: %v", err)
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Errorf("[processExperimentData] Mysql err: %v", err)
		return utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	return nil
}

// ExperimentDataToSqlData 解析数据并返回对应的实体对象
func ExperimentDataToSqlData(data []map[string]interface{}, info *entity.Task) ([]*entity.Experiment,
	[]*entity.ExperimentSteps,
	[]*entity.MaterialGroups, []*entity.Materials, error) {
	var experiments []*entity.Experiment
	var experimentSteps []*entity.ExperimentSteps
	var materialGroups []*entity.MaterialGroups
	var materials []*entity.Materials

	for _, row := range data {
		experiment := &entity.Experiment{
			ID:            uuid.NewString(),
			Status:        "pending_review",
			EntryCategory: "file_import",
			UserID:        info.UserID,
			TaskID:        info.ID,
		}
		materialGroupsResin := &entity.MaterialGroups{
			ID:                    uuid.NewString(),
			ExperimentID:          experiment.ID,
			MaterialGroupCategory: "resin",
			Sort:                  1,
			UserID:                info.UserID,
		}
		materialGroupsHardener := &entity.MaterialGroups{
			ID:                    uuid.NewString(),
			ExperimentID:          experiment.ID,
			MaterialGroupCategory: "hardener",
			Sort:                  2,
			UserID:                info.UserID,
		}
		stepCondition := ""
		if val, ok := row["固化温度与时长"].(string); ok {
			stepCondition = val
		}
		for key, value := range row {
			switch key {
			case "树脂编号":
				materialGroupsResin.MaterialGroupName = value.(string)
			case "固化剂编号":
				materialGroupsHardener.MaterialGroupName = value.(string)
			case "配比":
				resinProportion, hardenerProportion, err := convertRatioToProportion(value.(string))
				if err != nil {
					return nil, nil, nil, nil, fmt.Errorf("convertRatioToProportion error: %v", err)
				}
				materialGroupsResin.Proportion, materialGroupsHardener.Proportion = resinProportion, hardenerProportion
			case "树脂混合物测试项", "固化剂混合物测试项", "树脂/固化剂混合物测试项":
				jData, _ := json.Marshal(value)
				experimentSteps = append(experimentSteps, &entity.ExperimentSteps{
					ID:           uuid.NewString(),
					ExperimentID: experiment.ID,
					StepCategory: getStepCategoryByKey(key),
					StepName:     key,
					Sort:         getSort(key),
					ResultValue:  datatypes.JSON(jData),
					UserID:       info.UserID,
				})
			case "力学性能测试":
				sort := getSort(key)
				i := 0
				for k, v := range value.(map[string]interface{}) {
					if v != nil && v != "null" && v != "" && v != 0 {
						jData, _ := json.Marshal(map[string]interface{}{"value": v})
						experimentSteps = append(experimentSteps, &entity.ExperimentSteps{
							ID:            uuid.NewString(),
							ExperimentID:  experiment.ID,
							StepCategory:  "mechanical_performance",
							StepName:      k,
							ResultValue:   datatypes.JSON(jData),
							Sort:          sort + i,
							StepCondition: &stepCondition,
							UserID:        info.UserID,
						})
					}
					i++
				}
			case "树脂配方", "固化剂配方":
				sort := getSort(key)
				groupID := materialGroupsResin.ID
				if key == "固化剂配方" {
					groupID = materialGroupsHardener.ID
				}
				i := 0
				for k, v := range value.(map[string]interface{}) {
					if v == 0 || v == "" || v == "null" || v == nil {
						continue
					}
					if proportion, ok := v.(float64); ok && proportion > 0 {
						materials = append(materials, &entity.Materials{
							ID:              uuid.NewString(),
							MaterialName:    k,
							Sort:            sort + i,
							MaterialGroupID: groupID,
							Proportion:      proportion,
							UserID:          info.UserID,
						})
					}
					i++
				}
			}
		}

		experiments = append(experiments, experiment)
		materialGroups = append(materialGroups, materialGroupsResin, materialGroupsHardener)
	}

	return experiments, experimentSteps, materialGroups, materials, nil
}

// getStepCategoryByKey 根据键获取实验步骤分类
func getStepCategoryByKey(key string) string {
	switch key {
	case "树脂混合物测试项":
		return "resin_mixing"
	case "固化剂混合物测试项":
		return "hardener_mixing"
	case "树脂/固化剂混合物测试项":
		return "resin_hardener_mixing"
	}
	return ""
}

// getSort 获取排序值
func getSort(key string) int {
	switch key {
	case "树脂混合物测试项":
		return 1
	case "固化剂混合物测试项":
		return 2
	case "树脂/固化剂混合物测试项":
		return 3
	case "力学性能测试":
		return 4
	case "resin":
		return 1
	case "hardener":
		return 2
	case "树脂配方":
		return 1
	case "固化剂配方":
		return 100
	}
	return 0
}

// convertRatioToProportion 将"100/30"格式的字符串转换为百分比小数
func convertRatioToProportion(input string) (float64, float64, error) {
	if input == "" {
		return 0, 0, nil
	}

	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("输入格式错误: 需要形如 '100/30' 的字符串")
	}

	numerator, err1 := strconv.ParseFloat(parts[0], 64)
	denominator, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("数值转换失败: %v, %v", err1, err2)
	}

	// **优化**：避免分母为 0
	total := numerator + denominator
	if total == 0 {
		return 0, 0, fmt.Errorf("无效比例: 100/0")
	}

	return numerator / total * 100, denominator / total * 100, nil
}

// generateExperimentSignature 生成实验密钥
func generateExperimentSignature(experiment *entity.Experiment, experimentSteps []*entity.ExperimentSteps,
	materialGroups []*entity.MaterialGroups, materials []*entity.Materials) string {

	// 建立实验步骤映射
	stepMap := make(map[string][]*entity.ExperimentSteps)
	for _, step := range experimentSteps {
		stepMap[step.ExperimentID] = append(stepMap[step.ExperimentID], step)
	}

	// 建立实验-材料组映射
	materialGroupsMap := make(map[string][]*entity.MaterialGroups)
	for _, materialGroup := range materialGroups {
		materialGroupsMap[materialGroup.ExperimentID] = append(materialGroupsMap[materialGroup.ExperimentID], materialGroup)
	}

	// 建立材料组-材料映射
	materialMap := make(map[string][]*entity.Materials)
	for _, material := range materials {
		materialMap[material.MaterialGroupID] = append(materialMap[material.MaterialGroupID], material)
	}

	var signatureElements []string

	// 处理实验步骤
	if steps, exists := stepMap[experiment.ID]; exists {
		sort.SliceStable(steps, func(i, j int) bool {
			return strings.ToLower(strings.TrimSpace(steps[i].StepName)) <
				strings.ToLower(strings.TrimSpace(steps[j].StepName))
		})
		for _, step := range steps {
			signatureElements = append(signatureElements, step.StepName, step.ResultValue.String(), *step.StepCondition)
		}
	}

	// 处理材料组
	if materialGroups, exists := materialGroupsMap[experiment.ID]; exists {
		sort.SliceStable(materialGroups, func(i, j int) bool {
			return strings.ToLower(strings.TrimSpace(materialGroups[i].MaterialGroupName)) <
				strings.ToLower(strings.TrimSpace(materialGroups[j].MaterialGroupName))
		})
		for _, materialGroup := range materialGroups {
			signatureElements = append(signatureElements, fmt.Sprintf("%s%f",
				materialGroup.MaterialGroupName,
				materialGroup.Proportion))

			// 处理材料
			if materials, exists := materialMap[materialGroup.ID]; exists {
				sort.SliceStable(materials, func(i, j int) bool {
					return strings.ToLower(strings.TrimSpace(materials[i].MaterialName)) <
						strings.ToLower(strings.TrimSpace(materials[j].MaterialName))
				})
				for _, material := range materials {
					signatureElements = append(signatureElements, fmt.Sprintf("%s%f",
						material.MaterialName,
						material.Proportion))
				}
			}
		}
	}

	// 计算 MD5 哈希
	signature := utils.Md5Hash(strings.Join(signatureElements, "|"))

	return signature
}
