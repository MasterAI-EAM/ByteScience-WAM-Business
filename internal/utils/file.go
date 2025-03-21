package utils

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"reflect"
)

func ExcelToJson(filePath, outputFile string) ([]string, error) {
	// 打开 Excel 文件
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开 Excel 文件: %v", err)
	}

	// 获取指定 sheet
	rows, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		return nil, fmt.Errorf("无法获取 sheet: %v", err)
	}

	// 初始化一个列表来存储所有列的 JSON 字符串
	var jsonStrings []string

	// 处理每一列
	numCols := len(rows[0]) // 假设所有行列数相同
	for colIndex := 0; colIndex < numCols; colIndex++ {
		var columnData []string
		for _, row := range rows {
			if colIndex < len(row) {
				columnData = append(columnData, row[colIndex])
			} else {
				columnData = append(columnData, "")
			}
		}

		// 构造 JSON 格式的列数据
		columnName := "column"
		columnDict := map[string][]string{columnName: columnData}
		jsonBytes, err := json.Marshal(columnDict)
		if err != nil {
			return nil, fmt.Errorf("JSON 编码失败: %v", err)
		}
		jsonStr := string(jsonBytes)

		// 添加到 jsonStrings 列表
		jsonStrings = append(jsonStrings, jsonStr)
	}

	// 写入 JSON 文件
	err = WriteJSONToFile(jsonStrings, outputFile)
	if err != nil {
		return nil, err
	}

	// 返回包含所有列数据的 JSON 字符串列表
	return jsonStrings, nil
}

// WriteJSONToFile 将任意类型的列表（支持 []string 和 []map[string]interface{}）写入文件
func WriteJSONToFile(data interface{}, outputFile string) error {
	// 获取输出文件的目录
	dir := filepath.Dir(outputFile)

	// 确保目录存在
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("无法创建输出目录: %v", err)
	}

	// 创建输出文件
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %v", err)
	}
	defer output.Close()

	// 获取传入数据的类型
	dataType := reflect.TypeOf(data)

	// 判断数据类型并进行相应处理
	switch dataType.Kind() {
	case reflect.Slice:
		// 判断是否是 []string 类型
		if dataType.Elem().Kind() == reflect.String {
			// 将 []string 转为 json 字符串写入文件
			jsonStrings, ok := data.([]string)
			if !ok {
				return fmt.Errorf("数据类型转换失败: 期望 []string 类型")
			}
			for _, jsonStr := range jsonStrings {
				_, err := output.WriteString(jsonStr + "\n")
				if err != nil {
					return fmt.Errorf("写入 JSON 文件失败: %v", err)
				}
			}
		} else if dataType.Elem().Kind() == reflect.Map {
			// 判断是否是 []map[string]interface{} 类型
			dataMaps, ok := data.([]map[string]interface{})
			if !ok {
				return fmt.Errorf("数据类型转换失败: 期望 []map[string]interface{} 类型")
			}
			// 将每个 map 转换为 JSON 并写入文件
			for _, item := range dataMaps {
				jsonBytes, err := json.Marshal(item)
				if err != nil {
					return fmt.Errorf("将 map 转换为 JSON 失败: %v", err)
				}
				_, err = output.WriteString(string(jsonBytes) + "\n")
				if err != nil {
					return fmt.Errorf("写入 JSON 文件失败: %v", err)
				}
			}
		} else {
			return fmt.Errorf("不支持的数据类型: %v", dataType.Kind())
		}
	default:
		return fmt.Errorf("不支持的数据类型: %v", dataType.Kind())
	}

	return nil
}
