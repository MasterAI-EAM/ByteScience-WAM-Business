package entity

import (
	"time"
)

/******sql******
CREATE TABLE `task` (
  `id` varchar(36) NOT NULL COMMENT '文件id',
  `batch_id` varchar(36) NOT NULL COMMENT '批次号',
  `user_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户id',
  `file_name` varchar(255) NOT NULL COMMENT '文件名称',
  `file_path` varchar(255) NOT NULL COMMENT '文件路径',
  `json_file_path` varchar(255) NOT NULL COMMENT '硬代码json文件路径',
  `ai_file_path` varchar(255) NOT NULL COMMENT 'ai处理后json文件路径',
  `status` enum('pending','processing','success','failure') NOT NULL COMMENT '任务状态（pending=待处理, processing=处理中, success=成功, failure=失败）',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '任务状态描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_batch_status` (`batch_id`,`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='任务表'
******sql******/
// Task 任务表
type Task struct {
	ID           string    `gorm:"primaryKey;column:id;type:varchar(36);not null"`                                                                       // 文件id
	BatchID      string    `gorm:"index:idx_batch_status;column:batch_id;type:varchar(36);not null"`                                                     // 批次号
	UserID       string    `gorm:"column:user_id;type:varchar(36);not null"`                                                                             // 用户id
	FileName     string    `gorm:"column:file_name;type:varchar(255);not null"`                                                                          // 文件名称
	FilePath     string    `gorm:"column:file_path;type:varchar(255);not null"`                                                                          // 文件路径
	JSONFilePath string    `gorm:"column:json_file_path;type:varchar(255);not null"`                                                                     // 硬代码json文件路径
	AiFilePath   string    `gorm:"column:ai_file_path;type:varchar(255);not null"`                                                                       // ai处理后json文件路径
	Status       string    `gorm:"index:idx_status;index:idx_batch_status;column:status;type:enum('pending','processing','success','failure');not null"` // 任务状态（pending=待处理, processing=处理中, success=成功, failure=失败）
	Remark       string    `gorm:"column:remark;type:varchar(512);not null;default:''"`                                                                  // 任务状态描述
	CreatedAt    time.Time `gorm:"index:idx_created_at;column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`                              // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`                                                   // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *Task) TableName() string {
	return "task"
}

// TaskColumns get sql column name.获取数据库列名
var TaskColumns = struct {
	ID           string
	BatchID      string
	UserID       string
	FileName     string
	FilePath     string
	JSONFilePath string
	AiFilePath   string
	Status       string
	Remark       string
	CreatedAt    string
	UpdatedAt    string
}{
	ID:           "id",
	BatchID:      "batch_id",
	UserID:       "user_id",
	FileName:     "file_name",
	FilePath:     "file_path",
	JSONFilePath: "json_file_path",
	AiFilePath:   "ai_file_path",
	Status:       "status",
	Remark:       "remark",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}
