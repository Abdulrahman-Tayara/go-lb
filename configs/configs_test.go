package configs

import (
	"reflect"
	"tayara/go-lb/models"
	"tayara/go-lb/strategy"
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
						Name:      "server1",
						Url:       "http://localhost:8080",
						HealthUrl: "http://localhost:8080/health",
					},
				},
				LoadBalancerStrategy: "round_robin",
				StrategyConfigs: strategy.Configs{
					StickySessionCookieName: "example",
					StickySessionTTLSeconds: 100,
				},
				HealthCheckIntervalSeconds: 5,
				RateLimiterEnabled:         true,
				RateLimitIntervalSeconds:   10,
				RateLimitTokens:            10,
				LogLevel:                   "info",
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
						Name:      "server1",
						Url:       "http://localhost:8080",
						HealthUrl: "http://localhost:8080/health",
					},
					{
						Name:      "server2",
						Url:       "http://localhost:8081",
						HealthUrl: "http://localhost:8081/health",
					},
				},
				Routing: models.Routing{
					DefaultServer: "server2",
					Rules: models.RouteRules{
						{
							Conditions: []models.RouteCondition{
								{
									PathPrefix: "/api/v1",
									Method:     "GET",
									Headers: map[string]string{
										"useragent": "Mobile",
									},
								},
							},
							Action: models.RouteAction{
								RouteTo: "server1",
							},
						},
					},
				},
				LoadBalancerStrategy: "round_robin",
				StrategyConfigs: strategy.Configs{
					StickySessionCookieName: "example",
					StickySessionTTLSeconds: 100,
				},
				HealthCheckIntervalSeconds: 3,
				RateLimiterEnabled:         true,
				RateLimitIntervalSeconds:   10,
				RateLimitTokens:            10,
				LogLevel:                   "info",
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
