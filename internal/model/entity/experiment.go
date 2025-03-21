package entity

import (
	"time"
)

/******sql******
CREATE TABLE `experiment` (
  `id` varchar(36) NOT NULL COMMENT '实验id',
  `signature` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配方唯一签名（哈希值）',
  `task_id` varchar(36) NOT NULL COMMENT '任务id',
  `experiment_name` varchar(255) NOT NULL COMMENT '实验名称',
  `bs_type` enum('high_performance','wind_power') NOT NULL DEFAULT 'high_performance' COMMENT '业务类型（high_performance=高功能, wind_power=风电）',
  `entry_category` enum('file_import','manual_entry') NOT NULL DEFAULT 'file_import' COMMENT '录入类别（file_import=文件导入, manual_entry=页面录入）',
  `experimenter` varchar(128) NOT NULL COMMENT '实验者',
  `sort` int NOT NULL AUTO_INCREMENT COMMENT '排序字段（自增）',
  `user_id` varchar(36) NOT NULL COMMENT '用户id',
  `start_time` datetime DEFAULT NULL COMMENT '实验开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '实验结束时间',
  `status` enum('pending_review','approved','rejected') NOT NULL COMMENT '实验状态（pending_review=待审核, approved=审核通过, rejected=审核不通过）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `sort` (`sort`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='实验表'
******sql******/
// Experiment 实验表
type Experiment struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null"`                                              // 实验id
	Signature      string     `gorm:"column:signature;type:varchar(36);not null"`                                                  // 配方唯一签名（哈希值）
	TaskID         string     `gorm:"column:task_id;type:varchar(36);not null"`                                                    // 任务id
	ExperimentName string     `gorm:"column:experiment_name;type:varchar(255);not null"`                                           // 实验名称
	BsType         string     `gorm:"column:bs_type;type:enum('high_performance','wind_power');not null;default:high_performance"` // 业务类型（high_performance=高功能, wind_power=风电）
	EntryCategory  string     `gorm:"column:entry_category;type:enum('file_import','manual_entry');not null;default:file_import"`  // 录入类别（file_import=文件导入, manual_entry=页面录入）
	Experimenter   string     `gorm:"column:experimenter;type:varchar(128);not null"`                                              // 实验者
	Sort           int        `gorm:"unique;column:sort;type:int;not null"`                                                        // 排序字段（自增）
	UserID         string     `gorm:"index:idx_user_id;column:user_id;type:varchar(36);not null"`                                  // 用户id
	StartTime      *time.Time `gorm:"column:start_time;type:datetime"`                                                             // 实验开始时间
	EndTime        *time.Time `gorm:"column:end_time;type:datetime"`                                                               // 实验结束时间
	Status         string     `gorm:"column:status;type:enum('pending_review','approved','rejected');not null"`                    // 实验状态（pending_review=待审核, approved=审核通过, rejected=审核不通过）
	CreatedAt      time.Time  `gorm:"index:idx_created_at;column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`     // 创建时间
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`                          // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *Experiment) TableName() string {
	return "experiment"
}

// ExperimentColumns get sql column name.获取数据库列名
var ExperimentColumns = struct {
	ID             string
	Signature      string
	TaskID         string
	ExperimentName string
	BsType         string
	EntryCategory  string
	Experimenter   string
	Sort           string
	UserID         string
	StartTime      string
	EndTime        string
	Status         string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "id",
	Signature:      "signature",
	TaskID:         "task_id",
	ExperimentName: "experiment_name",
	BsType:         "bs_type",
	EntryCategory:  "entry_category",
	Experimenter:   "experimenter",
	Sort:           "sort",
	UserID:         "user_id",
	StartTime:      "start_time",
	EndTime:        "end_time",
	Status:         "status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}
