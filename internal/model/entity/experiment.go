package entity

import (
	"time"
)

/******sql******
CREATE TABLE `experiment` (
  `id` varchar(36) NOT NULL COMMENT '实验id',
  `experiment_signature` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配方唯一签名（哈希值）',
  `file_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '文件id',
  `experiment_name` varchar(255) NOT NULL COMMENT '实验名称',
  `entry_category` tinyint NOT NULL DEFAULT '1' COMMENT '1 文件导入 2 页面输入',
  `sort` int NOT NULL COMMENT '排序',
  `experimenter` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '实验者',
  `user_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作用户',
  `start_time` datetime DEFAULT NULL COMMENT '实验开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '实验结束时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_signature` (`experiment_signature`),
  KEY `file_id` (`file_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='实验表'
******sql******/
// Experiment 实验表
type Experiment struct {
	ID                  string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                                // 实验id
	ExperimentSignature string     `gorm:"unique;column:experiment_signature;type:varchar(36);not null" json:"experimentSignature"` // 配方唯一签名（哈希值）
	FileID              string     `gorm:"index:file_id;column:file_id;type:varchar(36);default:null" json:"fileId"`                // 文件id
	ExperimentName      string     `gorm:"column:experiment_name;type:varchar(255);not null" json:"experimentName"`                 // 实验名称
	EntryCategory       int8       `gorm:"column:entry_category;type:tinyint;not null;default:1" json:"entryCategory"`              // 1 文件导入 2 页面输入
	Sort                int        `gorm:"column:sort;type:int;not null" json:"sort"`                                               // 排序
	Experimenter        string     `gorm:"column:experimenter;type:varchar(128);not null" json:"experimenter"`                      // 实验者
	UserID              string     `gorm:"index:user_id;column:user_id;type:varchar(36);not null" json:"userId"`                    // 操作用户
	StartTime           *time.Time `gorm:"column:start_time;type:datetime;default:null" json:"startTime"`                           // 实验开始时间
	EndTime             *time.Time `gorm:"column:end_time;type:datetime;default:null" json:"endTime"`                               // 实验结束时间
	CreatedAt           time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`     // 创建时间
	UpdatedAt           time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`     // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *Experiment) TableName() string {
	return "experiment"
}

// ExperimentColumns get sql column name.获取数据库列名
var ExperimentColumns = struct {
	ID                  string
	ExperimentSignature string
	FileID              string
	ExperimentName      string
	EntryCategory       string
	Sort                string
	Experimenter        string
	UserID              string
	StartTime           string
	EndTime             string
	CreatedAt           string
	UpdatedAt           string
}{
	ID:                  "id",
	ExperimentSignature: "experiment_signature",
	FileID:              "file_id",
	ExperimentName:      "experiment_name",
	EntryCategory:       "entry_category",
	Sort:                "sort",
	Experimenter:        "experimenter",
	UserID:              "user_id",
	StartTime:           "start_time",
	EndTime:             "end_time",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}
