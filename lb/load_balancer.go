package lb

import (
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"slices"
	"sync"
	"tayara/go-lb/healthcheck"
	"tayara/go-lb/models"
	lbs "tayara/go-lb/strategy"
)

// ILoadBalancer interface
type ILoadBalancer interface {
	http.Handler
	healthcheck.IObserver
}

// -----------------------

type loadBalancer struct {
	servers  []*models.Server
	strategy lbs.ILoadBalancerStrategy

	sync.RWMutex
}

func (l *loadBalancer) ServerDown(server *models.Server) {
	defer l.Unlock()

	l.Lock()

	l.servers = slices.DeleteFunc(l.servers, func(s *models.Server) bool {
		return s.Equals(server)
	})

	l.strategy.UpdateServers(l.servers)

	slog.Info("server goes down", "server", *server)
}

func (l *loadBalancer) ServerUp(server *models.Server) {
	defer l.Unlock()

	l.Lock()

	// Prevent the duplication
	l.servers = slices.DeleteFunc(l.servers, func(s *models.Server) bool {
		return s.Equals(server)
	})

	l.servers = append(l.servers, server)

	l.strategy.UpdateServers(l.servers)

	slog.Info("server comes back", "server", *server)
}

func (l *loadBalancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	selectedServer := l.Next(request)

	if selectedServer == nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	parsedUrl, _ := url2.Parse(selectedServer.Url)

	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)

	proxy.ServeHTTP(writer, request)
}

func (l *loadBalancer) Next(request *http.Request) *models.Server {
	defer l.RUnlock()

	l.RLock()

	return l.strategy.Next(request)
}

func NewLoadBalancer(servers []*models.Server, strategy lbs.ILoadBalancerStrategy) ILoadBalancer {
	strategy.UpdateServers(servers)
	return &loadBalancer{
		servers:  servers,
		strategy: strategy,
	}
}
