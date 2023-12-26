package strategy

import (
	"math/rand"
	"net/http"
	"sync"
	"tayara/go-lb/models"
	"time"
)

type RandomStrategy struct {
	servers []*models.Server

	sync.RWMutex
}

func NewRandomStrategy() ILoadBalancerStrategy {
	return &RandomStrategy{}
}

func (r *RandomStrategy) Next(request *http.Request) *models.Server {
	defer r.RUnlock()

	r.RLock()

	if len(r.servers) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	randomServer := rand.Intn(len(r.servers))

	return r.servers[randomServer]
}

func (r *RandomStrategy) UpdateServers(servers []*models.Server) {
	defer r.Unlock()

	defer r.Lock()

	r.servers = servers
}
