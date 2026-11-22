// config 包负责从 SOML 文件加载 HTUNK 配置。
package config

import (
	"fmt"
	"os"

	"github.com/soml-lang/soml-go"
)

// Config 保存配置文件的顶层字段。
type Config struct {
	Log Log `soml:"log"`
}

// Log 保存日志配置。
type Log struct {
	Enable    bool   `soml:"enable"`
	Level     string `soml:"level"`
	Output    string `soml:"output"`
	Timestamp bool   `soml:"timestamp"`
}

// Load 读取配置文件并检查已知字段。
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read configuration: %w", err)
	}

	cfg := Config{Log: Log{Enable: true, Level: "info"}}
	if err := soml.Unmarshal(data, &cfg, soml.RejectUnknownKeys()); err != nil {
		return Config{}, fmt.Errorf("parse configuration %q: %w", path, err)
	}
	if cfg.Log.Enable {
		switch cfg.Log.Level {
		case "debug", "info", "warn", "error":
		default:
			return Config{}, fmt.Errorf("configuration %q: invalid log.level %q", path, cfg.Log.Level)
		}
	}
	return cfg, nil
}
