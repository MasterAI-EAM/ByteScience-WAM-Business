package entity

import (
	"time"
)

/******sql******
CREATE TABLE `experiment_files` (
  `id` varchar(36) NOT NULL COMMENT '文件id',
  `file_name` varchar(255) NOT NULL COMMENT '文件名称',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='实验来源文件表'
******sql******/
// ExperimentFiles 实验来源文件表
type ExperimentFiles struct {
	ID        string    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                            // 文件id
	FileName  string    `gorm:"column:file_name;type:varchar(255);not null" json:"fileName"`                         // 文件名称
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"` // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *ExperimentFiles) TableName() string {
	return "experiment_files"
}

// ExperimentFilesColumns get sql column name.获取数据库列名
var ExperimentFilesColumns = struct {
	ID        string
	FileName  string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	FileName:  "file_name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
