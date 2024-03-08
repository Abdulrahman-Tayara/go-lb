package strategy

import (
	"tayara/go-lb/models"
	"tayara/go-lb/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedRoundRobin(t *testing.T) {
	strategy := NewWeightedRoundRobinStrategy()

	servers := []*models.Server{
		{
			Name:   "Server1",
			Url:    "http://localhost:8080",
			Weight: 2,
		},
		{
			Name:   "Server2",
			Url:    "http://localhost:8081",
			Weight: 3,
		},
		{
			Name:   "Server3",
			Url:    "http://localhost:8082",
			Weight: 1,
		},
	}

	strategy.UpdateServers(servers)

	requestsCount := 6

	actualCounts := map[string]int{}

	for i := 0; i < requestsCount; i++ {
		server := strategy.Next(nil)

		actualCounts[server.Name] = utils.GetOrDefault(actualCounts, server.Name, 0) + 1

		strategy.RequestServed(server, nil)
	}

	expectedCounts := map[string]int{
		servers[0].Name: 2,
		servers[1].Name: 3,
		servers[2].Name: 1,
	}

	assert.Equal(t, expectedCounts, actualCounts)
}
