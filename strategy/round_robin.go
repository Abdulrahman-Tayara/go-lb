package strategy

import (
	"net/http"
	"sync"
	"tayara/go-lb/models"
)

type RoundRobinStrategy struct {
	index int

	servers []*models.Server

	sync.RWMutex
}

func NewRoundRobinStrategy() ILoadBalancerStrategy {
	return &RoundRobinStrategy{
		index: 0,
	}
}

func (s *RoundRobinStrategy) Next(request *http.Request) *models.Server {
	defer s.RUnlock()

	s.RLock()

	return s.selectServer(s.servers)
}

func (s *RoundRobinStrategy) selectServer(servers []*models.Server) *models.Server {
	if len(servers) == 0 {
		return nil
	}

	s.index = s.index % len(servers)

	selectedServer := s.servers[s.index%len(servers)]

	s.index += 1

	return selectedServer
}

func (s *RoundRobinStrategy) UpdateServers(servers []*models.Server) {
	defer s.Unlock()

	s.Lock()

	s.servers = servers
	s.index = 0
}

func (s *RoundRobinStrategy) RequestServed(server *models.Server, request *http.Request) {
}
