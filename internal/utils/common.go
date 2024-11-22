package utils

import "regexp"

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
