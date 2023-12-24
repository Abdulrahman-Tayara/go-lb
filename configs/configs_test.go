package configs

import (
	"reflect"
	"tayara/go-lb/models"
	"testing"
)

func TestLoadConfigs(t *testing.T) {
	tests := []struct {
		name    string
		input   string // filepath
		want    Configs
		wantErr bool
	}{
		{
			"Load not found configs file",
			"./not_found.yaml",
			Configs{},
			true,
		},
		{
			"Load json file",
			"configs_test.json",
			Configs{
				Port: "9090",
				Servers: []*models.Server{
					{
						Url:       "http://localhost:8080",
						HealthUrl: "http://localhost:8080/health",
					},
				},
				LoadBalancerStrategy: "round_robin",
			},
			false,
		},
		{
			"Load yaml file",
			"configs_test.yml",
			Configs{
				Port: "8900",
				Servers: []*models.Server{
					{
						Url:       "http://localhost:8080",
						HealthUrl: "http://localhost:8080/health",
					},
					{
						Url:       "http://localhost:8081",
						HealthUrl: "http://localhost:8081/health",
					},
				},
				LoadBalancerStrategy: "round_robin",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := LoadConfigs(tt.input)

			if tt.wantErr && gotErr == nil {
				t.Errorf("execpted error")
			} else if !tt.wantErr && gotErr != nil {
				t.Errorf("exepcted %v, got error %v", tt.want, gotErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.want, *got) {
				t.Errorf("exepcted %v, got %v", tt.want, got)
			}
		})
	}
}
