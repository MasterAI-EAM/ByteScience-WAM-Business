package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword 对明文密码进行加密处理
// 参数:
//   - plainPassword: 明文密码，不能为空
//
// 返回:
//   - string: 加密后的密码字符串
//   - error: 处理过程中可能出现的错误
//
// 使用场景: 在用户注册或更新密码时，将明文密码加密存储到数据库。
func EncryptPassword(plainPassword string) (string, error) {
	if plainPassword == "" {
		return "", errors.New("password cannot be empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 验证明文密码是否匹配加密密码
// 参数:
//   - plainPassword: 明文密码
//   - hashedPassword: 数据库中存储的加密密码
//
// 返回:
//   - bool: 验证结果，true 表示匹配，false 表示不匹配
//   - error: 处理过程中可能出现的错误
//
// 使用场景: 在用户登录或身份验证时，比较用户输入的密码与数据库中的加密密码。
//
// 注意: 如果返回的 bool 为 false，说明密码不匹配；若返回 error，则为系统性错误。
func VerifyPassword(plainPassword, hashedPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
