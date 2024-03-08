package strategy

import (
	"net/http"
	"slices"
	"sync"
	"tayara/go-lb/models"
)

type WeightedRoundRobinStrategy struct {
	weightsBucket map[string]int

	servers []*models.Server

	sync.RWMutex
}

func NewWeightedRoundRobinStrategy() ILoadBalancerStrategy {
	return &WeightedRoundRobinStrategy{
		weightsBucket: make(map[string]int),
	}
}

func (s *WeightedRoundRobinStrategy) Next(request *http.Request) *models.Server {
	defer s.RUnlock()

	s.RLock()

	return s.selectServer(s.servers)
}

func (s *WeightedRoundRobinStrategy) selectServer(servers []*models.Server) *models.Server {
	if len(servers) == 0 {
		return nil
	}

	s.ensureWightsBucketValid()

	selectedServer := s.findMaxServerWeight()

	return selectedServer
}

func (s *WeightedRoundRobinStrategy) findMaxServerWeight() *models.Server {
	var selectedServerName string
	var maxWeight = -1000

	for serverName, tokens := range s.weightsBucket {
		if tokens > maxWeight {
			maxWeight = tokens
			selectedServerName = serverName
		}
	}

	if index := slices.IndexFunc(s.servers, func(s *models.Server) bool {
		return s.Name == selectedServerName
	}); index >= 0 && index < len(s.servers) {
		return s.servers[index]
	} else {
		return nil
	}
}

func (s *WeightedRoundRobinStrategy) ensureWightsBucketValid() {
	for _, tokens := range s.weightsBucket {
		if tokens > 0 {
			return
		}
	}

	s.resetBucket()
}

func (s *WeightedRoundRobinStrategy) UpdateServers(servers []*models.Server) {
	defer s.Unlock()

	s.Lock()

	s.servers = servers

	s.resetBucket()
}

func (s *WeightedRoundRobinStrategy) resetBucket() {
	clear(s.weightsBucket)

	for _, server := range s.servers {
		s.weightsBucket[server.Name] = server.Weight
	}
}

func (s *WeightedRoundRobinStrategy) RequestServed(server *models.Server, request *http.Request) {
	if tokens, ok := s.weightsBucket[server.Name]; ok {
		s.weightsBucket[server.Name] = max(tokens-1, 0)
	}
}
