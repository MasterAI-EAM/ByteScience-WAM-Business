package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

// PrettyJSONFormatter 是一个自定义的 JSON 格式化器，用于格式化输出的 JSON
type PrettyJSONFormatter struct {
	logrus.JSONFormatter
}

// Format 实现 logrus.Formatter 接口的 Format 方法
func (f *PrettyJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 使用 logrus.JSONFormatter 将日志条目格式化为紧凑的 JSON
	data, err := f.JSONFormatter.Format(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to format log entry as JSON: %v", err)
	}

	// 创建一个 bytes.Buffer，用于接收格式化后的 JSON
	var prettyJSON bytes.Buffer

	// 将紧凑的 JSON 数据进行格式化（添加缩进）
	err = json.Indent(&prettyJSON, data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to pretty print JSON: %v", err)
	}

	// 返回格式化后的 JSON 数据，并在末尾加上换行符
	return append(prettyJSON.Bytes(), '\n'), nil
}
