package entity

import (
	"time"
)

/******sql******
CREATE TABLE `experiment` (
  `id` varchar(36) NOT NULL COMMENT '实验id',
  `file_id` varchar(36) NOT NULL COMMENT '文件id',
  `experiment_name` varchar(255) NOT NULL COMMENT '实验名称',
  `sort` bigint NOT NULL COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `file_id` (`file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='实验表'
******sql******/
// Experiment 实验表
type Experiment struct {
	ID             string    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                            // 实验id
	FileID         string    `gorm:"index:file_id;column:file_id;type:varchar(36);not null" json:"fileId"`                // 文件id
	ExperimentName string    `gorm:"column:experiment_name;type:varchar(255);not null" json:"experimentName"`             // 实验名称
	Sort           int64     `gorm:"column:sort;type:bigint;not null" json:"sort"`                                        // 排序
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"` // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"` // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *Experiment) TableName() string {
	return "experiment"
}

// ExperimentColumns get sql column name.获取数据库列名
var ExperimentColumns = struct {
	ID             string
	FileID         string
	ExperimentName string
	Sort           string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "id",
	FileID:         "file_id",
	ExperimentName: "experiment_name",
	Sort:           "sort",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}
