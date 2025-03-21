package data

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/model/dto"
	"ByteScience-WAM-Business/internal/model/dto/data"
	"ByteScience-WAM-Business/internal/service"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/logger"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TaskApi struct {
	service *service.TaskService
}

// NewTaskApi 创建 TaskApi 实例并初始化依赖项
func NewTaskApi() *TaskApi {
	service := service.NewTaskService()
	return &TaskApi{service: service}
}

// Add 创建任务
// @Summary 上传任务文件
// @Description 接收并存储上传的任务文件，支持多文件上传，限制最大 100MB/文件，最多 20 个文件
// @Tags 任务管理
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "上传的文件，支持多个"
// @Success 200 {object} dto.Empty "文件上传成功"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，可能是文件上传失败或格式不符合要求"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是文件解析或存储失败"
// @Router /data/task [post]
func (api *TaskApi) Add(ctx *gin.Context, req *dto.Empty) (res *dto.Empty, err error) {
	userId, exists := ctx.Get("userId")
	if !exists {
		logger.Logger.Errorf("[Import] User id does not exist")
		return nil, utils.NewBusinessError(utils.UserNotFoundCode, "")
	}

	// 获取表单中的多个文件
	form, err := ctx.MultipartForm()
	if err != nil {
		logger.Logger.Errorf("[Import] ctx.MultipartForm Error: %v", err)
		return nil, utils.NewBusinessError(utils.FileParsingFailedCode, "")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return nil, utils.NewBusinessError(utils.FileParsingFailedCode, ",No files uploaded")
	}

	if len(files) > 20 {
		return nil, utils.NewBusinessError(utils.FileParsingFailedCode, ",The number of files cannot exceed 20")
	}

	var savedFiles []string
	fileHashes := make(map[string]bool) // 用于存储已处理文件的哈希值

	// 遍历上传的文件
	for _, file := range files {
		// 文件不能超出100MB
		if file.Size > 100*1024*1024 {
			logger.Logger.Warnf("File size exceeds 100 MB: %d bytes", file.Size)
			return nil, utils.NewBusinessError(utils.FileTooLargeCode, ",File size exceeds 100 MB")
		}

		// 计算文件哈希
		src, err := file.Open()
		if err != nil {
			src.Close()
			logger.Logger.Errorf("Failed to open file: %v", err)
			return nil, utils.NewBusinessError(utils.FileParsingFailedCode, ",Failed to open file")
		}

		hash := sha256.New()
		if _, err := io.Copy(hash, src); err != nil {
			src.Close()
			logger.Logger.Errorf("Failed to calculate file hash: %v", err)
			return nil, utils.NewBusinessError(utils.FileParsingFailedCode, ",Failed to calculate file hash")
		}
		fileHash := hex.EncodeToString(hash.Sum(nil))
		src.Close()

		// 检查是否已存在相同文件
		if fileHashes[fileHash] {
			logger.Logger.Warnf("Duplicate file detected: %s", file.Filename)
			return nil, utils.NewBusinessError(utils.FileParsingFailedCode, ",Duplicate file detected")
		}
		fileHashes[fileHash] = true // 记录文件哈希值

		fileName := filepath.Base(file.Filename)
		ext := filepath.Ext(fileName) // 获取文件后缀（如 .txt, .jpg）
		if ext == "" {
			ext = "unknown" // 如果没有后缀，设为 "unknown"
		} else {
			ext = ext[1:] // 去掉前面的 "."
		}

		// 目标目录（按文件后缀分类存储）
		dirPath := filepath.Join(conf.GlobalConf.File.TaskPath, ext)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			logger.Logger.Errorf("Failed to create directory: %s, error: %v", dirPath, err)
			return nil, utils.NewBusinessError(utils.FileTooLargeCode, ",Failed to create directory")
		}

		fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName)) // 去除原始文件名的扩展名
		savePath := filepath.Join(dirPath, fmt.Sprintf("%s_%s%s", fileNameWithoutExt, time.Now().Format("20060102150405"), filepath.Ext(fileName)))

		// 保存文件到本地
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			logger.Logger.Errorf("File saving failed: %v", err)
			// 发生错误时删除已保存的文件，保证要么全部成功要么一个不保存
			for _, savedFile := range savedFiles {
				_ = os.Remove(savedFile)
			}
			return nil, utils.NewBusinessError(utils.FileTooLargeCode, ",File saving failed")
		}

		savedFiles = append(savedFiles, savePath)
	}

	res, err = api.service.Add(ctx, userId.(string), savedFiles)
	if err != nil {
		// 发生错误时删除已保存的文件，保证要么全部成功要么一个不保存
		for _, savedFile := range savedFiles {
			if err := os.Remove(savedFile); err != nil {
				logger.Logger.Errorf("Failed to delete file: %s, error: %v", savedFile, err)
			}
		}
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}

	return
}

// List 获取任务列表
// @Summary 获取任务列表
// @Description 支持分页查询任务列表，可按任务名称筛选
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param req body data.TaskListRequest true "分页参数及筛选条件"
// @Success 200 {object} data.TaskListResponse "成功返回任务列表"
// @Failure 400 {object} dto.ErrorResponse "请求参数错误，如分页参数错误、筛选条件不符合要求"
// @Failure 500 {object} dto.ErrorResponse "服务器内部错误，可能是数据库查询失败"
// @Router /data/task [get]
func (api *TaskApi) List(ctx *gin.Context, req *data.TaskListRequest) (res *data.TaskListResponse, err error) {
	res, err = api.service.List(ctx, req)
	if err != nil {
		logger.Logger.Errorf("[List] List Error: %v", err)
		if businessErr, ok := err.(*utils.BusinessError); ok {
			return nil, businessErr
		}
		return nil, utils.NewBusinessError(utils.InternalError, "")
	}
	return
}
