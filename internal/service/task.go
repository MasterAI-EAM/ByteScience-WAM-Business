package service

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/dao"
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/data"
	"ByteScience-WAM-Business/internal/model/entity"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/logger"
	"context"
	"github.com/google/uuid"
	"path"
	"strings"
	"time"
)

type TaskService struct {
	dao *dao.TaskDao
}

// NewTaskService 创建一个新的 TaskService 实例
func NewTaskService() *TaskService {
	return &TaskService{
		dao: dao.NewTaskDao(),
	}
}

// Add 创建任务
func (ts TaskService) Add(ctx context.Context, userId string, uploadedFiles []string) (*dto.Empty, error) {
	// 生成批次号
	batchId := uuid.NewString()
	tasks := make([]*entity.Task, 0, len(uploadedFiles))

	// 遍历上传的文件，构建 Task 结构体
	for _, filePath := range uploadedFiles {
		path, name := extractFileDetails(filePath, conf.GlobalConf.File.TaskPath)
		task := &entity.Task{
			ID:           uuid.NewString(),
			BatchID:      batchId,
			UserID:       userId,
			FileName:     name,
			FilePath:     path,
			JSONFilePath: "",
			AiFilePath:   "",
			Status:       dao.TaskStatusPending,
			Remark:       "",
		}
		tasks = append(tasks, task)
	}

	// 插入数据库
	err := ts.dao.BatchInsert(ctx, nil, tasks)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// List 查询任务及其相关数据，包括任务步骤和材料组
func (ts *TaskService) List(ctx context.Context, req *data.TaskListRequest) (*data.TaskListResponse, error) {
	res := &data.TaskListResponse{
		Total: 0,
		List:  make([]data.TaskData, 0),
	}

	// 查询实验及其基础信息
	tasks, total, err := ts.dao.Query(ctx, req.Page, req.PageSize, map[string]interface{}{
		entity.TaskColumns.Status:   req.Status,
		entity.TaskColumns.FileName: req.FileName,
		"startTime":                 req.StartTime,
		"endTIme":                   req.EndTime,
	})
	if err != nil {
		logger.Logger.Errorf("[TaskService List] Mysql err: %v", err)
		return nil, utils.NewBusinessError(utils.DatabaseErrorCode, "")
	}

	// 构造返回的任务列表数据
	taskDataList := make([]data.TaskData, len(tasks))
	for i, task := range tasks {
		// 直接将任务映射为 TaskData 结构体
		taskDataList[i] = data.TaskData{
			TaskID:       task.ID,
			BatchID:      task.BatchID,
			FileName:     task.FileName,
			FilePath:     task.FilePath,
			JSONFilePath: task.JSONFilePath,
			AiFilePath:   task.AiFilePath,
			Status:       task.Status,
			Remark:       task.Remark,
			CreatedAt:    task.CreatedAt.Format(time.DateTime),
			UpdatedAt:    task.UpdatedAt.Format(time.DateTime),
		}
	}

	res.Total = total
	res.List = taskDataList

	return res, nil
}

// 提取文件路径和文件名（带时间戳的文件名）
func extractFileDetails(basePath string, dir string) (string, string) {
	// 从 basePath 中提取相对路径部分
	relativePath := strings.TrimPrefix(basePath, dir+"/")

	// 提取文件路径（去掉文件名部分）
	filePath := path.Dir(relativePath)

	// 提取文件名（带时间戳）
	fileNameWithTimestamp := path.Base(relativePath)

	// 去掉时间戳部分，保留文件名
	fileName := strings.TrimSuffix(fileNameWithTimestamp, path.Ext(fileNameWithTimestamp))

	return filePath + "/" + fileNameWithTimestamp, fileName + path.Ext(fileNameWithTimestamp)
}
