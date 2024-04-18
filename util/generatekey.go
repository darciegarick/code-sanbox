package util

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

// GenerateAccessKey 生成 用户的唯一标识（访问密钥）
func GenerateAccessKey() string {
	return uuid.NewString()
}

// GenerateSecureSecretKey 生成一个指定长度的安全的secretKey
func GenerateSecureKey() (string, error) {
	var byteLength int = 32
	// 创建一个足够大的空间来保存随机字节
	randomBytes := make([]byte, byteLength)
	// 使用crypto/rand生成随机字节
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果生成过程中发生错误，返回错误
		return "", err
	}
	// 将随机字节序列编码为十六进制字符串
	secretKey := hex.EncodeToString(randomBytes)
	return secretKey, nil
}

// GenerateSign 根据 用户标识、密钥、随机数、时间戳生成的
func GenerateSign(accessKey, secretKey, nonce string, timestamp int64) string {
	// 将参数按照一定规则拼接在一起
	data := fmt.Sprintf("%s%s%d%s", accessKey, secretKey, timestamp, nonce)
	// 使用 HMAC-SHA256 算法计算哈希值
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))
	// 返回签名结果
	return signature
}
