package configs

import (
	"github.com/spf13/viper"
	"tayara/go-lb/models"
)

type Configs struct {
	Port                 string           `mapstructure:"port" json:"port" yaml:"port"`
	LoadBalancerStrategy string           `mapstructure:"load_balancer_strategy" json:"load_balancer_strategy" yaml:"load_balancer_strategy"`
	Servers              []*models.Server `mapstructure:"servers" json:"servers" yaml:"servers"`
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

	return &configs, nil
}
