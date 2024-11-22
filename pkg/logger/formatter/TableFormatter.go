package formatter

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

// TableFormatter 自定义的表格格式化器
type TableFormatter struct {
	Headers []string // 表头
}

// Format 实现 logrus.Formatter 接口的 Format 方法
func (f *TableFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var builder strings.Builder

	// 每列的内容
	timeStr := entry.Time.Format("2006-01-02 15:04:05")
	levelStr := entry.Level.String()
	messageStr := entry.Message

	// 计算每列的最大宽度
	maxWidths := make([]int, len(f.Headers))
	columns := []string{timeStr, levelStr, messageStr}
	for i, header := range f.Headers {
		maxWidths[i] = len(header) // 表头宽度
		if len(columns[i]) > maxWidths[i] {
			maxWidths[i] = len(columns[i]) // 内容宽度
		}
	}

	// 动态生成每列的分隔符线
	var separator strings.Builder
	separator.WriteString("+")
	for _, width := range maxWidths {
		separator.WriteString(strings.Repeat("-", width+2) + "+")
	}
	separator.WriteString("\n")

	// 输出表格的标题行
	builder.WriteString(separator.String())
	for i, header := range f.Headers {
		builder.WriteString(fmt.Sprintf("| %-*s ", maxWidths[i], header))
	}
	builder.WriteString("|\n")
	builder.WriteString(separator.String())

	// 输出日志字段行
	for i, col := range columns {
		builder.WriteString(fmt.Sprintf("| %-*s ", maxWidths[i], col))
	}
	builder.WriteString("|\n")
	builder.WriteString(separator.String())

	return []byte(builder.String()), nil
}
