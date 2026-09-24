package config

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/hkdf"
)

const (
	// EnvCredentialEncryptionKey 通过环境变量注入凭据加密主密钥。
	EnvCredentialEncryptionKey = "GROK2API_CREDENTIAL_ENCRYPTION_KEY"
	jwtDerivationInfo          = "grok2api-jwt-v1"
)

func applySecrets(cfg *Config) error {
	if envKey := strings.TrimSpace(os.Getenv(EnvCredentialEncryptionKey)); envKey != "" {
		cfg.Secrets.CredentialEncryptionKey = envKey
	}
	if !validCredentialEncryptionKey(cfg.Secrets.CredentialEncryptionKey) {
		if cfg.Secrets.CredentialEncryptionKey == "" || isExampleSecret(cfg.Secrets.CredentialEncryptionKey) {
			return errors.New("请设置环境变量 GROK2API_CREDENTIAL_ENCRYPTION_KEY（openssl rand -base64 32）或在配置文件中提供 secrets.credentialEncryptionKey")
		}
		return errors.New("secrets.credentialEncryptionKey 必须是 Base64 编码的 32 字节密钥")
	}
	jwtSecret, err := deriveJWTSecret(cfg.Secrets.CredentialEncryptionKey)
	if err != nil {
		return err
	}
	cfg.Secrets.JWTSecret = jwtSecret
	return nil
}

func deriveJWTSecret(credentialEncryptionKey string) (string, error) {
	key, err := decodeCredentialEncryptionKey(credentialEncryptionKey)
	if err != nil {
		return "", err
	}
	reader := hkdf.New(sha256.New, key, nil, []byte(jwtDerivationInfo))
	derived := make([]byte, 32)
	if _, err := io.ReadFull(reader, derived); err != nil {
		return "", fmt.Errorf("派生 jwtSecret: %w", err)
	}
	return hex.EncodeToString(derived), nil
}

func decodeCredentialEncryptionKey(value string) ([]byte, error) {
	decoded, err := decodeBase64Key(value)
	if err != nil {
		return nil, fmt.Errorf("解析凭据加密密钥: %w", err)
	}
	if len(decoded) != 32 {
		return nil, errors.New("凭据加密密钥必须是 32 字节")
	}
	return decoded, nil
}