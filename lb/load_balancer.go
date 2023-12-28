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
	"time"
)

const (
	numberOfRetries = 5
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
	if slices.ContainsFunc(l.servers, func(s *models.Server) bool {
		return s.Equals(server)
	}) {
		return
	}

	l.servers = append(l.servers, server)

	l.strategy.UpdateServers(l.servers)

	slog.Info("server comes back", "server", *server)
}

func (l *loadBalancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	selectedServer := l.nextWithRetries(request, numberOfRetries)

	if selectedServer == nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	t := time.Now()

	parsedUrl, _ := url2.Parse(selectedServer.Url)

	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)

	proxy.ServeHTTP(writer, request)

	l.strategy.RequestServed(selectedServer, request)

	slog.Info("request served in", "server", selectedServer.Url, "time", time.Since(t))
}

func (l *loadBalancer) nextWithRetries(request *http.Request, retries int) *models.Server {
	for i := 0; i < retries; i++ {
		server := l.Next(request)

		if server != nil {
			return server
		}

		time.Sleep(time.Second * time.Duration(i+1))
	}

	return nil
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
