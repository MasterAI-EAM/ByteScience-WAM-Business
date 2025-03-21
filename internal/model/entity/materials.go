package entity

/******sql******
CREATE TABLE `materials` (
  `id` varchar(36) NOT NULL COMMENT '材料id',
  `material_name` varchar(255) NOT NULL COMMENT '材料名称',
  `material_group_id` varchar(36) NOT NULL COMMENT '实验材料组id',
  `proportion` decimal(10,2) NOT NULL COMMENT '材料在组内的占比（%）',
  `sort` int NOT NULL,
  `user_id` varchar(36) NOT NULL COMMENT '用户id',
  PRIMARY KEY (`id`),
  KEY `idx_material_group_id` (`material_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='材料表'
******sql******/
// Materials 材料表
type Materials struct {
	ID              string  `gorm:"primaryKey;column:id;type:varchar(36);not null"`                                 // 材料id
	MaterialName    string  `gorm:"column:material_name;type:varchar(255);not null"`                                // 材料名称
	MaterialGroupID string  `gorm:"index:idx_material_group_id;column:material_group_id;type:varchar(36);not null"` // 实验材料组id
	Proportion      float64 `gorm:"column:proportion;type:decimal(10,2);not null"`                                  // 材料在组内的占比（%）
	Sort            int     `gorm:"column:sort;type:int;not null"`
	UserID          string  `gorm:"column:user_id;type:varchar(36);not null"` // 用户id
}

// TableName get sql table name.获取数据库表名
func (m *Materials) TableName() string {
	return "materials"
}

// MaterialsColumns get sql column name.获取数据库列名
var MaterialsColumns = struct {
	ID              string
	MaterialName    string
	MaterialGroupID string
	Proportion      string
	Sort            string
	UserID          string
}{
	ID:              "id",
	MaterialName:    "material_name",
	MaterialGroupID: "material_group_id",
	Proportion:      "proportion",
	Sort:            "sort",
	UserID:          "user_id",
}
