package entity

import (
	"time"
)

/******sql******
CREATE TABLE `operation_record` (
  `id` varchar(36) NOT NULL COMMENT 'id',
  `op_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作类型（RecordTypeData）',
  `user_id` varchar(36) NOT NULL COMMENT '操作用户',
  `desc` varchar(256) NOT NULL COMMENT '操作描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `index_user_createat` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='操作记录'
******sql******/
// OperationRecord 操作记录
type OperationRecord struct {
	ID        string    `gorm:"primaryKey;column:id;type:varchar(36);not null"`                                               // id
	OpType    string    `gorm:"column:op_type;type:varchar(64);not null"`                                                     // 操作类型（RecordTypeData）
	UserID    string    `gorm:"index:index_user_createat;column:user_id;type:varchar(36);not null"`                           // 操作用户
	Desc      string    `gorm:"column:desc;type:varchar(256);not null"`                                                       // 操作描述
	CreatedAt time.Time `gorm:"index:index_user_createat;column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"` // 创建时间
}

// TableName get sql table name.获取数据库表名
func (m *OperationRecord) TableName() string {
	return "operation_record"
}

// OperationRecordColumns get sql column name.获取数据库列名
var OperationRecordColumns = struct {
	ID        string
	OpType    string
	UserID    string
	Desc      string
	CreatedAt string
}{
	ID:        "id",
	OpType:    "op_type",
	UserID:    "user_id",
	Desc:      "desc",
	CreatedAt: "created_at",
}
