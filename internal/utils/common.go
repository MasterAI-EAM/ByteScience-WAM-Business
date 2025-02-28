package utils

import (
	"regexp"
	"time"
)

// IdentifyType 用于判断标识符类型（用户名、邮箱或手机号）
func IdentifyType(identifier string) string {
	// 判断是否为邮箱
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, identifier); matched {
		return "email"
	}

	// 判断是否为手机号（国际格式，E.164）
	phoneRegex := `^\+?[1-9]\d{1,14}$`
	if matched, _ := regexp.MatchString(phoneRegex, identifier); matched {
		return "phone"
	}

	// 如果都不匹配，默认为用户名
	return "username"
}

// Contains 判断切片是否包含某个元素，支持多种类型
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Difference 找出a中但不在b中的元素，支持泛型类型
func Difference[T comparable](a, b []T) []T {
	// 创建一个哈希集合来存储b的元素
	bSet := make(map[T]struct{})
	for _, item := range b {
		bSet[item] = struct{}{}
	}

	// 遍历数组a，找出不在b中的元素
	var result []T
	for _, item := range a {
		if _, found := bSet[item]; !found {
			result = append(result, item)
		}
	}
	return result
}

// RemoveDuplicates 去重函数，适用于任何可比较类型
func RemoveDuplicates[T comparable](arr []T) []T {
	seen := make(map[T]struct{}) // 用于存储已遇到的元素
	result := make([]T, 0, len(arr))

	for _, item := range arr {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z")
}
