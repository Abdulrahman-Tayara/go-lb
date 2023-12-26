package strategy

import (
	"math"
	"net/http"
	"sync"
	"tayara/go-lb/models"
)

type LeastConnectionsStrategy struct {
	servers []*models.Server

	serversConnections map[*models.Server]int64

	sync.RWMutex
}

func NewLeastConnectionsStrategy() ILoadBalancerStrategy {
	return &LeastConnectionsStrategy{
		serversConnections: make(map[*models.Server]int64),
	}
}

func (l *LeastConnectionsStrategy) getLeastConnections() *models.Server {
	minConnections := int64(math.MaxInt64)
	var selectedServer *models.Server

	for server, connections := range l.serversConnections {
		if connections < minConnections {
			minConnections = connections
			selectedServer = server
		}
	}

	return selectedServer
}

func (l *LeastConnectionsStrategy) increaseConnections(s *models.Server) {
	if _, exists := l.serversConnections[s]; !exists {
		l.serversConnections[s] = 0
	}
	l.serversConnections[s]++
}

func (l *LeastConnectionsStrategy) decreaseConnections(s *models.Server) {
	if _, exists := l.serversConnections[s]; exists {
		l.serversConnections[s] = max(0, l.serversConnections[s]-1)
	}
}

func (l *LeastConnectionsStrategy) Next(request *http.Request) *models.Server {
	defer l.Unlock()

	defer l.Lock()

	selectedServer := l.getLeastConnections()
	if selectedServer == nil {
		return nil
	}

	l.increaseConnections(selectedServer)

	return selectedServer
}

func (l *LeastConnectionsStrategy) UpdateServers(servers []*models.Server) {
	defer l.Unlock()
	l.Lock()

	contains := func(s *models.Server, list []*models.Server) bool {
		for _, server := range list {
			if s.Equals(server) {
				return true
			}
		}
		return false
	}

	for _, server := range l.servers {
		// If the server isn't exist in the new servers, so delete its connections record
		if !contains(server, servers) {
			delete(l.serversConnections, server)
		}
	}

	l.servers = servers
}

func (l *LeastConnectionsStrategy) RequestServed(server *models.Server, request *http.Request) {
	defer l.Unlock()
	l.Lock()

	l.decreaseConnections(server)
}
