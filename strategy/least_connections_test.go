package strategy

import (
	"github.com/stretchr/testify/assert"
	"tayara/go-lb/models"
	"testing"
)

func TestLeastConnections_getLeastConnections(t *testing.T) {
	servers := []*models.Server{
		{
			Url: "http://localhost:8080",
		},
		{
			Url: "http://localhost:8081",
		},
	}

	s := NewLeastConnectionsStrategy().(*LeastConnectionsStrategy)
	s.UpdateServers(servers)

	server1Connections := 3
	server2Connections := 2

	for i := 0; i < server1Connections; i++ {
		s.increaseConnections(servers[0])
	}

	for i := 0; i < 100; i++ {
		s.decreaseConnections(servers[1])
	}

	for i := 0; i < server2Connections; i++ {
		s.increaseConnections(servers[1])
	}

	assert.Equal(t, int64(server1Connections), s.serversConnections[servers[0]])
	assert.Equal(t, int64(server2Connections), s.serversConnections[servers[1]])

	for i := 0; i < 100; i++ {
		s.decreaseConnections(servers[1])
	}

	assert.Equal(t, int64(0), s.serversConnections[servers[1]])
}

func TestLeastConnections_UpdateServers(t *testing.T) {
	servers := []*models.Server{
		{
			Url: "http://localhost:8080",
		},
		{
			Url: "http://localhost:8081",
		},
	}

	s := NewLeastConnectionsStrategy().(*LeastConnectionsStrategy)
	s.UpdateServers(servers)

	for _, server := range servers {
		s.increaseConnections(server)
	}

	for _, server := range servers {
		assert.Equal(t, int64(1), s.serversConnections[server])
	}

	newServers := []*models.Server{
		{
			Url: "http://localhost:8081",
		},
	}

	s.UpdateServers(newServers)

	assert.Equal(t, int64(0), s.serversConnections[servers[0]])
	assert.Equal(t, int64(1), s.serversConnections[servers[1]])
}

func TestLeastConnections_UpdateServers_EmptyList(t *testing.T) {
	servers := []*models.Server{
		{
			Url: "http://localhost:8080",
		},
		{
			Url: "http://localhost:8081",
		},
	}

	s := NewLeastConnectionsStrategy().(*LeastConnectionsStrategy)
	s.UpdateServers(servers)

	var newServers []*models.Server

	s.UpdateServers(newServers)

	assert.Equal(t, 0, len(s.serversConnections))
}
