package lb

import (
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"tayara/go-lb/models"
	lbs "tayara/go-lb/strategy"
)

// ILoadBalancer interface
type ILoadBalancer interface {
	http.Handler
}

// -----------------------

type loadBalancer struct {
	servers  []*models.Server
	strategy lbs.ILoadBalancerStrategy
}

func (l *loadBalancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	selectedServer := l.Next(request)

	parsedUrl, _ := url2.Parse(selectedServer.Url)

	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)

	proxy.ServeHTTP(writer, request)
}

func (l *loadBalancer) Next(request *http.Request) *models.Server {
	return l.strategy.Next(request)
}

func NewLoadBalancer(servers []*models.Server, strategy lbs.ILoadBalancerStrategy) ILoadBalancer {
	strategy.UpdateServers(servers)
	return &loadBalancer{
		servers:  servers,
		strategy: strategy,
	}
}
