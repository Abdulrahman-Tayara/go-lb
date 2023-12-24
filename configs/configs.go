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
	HealthCheckIntervalSeconds int              `mapstructure:"health_check_interval_seconds" json:"health_check_interval_seconds" yaml:"health_check_interval_seconds"`
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
