package entity

import (
	"time"
)

/******sql******
CREATE TABLE `materials` (
  `id` varchar(36) NOT NULL COMMENT '材料id',
  `material_name` varchar(255) NOT NULL COMMENT '材料名称',
  `experiment_material_group_id` varchar(36) NOT NULL COMMENT '实验材料组id',
  `material_group_name` varchar(255) NOT NULL COMMENT '材料组名称',
  `percentage` decimal(10,2) NOT NULL COMMENT '材料在组内的占比（%）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `experiment_material_group_id` (`experiment_material_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='材料表'
******sql******/
// Materials 材料表
type Materials struct {
	ID                        string    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                                                                          // 材料id
	MaterialName              string    `gorm:"column:material_name;type:varchar(255);not null" json:"materialName"`                                                               // 材料名称
	ExperimentMaterialGroupID string    `gorm:"index:experiment_material_group_id;column:experiment_material_group_id;type:varchar(36);not null" json:"experimentMaterialGroupId"` // 实验材料组id
	MaterialGroupName         string    `gorm:"column:material_group_name;type:varchar(255);not null" json:"materialGroupName"`                                                    // 材料组名称
	Percentage                float64   `gorm:"column:percentage;type:decimal(10,2);not null" json:"percentage"`                                                                   // 材料在组内的占比（%）
	CreatedAt                 time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`                                               // 创建时间
	UpdatedAt                 time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`                                               // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *Materials) TableName() string {
	return "materials"
}

// MaterialsColumns get sql column name.获取数据库列名
var MaterialsColumns = struct {
	ID                        string
	MaterialName              string
	ExperimentMaterialGroupID string
	MaterialGroupName         string
	Percentage                string
	CreatedAt                 string
	UpdatedAt                 string
}{
	ID:                        "id",
	MaterialName:              "material_name",
	ExperimentMaterialGroupID: "experiment_material_group_id",
	MaterialGroupName:         "material_group_name",
	Percentage:                "percentage",
	CreatedAt:                 "created_at",
	UpdatedAt:                 "updated_at",
}
