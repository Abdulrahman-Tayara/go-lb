package lb

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"tayara/go-lb/models"
	"tayara/go-lb/strategy"
	"testing"
)

func TestLoadBalancer_ServerDown(t *testing.T) {
	servers := []*models.Server{
		{
			Url: "http://loclhost:8080",
		},
		{
			Url: "http://loclhost:8081",
		},
	}

	lb := NewLoadBalancer(slices.Clone(servers), nil, strategy.NewRoundRobinStrategy()).(*loadBalancer)

	assert.Equal(t, len(servers), len(lb.servers))

	lb.ServerDown(servers[0])
	lb.ServerDown(servers[0])

	assert.Equal(t, len(servers)-1, len(lb.servers))
	assert.True(t, lb.servers[0].Equals(servers[1]))
}

func TestLoadBalancer_ServerUp(t *testing.T) {
	servers := []*models.Server{
		{
			Url: "http://loclhost:8080",
		},
	}

	lb := NewLoadBalancer(slices.Clone(servers), nil, strategy.NewRoundRobinStrategy()).(*loadBalancer)

	assert.Equal(t, len(servers), len(lb.servers))

	lb.ServerUp(servers[0])
	assert.Equal(t, len(servers), len(lb.servers))

	lb.ServerUp(&models.Server{Url: "http://loclhost:8081"})

	assert.Equal(t, len(servers)+1, len(lb.servers))
}
