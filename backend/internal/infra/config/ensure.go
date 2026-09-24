package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultBootstrapUsername = "admin"
	DefaultBootstrapPassword = "grok2api"
)

type bootstrapFile struct {
	Frontend FrontendConfig `yaml:"frontend"`
	Database DatabaseConfig `yaml:"database"`
	Media    MediaConfig    `yaml:"media"`
}

// Ensure 在配置文件不存在时写入默认启动配置。
func Ensure(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("配置文件路径不能为空")
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("读取配置文件: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("创建配置目录: %w", err)
	}
	staticPath := "./frontend/dist"
	if info, err := os.Stat("/app/frontend/dist"); err == nil && info.IsDir() {
		staticPath = "/app/frontend/dist"
	}
	payload := bootstrapFile{
		Frontend: FrontendConfig{
			PublicAPIBaseURL: "http://127.0.0.1:8000",
			StaticPath:       staticPath,
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			SQLite: SQLiteDatabaseConfig{Path: "./data/backend.db"},
		},
		Media: MediaConfig{
			Driver: "local",
			Local:  LocalMediaConfig{Path: "./data/media"},
		},
	}
	data, err := yaml.Marshal(payload)
	if err != nil {
		return fmt.Errorf("生成默认配置: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("写入默认配置: %w", err)
	}
	return nil
}