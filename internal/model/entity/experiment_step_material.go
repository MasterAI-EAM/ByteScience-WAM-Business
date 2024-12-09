package entity

/******sql******
CREATE TABLE `experiment_step_material` (
  `experiment_step_id` varchar(36) NOT NULL COMMENT '实验步骤ID',
  `experiment_material_group_id` varchar(36) NOT NULL COMMENT '实验材料组ID',
  `proportion` decimal(10,2) NOT NULL COMMENT '实验材料组在步骤中的占比（%）',
  PRIMARY KEY (`experiment_step_id`,`experiment_material_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='实验步骤与材料组使用表'
******sql******/
// ExperimentStepMaterial 实验步骤与材料组使用表
type ExperimentStepMaterial struct {
	ExperimentStepID          string  `gorm:"primaryKey;column:experiment_step_id;type:varchar(36);not null" json:"experimentStepId"`                    // 实验步骤ID
	ExperimentMaterialGroupID string  `gorm:"primaryKey;column:experiment_material_group_id;type:varchar(36);not null" json:"experimentMaterialGroupId"` // 实验材料组ID
	Proportion                float64 `gorm:"column:proportion;type:decimal(10,2);not null" json:"proportion"`                                           // 实验材料组在步骤中的占比（%）
}

// TableName get sql table name.获取数据库表名
func (m *ExperimentStepMaterial) TableName() string {
	return "experiment_step_material"
}

// ExperimentStepMaterialColumns get sql column name.获取数据库列名
var ExperimentStepMaterialColumns = struct {
	ExperimentStepID          string
	ExperimentMaterialGroupID string
	Proportion                string
}{
	ExperimentStepID:          "experiment_step_id",
	ExperimentMaterialGroupID: "experiment_material_group_id",
	Proportion:                "proportion",
}
