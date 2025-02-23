package entity

import (
	"time"
)

/******sql******
CREATE TABLE `recipes` (
  `id` varchar(36) NOT NULL COMMENT '配方ID',
  `recipe_name` varchar(255) NOT NULL COMMENT '配方名称',
  `recipe_signature` varchar(64) NOT NULL COMMENT '配方唯一签名（哈希值）',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_recipe_signature` (`recipe_signature`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='配方表'
******sql******/
// Recipes 配方表
type Recipes struct {
	ID              string    `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                            // 配方ID
	RecipeName      string    `gorm:"column:recipe_name;type:varchar(255);not null" json:"recipeName"`                     // 配方名称
	RecipeSignature string    `gorm:"unique;column:recipe_signature;type:varchar(64);not null" json:"recipeSignature"`     // 配方唯一签名（哈希值）
	Sort            int       `gorm:"column:sort;type:int;not null;default:0" json:"sort"`                                 // 排序
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"` // 创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"` // 修改时间
}

// TableName get sql table name.获取数据库表名
func (m *Recipes) TableName() string {
	return "recipes"
}

// RecipesColumns get sql column name.获取数据库列名
var RecipesColumns = struct {
	ID              string
	RecipeName      string
	RecipeSignature string
	Sort            string
	CreatedAt       string
	UpdatedAt       string
}{
	ID:              "id",
	RecipeName:      "recipe_name",
	RecipeSignature: "recipe_signature",
	Sort:            "sort",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}
