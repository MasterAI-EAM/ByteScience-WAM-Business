package entity

/******sql******
CREATE TABLE `recipe_material_groups` (
  `recipe_id` varchar(36) NOT NULL COMMENT '配方ID',
  `experiment_material_group_id` varchar(36) NOT NULL COMMENT '实验材料组ID',
  `proportion` decimal(10,2) NOT NULL COMMENT '材料组在配方中的占比（%）',
  PRIMARY KEY (`recipe_id`,`experiment_material_group_id`),
  KEY `idx_recipe_material_group_id` (`experiment_material_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='配方与实验材料组关系表'
******sql******/
// RecipeMaterialGroups 配方与实验材料组关系表
type RecipeMaterialGroups struct {
	RecipeID                  string  `gorm:"primaryKey;column:recipe_id;type:varchar(36);not null" json:"recipeId"`                                                                        // 配方ID
	ExperimentMaterialGroupID string  `gorm:"primaryKey;index:idx_recipe_material_group_id;column:experiment_material_group_id;type:varchar(36);not null" json:"experimentMaterialGroupId"` // 实验材料组ID
	Proportion                float64 `gorm:"column:proportion;type:decimal(10,2);not null" json:"proportion"`                                                                              // 材料组在配方中的占比（%）
}

// TableName get sql table name.获取数据库表名
func (m *RecipeMaterialGroups) TableName() string {
	return "recipe_material_groups"
}

// RecipeMaterialGroupsColumns get sql column name.获取数据库列名
var RecipeMaterialGroupsColumns = struct {
	RecipeID                  string
	ExperimentMaterialGroupID string
	Proportion                string
}{
	RecipeID:                  "recipe_id",
	ExperimentMaterialGroupID: "experiment_material_group_id",
	Proportion:                "proportion",
}
