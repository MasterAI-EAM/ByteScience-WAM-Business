package entity

import (
	"gorm.io/datatypes"
	"time"
)

/******sql******
CREATE TABLE `experiment_steps` (
  `id` varchar(36) NOT NULL COMMENT '实验步骤id',
  `experiment_id` varchar(36) NOT NULL COMMENT '实验id',
  `step_category` enum('resin_mixing','hardener_mixing','resin_hardener_mixing','mechanical_performance') NOT NULL DEFAULT 'resin_mixing' COMMENT '实验步骤类别（resin_mixing=树脂混合, hardener_mixing=固化剂混合, resin_hardener_mixing=树脂/固化剂混合, mechanical_performance=力学性能）',
  `step_name` varchar(255) NOT NULL COMMENT '步骤名称',
  `sort` int NOT NULL,
  `result_value` json DEFAULT NULL COMMENT '步骤结果值，JSON 格式存储多个测试项',
  `step_condition` varchar(255) DEFAULT '' COMMENT '步骤实验条件',
  `user_id` varchar(36) NOT NULL COMMENT '用户id',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_experiment_id` (`experiment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='实验步骤表'
******sql******/
// ExperimentSteps 实验步骤表
type ExperimentSteps struct {
	ID            string         `gorm:"primaryKey;column:id;type:varchar(36);not null"`                                                                                                  // 实验步骤id
	ExperimentID  string         `gorm:"index:idx_experiment_id;column:experiment_id;type:varchar(36);not null"`                                                                          // 实验id
	StepCategory  string         `gorm:"column:step_category;type:enum('resin_mixing','hardener_mixing','resin_hardener_mixing','mechanical_performance');not null;default:resin_mixing"` // 实验步骤类别（resin_mixing=树脂混合, hardener_mixing=固化剂混合, resin_hardener_mixing=树脂/固化剂混合, mechanical_performance=力学性能）
	StepName      string         `gorm:"column:step_name;type:varchar(255);not null"`                                                                                                     // 步骤名称
	Sort          int            `gorm:"column:sort;type:int;not null"`
	ResultValue   datatypes.JSON `gorm:"column:result_value;type:json"`                                      // 步骤结果值，JSON 格式存储多个测试项
	StepCondition *string        `gorm:"column:step_condition;type:varchar(255);default:''"`                 // 步骤实验条件
	UserID        string         `gorm:"column:user_id;type:varchar(36);not null"`                           // 用户id
	CreatedAt     time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"` // 创建时间
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP"` // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *ExperimentSteps) TableName() string {
	return "experiment_steps"
}

// ExperimentStepsColumns get sql column name.获取数据库列名
var ExperimentStepsColumns = struct {
	ID            string
	ExperimentID  string
	StepCategory  string
	StepName      string
	Sort          string
	ResultValue   string
	StepCondition string
	UserID        string
	CreatedAt     string
	UpdatedAt     string
}{
	ID:            "id",
	ExperimentID:  "experiment_id",
	StepCategory:  "step_category",
	StepName:      "step_name",
	Sort:          "sort",
	ResultValue:   "result_value",
	StepCondition: "step_condition",
	UserID:        "user_id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}
