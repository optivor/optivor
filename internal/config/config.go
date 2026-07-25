package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Buckets   []BucketConfig  `mapstructure:"buckets"`
	Cache     CacheConfig     `mapstructure:"cache"`
	Image     ImageConfig     `mapstructure:"image"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
	Remote    RemoteConfig    `mapstructure:"remote"`
	Presets   map[string]PresetConfig `mapstructure:"presets"`
	Crawler   CrawlerConfig   `mapstructure:"crawler"`
}

type RemoteConfig struct {
	Enabled        bool     `mapstructure:"enabled"`
	AllowedDomains []string `mapstructure:"allowed_domains"`
}

type PresetConfig struct {
	Width  int    `mapstructure:"w"`
	Height int    `mapstructure:"h"`
	Format string `mapstructure:"f"`
	Fit    string `mapstructure:"fit"`
	Quality int   `mapstructure:"q"`
}

type CrawlerConfig struct {
	Enabled               bool `mapstructure:"enabled"`
	MaxConcurrencyPerVariant int  `mapstructure:"max_concurrency_per_variant"`
}

type BucketConfig struct {
	Name            string `mapstructure:"name"`
	Provider        string `mapstructure:"provider"`
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	Region          string `mapstructure:"region"`
	AccountID       string `mapstructure:"account_id"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Access          string `mapstructure:"access"`
	Fallback        string `mapstructure:"fallback"`
}

type TelemetryConfig struct {
	Enabled       bool    `mapstructure:"enabled"`
	OTLPEndpoint  string  `mapstructure:"otlp_endpoint"`
	ServiceName   string  `mapstructure:"service_name"`
	SamplingRatio float64 `mapstructure:"sampling_ratio"`
}

