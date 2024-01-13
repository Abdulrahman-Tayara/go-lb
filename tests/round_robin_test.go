package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	lb2 "tayara/go-lb/lb"
	"tayara/go-lb/models"
	"tayara/go-lb/strategy"
	"testing"
)

func TestLoadBalancerWithRoundRobin(t *testing.T) {
	servers := []*models.Server{
		{
			Url: "http://localhost:7070",
		},
		{
			Url: "http://localhost:7071",
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

	lb := lb2.NewLoadBalancer(
		servers,
		nil,
		strategy.NewRoundRobinStrategy(),
	)

	lbServer := models.Server{Url: "http://localhost:9080"}

	loadBalancerHttpServer := setupServers(func(server *models.Server) http.Handler {
		return lb
	}, &lbServer)

	defer loadBalancerHttpServer.Close()

	// Tests
	numberOfCalls := 10

	wantCounters := map[string]int{
		servers[0].Url: 5,
		servers[1].Url: 5,
	}

	for i := 0; i < numberOfCalls; i++ {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:9080", nil)
		_, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)

	}

	assert.Equal(t, wantCounters, counters)
}
