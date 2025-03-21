package entity

/******sql******
CREATE TABLE `material_groups` (
  `id` varchar(36) NOT NULL COMMENT '材料组id',
  `experiment_id` varchar(36) NOT NULL COMMENT '实验id',
  `material_group_category` enum('resin','hardener') NOT NULL DEFAULT 'resin' COMMENT '材料组类别（resin=树脂, hardener=固化剂）',
  `material_group_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '材料组名称(编号)',
  `proportion` decimal(10,2) NOT NULL COMMENT '材料组在配方中的占比',
  `sort` int NOT NULL,
  `parent_id` varchar(36) DEFAULT '' COMMENT '父级材料组id，无id表示顶级',
  `user_id` varchar(36) NOT NULL COMMENT '用户id',
  PRIMARY KEY (`id`),
  KEY `idx_experiment_id` (`experiment_id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='材料组表，支持层级结构'
******sql******/
// MaterialGroups 材料组表，支持层级结构
type MaterialGroups struct {
	ID                    string  `gorm:"primaryKey;column:id;type:varchar(36);not null"`                                      // 材料组id
	ExperimentID          string  `gorm:"index:idx_experiment_id;column:experiment_id;type:varchar(36);not null"`              // 实验id
	MaterialGroupCategory string  `gorm:"column:material_group_category;type:enum('resin','hardener');not null;default:resin"` // 材料组类别（resin=树脂, hardener=固化剂）
	MaterialGroupName     string  `gorm:"column:material_group_name;type:varchar(255);not null"`                               // 材料组名称(编号)
	Proportion            float64 `gorm:"column:proportion;type:decimal(10,2);not null"`                                       // 材料组在配方中的占比
	Sort                  int     `gorm:"column:sort;type:int;not null"`
	ParentID              *string `gorm:"index:idx_parent_id;column:parent_id;type:varchar(36);default:''"` // 父级材料组id，无id表示顶级
	UserID                string  `gorm:"column:user_id;type:varchar(36);not null"`                         // 用户id
}

// TableName get sql table name.获取数据库表名
func (m *MaterialGroups) TableName() string {
	return "material_groups"
}

// MaterialGroupsColumns get sql column name.获取数据库列名
var MaterialGroupsColumns = struct {
	ID                    string
	ExperimentID          string
	MaterialGroupCategory string
	MaterialGroupName     string
	Proportion            string
	Sort                  string
	ParentID              string
	UserID                string
}{
	ID:                    "id",
	ExperimentID:          "experiment_id",
	MaterialGroupCategory: "material_group_category",
	MaterialGroupName:     "material_group_name",
	Proportion:            "proportion",
	Sort:                  "sort",
	ParentID:              "parent_id",
	UserID:                "user_id",
}
