package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateClientAuthToken 生成客户端认证令牌
// 返回格式为 "sk-" + 48个随机字符的令牌
func GenerateClientAuthToken() (string, error) {
	// 生成 24 字节的随机数据 (24 * 2 = 48 hex characters)
	randomBytes := make([]byte, 24)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	
	// 转换为十六进制字符串
	randomHex := hex.EncodeToString(randomBytes)
	
	// 添加 "sk-" 前缀
	token := "sk-" + randomHex
	
	return token, nil
}