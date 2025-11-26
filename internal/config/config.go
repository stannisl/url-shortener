package config

import (
	"fmt"
	"time"

	cfg "github.com/wb-go/wbf/config"
)

const (
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 10 * time.Second
	DefaultShutdownTimeout = 15 * time.Second
	DefaultHTTPHost        = "0.0.0.0"
	DefaultHTTPPort        = "8080"
)

type Config struct {
	ServerConfig   ServerConfig   `mapstructure:"server"`
	LoggerConfig   LoggerConfig   `mapstructure:"logger"`
	PostgresConfig PostgresConfig `mapstructure:"postgres"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type LoggerConfig struct {
	LogLevel string `mapstructure:"log_level"`
}

type PostgresConfig struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	DatabaseUrl     string        `mapstructure:"DATABASE_URL"`
}

func LoadConfig(configPath string, envFilePath string) (*Config, error) {
	c := cfg.New()

	c.EnableEnv("")
	if err := c.LoadConfigFiles(configPath); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}
	if err := c.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	var conf Config
	if err := c.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	conf.PostgresConfig.DatabaseUrl = c.GetString("DATABASE_URL")

	return &conf, nil
}
