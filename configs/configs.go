package configs

import (
	"tayara/go-lb/models"
	"tayara/go-lb/strategy"

	"github.com/spf13/viper"
)

const (
	defaultHealthCheckIntervalSeconds = 5
	defaultLogLevel                   = "info"
)

type Configs struct {
	Port                       string           `mapstructure:"port" json:"port" yaml:"port"`
	LoadBalancerStrategy       string           `mapstructure:"load_balancer_strategy" json:"load_balancer_strategy" yaml:"load_balancer_strategy"`
	StrategyConfigs            strategy.Configs `mapstructure:"strategy_configs" json:"strategy_configs" yaml:"strategy_configs"`
	Servers                    []*models.Server `mapstructure:"servers" json:"servers" yaml:"servers"`
	Routing                    models.Routing   `mapstructure:"routing" json:"routing" yaml:"routing"`
	HealthCheckIntervalSeconds int              `mapstructure:"health_check_interval_seconds" json:"health_check_interval_seconds" yaml:"health_check_interval_seconds"`
	RateLimiterEnabled         bool             `mapstructure:"rate_limiter_enabled" json:"rate_limiter_enabled" yaml:"rate_limiter_enabled"`
	RateLimitTokens            int              `mapstructure:"rate_limit_tokens" json:"rate_limit_tokens" yaml:"rate_limit_tokens"`
	RateLimitIntervalSeconds   int              `mapstructure:"rate_limit_interval_seconds" json:"rate_limit_interval_seconds" yaml:"rate_limit_interval_seconds"`
	TLSEnabled                 bool             `mapstructure:"tls_enabled" json:"tls_enabled" yaml:"tls_enabled"`
	TLSCertPath                string           `mapstructure:"tls_cert_path" json:"tls_cert_path" yaml:"tls_cert_path"`
	TLSKeyPath                 string           `mapstructure:"tls_key_path" json:"tls_key_path" yaml:"tls_key_path"`
	LogFile                    string           `mapstructure:"log_file" json:"log_file" yaml:"log_file"`
	LogLevel                   string           `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
}

func LoadConfigs(path string) (*Configs, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var configs Configs

	err = viper.Unmarshal(&configs)
	if err != nil {
		return nil, err
	}

	if configs.HealthCheckIntervalSeconds == 0 {
		configs.HealthCheckIntervalSeconds = defaultHealthCheckIntervalSeconds
	}

	if configs.LogLevel == "" {
		configs.LogLevel = defaultLogLevel
	}

	return &configs, nil
}
