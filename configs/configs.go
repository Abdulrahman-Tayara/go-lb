package configs

import (
	"github.com/spf13/viper"
	"tayara/go-lb/models"
)

const (
	defaultHealthCheckIntervalSeconds = 5
)

type Configs struct {
	Port                       string           `mapstructure:"port" json:"port" yaml:"port"`
	LoadBalancerStrategy       string           `mapstructure:"load_balancer_strategy" json:"load_balancer_strategy" yaml:"load_balancer_strategy"`
	Servers                    []*models.Server `mapstructure:"servers" json:"servers" yaml:"servers"`
	Routing                    models.Routing   `mapstructure:"routing" json:"routing" yaml:"routing"`
	HealthCheckIntervalSeconds int              `mapstructure:"health_check_interval_seconds" json:"health_check_interval_seconds" yaml:"health_check_interval_seconds"`
	RateLimiterEnabled         bool             `mapstructure:"rate_limiter_enabled" json:"rate_limiter_enabled" yaml:"rate_limiter_enabled"`
	RateLimitTokens            int              `mapstructure:"rate_limit_tokens" json:"rate_limit_tokens" yaml:"rate_limit_tokens"`
	RateLimitIntervalSeconds   int              `mapstructure:"rate_limit_interval_seconds" json:"rate_limit_interval_seconds" yaml:"rate_limit_interval_seconds"`
	TLSEnabled                 bool             `mapstructure:"tls_enabled" json:"tls_enabled" yaml:"tls_enabled"`
	TLSCertPath                string           `mapstructure:"tls_cert_path" json:"tls_cert_path" yaml:"tls_cert_path"`
	TLSKeyPath                 string           `mapstructure:"tls_key_path" json:"tls_key_path" yaml:"tls_key_path"`
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

	return &configs, nil
}
