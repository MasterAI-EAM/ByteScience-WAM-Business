package entity

import (
	"time"
)

/******sql******
CREATE TABLE `experiment_steps` (
  `id` varchar(36) NOT NULL COMMENT '实验步骤id',
  `experiment_id` varchar(36) NOT NULL COMMENT '实验ID',
  `recipe_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '配方id',
  `step_order` int NOT NULL COMMENT '步骤顺序',
  `step_name` varchar(255) NOT NULL COMMENT '步骤名称',
  `result_value` varchar(256) DEFAULT NULL COMMENT '步骤结果值',
  `experiment_condition` varchar(255) DEFAULT NULL COMMENT '步骤实验条件',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `experiment_id` (`experiment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='实验步骤表'
******sql******/
// ExperimentSteps 实验步骤表
type ExperimentSteps struct {
	ID                  string    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                               // 实验步骤id
	ExperimentID        string    `gorm:"index:experiment_id;column:experiment_id;type:varchar(36);not null" json:"experimentId"` // 实验ID
	RecipeID            string    `gorm:"column:recipe_id;type:varchar(36);not null" json:"recipeId"`                             // 配方id
	StepOrder           int       `gorm:"column:step_order;type:int;not null" json:"stepOrder"`                                   // 步骤顺序
	StepName            string    `gorm:"column:step_name;type:varchar(255);not null" json:"stepName"`                            // 步骤名称
	ResultValue         string    `gorm:"column:result_value;type:varchar(256);default:null" json:"resultValue"`                  // 步骤结果值
	ExperimentCondition string    `gorm:"column:experiment_condition;type:varchar(255);default:null" json:"experimentCondition"`  // 步骤实验条件
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`    // 创建时间
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`    // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *ExperimentSteps) TableName() string {
	return "experiment_steps"
}

// ExperimentStepsColumns get sql column name.获取数据库列名
var ExperimentStepsColumns = struct {
	ID                  string
	ExperimentID        string
	RecipeID            string
	StepOrder           string
	StepName            string
	ResultValue         string
	ExperimentCondition string
	CreatedAt           string
	UpdatedAt           string
}{
	ID:                  "id",
	ExperimentID:        "experiment_id",
	RecipeID:            "recipe_id",
	StepOrder:           "step_order",
	StepName:            "step_name",
	ResultValue:         "result_value",
	ExperimentCondition: "experiment_condition",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
}
