package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	lb2 "tayara/go-lb/lb"
	"tayara/go-lb/models"
	"tayara/go-lb/strategy"
	"testing"
)

func TestLoadBalancerWithCBR(t *testing.T) {
	servers := []*models.Server{
		{
			Name: "server1",
			Url:  "http://localhost:7070",
		},
		{
			Name: "server2",
			Url:  "http://localhost:7071",
		},
	}

	counters := make(map[string]int)
	for _, server := range servers {
		counters[server.Url] = 0
	}

	httpServers := setupServers(func(server *models.Server) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			counters[server.Url]++
		})
	}, servers...)

	defer httpServers.Close()

	routing := &models.Routing{
		Rules: models.RouteRules{
			{
				Conditions: []models.RouteCondition{
					{
						PathPrefix: "/api/v1",
						Headers: map[string]string{
							"My-Header": "my-value",
						},
					},
				},
				Action: models.RouteAction{RouteTo: "server1"},
			},
			{
				Conditions: []models.RouteCondition{
					{
						PathPrefix: "/api/v2",
						Method:     "POST",
					},
				},
				Action: models.RouteAction{RouteTo: "server2"},
			},
		},
		DefaultServer: "server1",
	}

	lb := lb2.NewLoadBalancer(
		servers,
		routing,
		strategy.NewRoundRobinStrategy(),
	)

	lbServer := models.Server{Url: "http://localhost:9080"}

	loadBalancerHttpServer := setupServers(func(server *models.Server) http.Handler {
		return lb
	}, &lbServer)

	defer loadBalancerHttpServer.Close()

	// Tests
	requests := []*http.Request{
		// Server1
		func() *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "http://localhost:9080/api/v1", nil)
			req.Header.Add("My-Header", "my-value")
			return req
		}(),
		// Default server
		func() *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "http://localhost:9080/api/v3", nil)
			return req
		}(),
		// Server2
		func() *http.Request {
			req, _ := http.NewRequest(http.MethodPost, "http://localhost:9080/api/v2", nil)
			return req
		}(),
	}

	wantCounters := map[string]int{
		servers[0].Url: 2,
		servers[1].Url: 1,
	}

	for i := 0; i < len(requests); i++ {
		_, err := http.DefaultClient.Do(requests[i])

		assert.NoError(t, err)

	}

	assert.Equal(t, wantCounters, counters)
}
