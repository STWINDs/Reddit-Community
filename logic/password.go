package logic

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 将明文密码加密为哈希值
func HashPassword(password string) (string, error) {
	// DefaultCost 目前是 10，是一个在性能和安全性之间很好的平衡
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword 检查明文密码和哈希值是否匹配
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
