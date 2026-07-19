package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Storage StorageConfig `mapstructure:"storage"`
	Cache   CacheConfig   `mapstructure:"cache"`
	Image   ImageConfig   `mapstructure:"image"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Image        ServerImage   `mapstructure:"image"`
}

type ServerImage struct {
	MaxWidth  int `mapstructure:"max_width"`
	MaxHeight int `mapstructure:"max_height"`
}

type StorageConfig struct {
	S3 S3Config `mapstructure:"s3"`
}

type S3Config struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

type CacheConfig struct {
	FS FSCacheConfig `mapstructure:"fs"`
}

type FSCacheConfig struct {
	Dir string `mapstructure:"dir"`
}

type ImageConfig struct {
	ContainBackgroundColor string `mapstructure:"contain_background_color"`
}

// SetDefaults sets sensible baseline defaults for Viper.
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.image.max_width", 5000)
	v.SetDefault("server.image.max_height", 5000)
	v.SetDefault("cache.fs.dir", "/tmp/optivor-cache")
	v.SetDefault("image.contain_background_color", "#ffffff")
	v.SetDefault("storage.s3.region", "us-east-1")
}

// Load reads configuration using Viper and validates required fields.
// If configFilePath is non-empty, it overrides default discovery.
func Load(configFilePath string) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	if configFilePath != "" {
		v.SetConfigFile(configFilePath)
	} else {
		v.SetConfigName("optivor")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	v.SetEnvPrefix("OPTIVOR")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) && configFilePath != "" {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := Validate(&cfg); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return &cfg, nil
}

// Validate checks that all required parameters are specified.
func Validate(cfg *Config) error {
	if cfg.Server.Port <= 0 {
		return errors.New("server.port must be > 0")
	}
	if cfg.Server.Image.MaxWidth <= 0 {
		return errors.New("server.image.max_width must be > 0")
	}
	if cfg.Server.Image.MaxHeight <= 0 {
		return errors.New("server.image.max_height must be > 0")
	}
	if cfg.Storage.S3.Endpoint == "" {
		return errors.New("storage.s3.endpoint is required")
	}
	if cfg.Storage.S3.Bucket == "" {
		return errors.New("storage.s3.bucket is required")
	}
	if cfg.Cache.FS.Dir == "" {
		return errors.New("cache.fs.dir is required")
	}
	if cfg.Image.ContainBackgroundColor == "" {
		cfg.Image.ContainBackgroundColor = "#ffffff"
	}
	return nil
}