type ServerConfig struct {
	Port           int             `mapstructure:"port"`
	ReadTimeout    time.Duration   `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration   `mapstructure:"write_timeout"`
	RequestTimeout time.Duration   `mapstructure:"request_timeout"`
	LogLevel       string          `mapstructure:"log_level"`
	LogFormat      string          `mapstructure:"log_format"`
	Image          ServerImage     `mapstructure:"image"`
	RateLimit      RateLimitConfig `mapstructure:"rate_limit"`
}

type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	RPS     int  `mapstructure:"rps"`
	Burst   int  `mapstructure:"burst"`
}

type ServerImage struct {
	MaxWidth  int `mapstructure:"max_width"`
	MaxHeight int `mapstructure:"max_height"`
}

type StorageConfig struct {
	Driver string   `mapstructure:"driver"`
	S3     S3Config `mapstructure:"s3"`
}

type S3Config struct {
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

type CacheConfig struct {
	Type  string           `mapstructure:"type"`
	FS    FSCacheConfig    `mapstructure:"fs"`
	Redis RedisCacheConfig `mapstructure:"redis"`
}

type RedisCacheConfig struct {
	Addr                       string        `mapstructure:"addr"`
	Password                   string        `mapstructure:"password"`
	DB                         int           `mapstructure:"db"`
	Prefix                     string        `mapstructure:"prefix"`
	TTL                        time.Duration `mapstructure:"ttl"`
	PoolSize                   int           `mapstructure:"pool_size"`
	MinIdleConns               int           `mapstructure:"min_idle_conns"`
	CircuitBreakerMaxFailures int           `mapstructure:"circuit_breaker_max_failures"`
	CircuitBreakerTimeout     time.Duration `mapstructure:"circuit_breaker_timeout"`
}

type FSCacheConfig struct {
	Dir       string `mapstructure:"dir"`
	MaxSizeMB int64  `mapstructure:"max_size_mb"`
}

type ImageConfig struct {
	ContainBackgroundColor string `mapstructure:"contain_background_color"`
	MaxPixels              int    `mapstructure:"max_pixels"`
	MaxDecodeMB            int    `mapstructure:"max_decode_mb"`
}

type AuthConfig struct {
	SignedURLs SignedURLsConfig `mapstructure:"signed_urls"`
	APIKeys    []APIKeyConfig   `mapstructure:"api_keys"`
}

type APIKeyConfig struct {
	Key     string   `mapstructure:"key"`
	Name    string   `mapstructure:"name"`
	Buckets []string `mapstructure:"buckets"`
	Scopes  []string `mapstructure:"scopes"`
}

type SignedURLsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Secret  string `mapstructure:"secret"`
	MaxAge  int    `mapstructure:"max_age"`
}

// SetDefaults sets sensible baseline defaults for Viper.
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30*time.Second)
	v.SetDefault("server.write_timeout", 30*time.Second)
	v.SetDefault("server.request_timeout", 30*time.Second)
	v.SetDefault("server.log_level", "info")
	v.SetDefault("server.log_format", "text")
	v.SetDefault("server.image.max_width", 5000)
	v.SetDefault("server.image.max_height", 5000)
	v.SetDefault("cache.type", "fs")
	v.SetDefault("cache.fs.dir", "/tmp/optivor-cache")
	v.SetDefault("cache.fs.max_size_mb", 1024)
	v.SetDefault("cache.redis.addr", "localhost:6379")
	v.SetDefault("cache.redis.prefix", "optivor:cache:")
	v.SetDefault("cache.redis.ttl", 24*time.Hour)
	v.SetDefault("cache.redis.pool_size", 10)
	v.SetDefault("cache.redis.min_idle_conns", 5)
	v.SetDefault("cache.redis.circuit_breaker_max_failures", 5)
	v.SetDefault("cache.redis.circuit_breaker_timeout", 30*time.Second)
	v.SetDefault("image.contain_background_color", "#ffffff")
	v.SetDefault("image.max_pixels", 25000000)
	v.SetDefault("image.max_decode_mb", 64)
	v.SetDefault("storage.driver", "s3")
	v.SetDefault("storage.s3.region", "us-east-1")
	v.SetDefault("auth.signed_urls.enabled", false)
	v.SetDefault("auth.signed_urls.max_age", 3600)
	v.SetDefault("server.rate_limit.enabled", true)
	v.SetDefault("server.rate_limit.rps", 10)
	v.SetDefault("server.rate_limit.burst", 20)
	v.SetDefault("telemetry.enabled", false)
	v.SetDefault("telemetry.service_name", "optivor")
	v.SetDefault("telemetry.sampling_ratio", 1.0)
	v.SetDefault("remote.enabled", true)
	v.SetDefault("remote.allowed_domains", []string{"*"})
	v.SetDefault("crawler.enabled", true)
	v.SetDefault("crawler.max_concurrency_per_variant", 10)
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
	if len(cfg.Buckets) > 0 {
		for i, b := range cfg.Buckets {
			if b.Name == "" {
				return fmt.Errorf("buckets[%d].name is required", i)
			}
			if b.Endpoint == "" {
				return fmt.Errorf("buckets[%d].endpoint is required", i)
			}
			if b.Bucket == "" {
				return fmt.Errorf("buckets[%d].bucket is required", i)
			}
		}
	} else {
		if cfg.Storage.S3.Endpoint == "" {
			return errors.New("storage.s3.endpoint is required")
		}
		if cfg.Storage.S3.Bucket == "" {
			return errors.New("storage.s3.bucket is required")
		}
	}
	if cfg.Cache.Type == "" {
		cfg.Cache.Type = "fs"
	}
	if cfg.Cache.Type != "fs" && cfg.Cache.Type != "redis" {
		return fmt.Errorf("cache.type must be either 'fs' or 'redis', got: %s", cfg.Cache.Type)
	}
	if cfg.Cache.Type == "fs" && cfg.Cache.FS.Dir == "" {
		return errors.New("cache.fs.dir is required when cache.type is 'fs'")
	}
	if cfg.Cache.Type == "redis" && cfg.Cache.Redis.Addr == "" {
		return errors.New("cache.redis.addr is required when cache.type is 'redis'")
	}
	if cfg.Image.ContainBackgroundColor == "" {
		cfg.Image.ContainBackgroundColor = "#ffffff"
	}
	if cfg.Auth.SignedURLs.Enabled && cfg.Auth.SignedURLs.Secret == "" {
		return errors.New("auth.signed_urls.secret is required when auth.signed_urls.enabled is true")
	}
	return nil
}
