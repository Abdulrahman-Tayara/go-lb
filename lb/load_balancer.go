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
	servers    []*models.Server
	strategy   lbs.ILoadBalancerStrategy
	cbr        *models.Routing
	cbrEnabled bool

	logger *slog.Logger

	sync.RWMutex
}

func (l *loadBalancer) ServerDown(server *models.Server) {
	defer l.Unlock()

	l.Lock()

	l.servers = slices.DeleteFunc(l.servers, func(s *models.Server) bool {
		return s.Equals(server)
	})

	l.strategy.UpdateServers(l.servers)

	l.logger.Info("server goes down", "server", *server)
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

	l.logger.Info("server comes back", "server", *server)
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

	l.logger.Info("request served in", "server", selectedServer.Url, "time", time.Since(t))
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

	if l.cbrEnabled {
		return l.getServerFromCbr(request)
	}

	return l.strategy.Next(request)
}

func (l *loadBalancer) getServerFromCbr(req *http.Request) *models.Server {
	serverName := l.cbr.Rules.GetRouteTo(&models.RequestProps{
		Method:  req.Method,
		Headers: req.Header,
		Path:    req.URL.Path,
	})

	if serverName == "" {
		serverName = l.cbr.DefaultServer
	}

	return l.serverByName(serverName)
}

func (l *loadBalancer) serverByName(name string) *models.Server {
	for _, server := range l.servers {
		if server.Name == name {
			return server
		}
	}

	return nil
}

type Options struct {
	Logger *slog.Logger
}

func NewLoadBalancer(
	servers []*models.Server,
	cbr *models.Routing,
	strategy lbs.ILoadBalancerStrategy,
	opts ...*Options,
) ILoadBalancer {
	opt := &Options{
		Logger: slog.Default(),
	}
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	strategy.UpdateServers(servers)
	return &loadBalancer{
		servers:    servers,
		strategy:   strategy,
		cbr:        cbr,
		cbrEnabled: cbr != nil && len(cbr.Rules) > 0,
		logger:     opt.Logger,
	}
}
