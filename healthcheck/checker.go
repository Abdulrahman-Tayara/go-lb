package healthcheck

import (
	"net/http"
	"sync"
	"tayara/go-lb/models"
	"time"
)

type IHealthChecker interface {
	IObservable
	Start(intervalSeconds int)
}

type HealthChecker struct {
	servers []*models.Server

	observers []IObserver

	// serversStatues if the value is true so the server is healthy otherwise the server is down
	serversStatues map[*models.Server]bool

	httpClient *http.Client

	sync.RWMutex
}

func (h *HealthChecker) Attach(observer IObserver) {
	h.observers = append(h.observers, observer)
}

func (h *HealthChecker) Detach(observer IObserver) {
	for i, item := range h.observers {
		if item == observer {
			h.observers = append(h.observers[:i], h.observers[i+1:]...)
			break
		}
	}
}

func (h *HealthChecker) Start(intervalSeconds int) {
	go func() {
		for {
			h.Check()
			time.Sleep(time.Duration(intervalSeconds) * time.Second)
		}
	}()
}

func (h *HealthChecker) Check() {
	for _, server := range h.servers {
		go func(s *models.Server) {
			h.checkAndNotify(s, h.isHealthy, h.didStatusChanged)
		}(server)
	}
}

func (h *HealthChecker) checkAndNotify(
	server *models.Server,
	healthCheckerFunc func(*models.Server) bool,
	didStatusChangedFunc func(*models.Server, bool) bool,
) {
	isHealthy := healthCheckerFunc(server)
	statusChanged := didStatusChangedFunc(server, isHealthy)

	if statusChanged {
		if isHealthy {
			h.notifyServerUp(server)
		} else {
			h.notifyServerDown(server)
		}
	}

	h.changeServerStatus(server, isHealthy)
}

func (h *HealthChecker) isHealthy(server *models.Server) bool {
	res, err := h.httpClient.Get(server.GetHealthUrl())
	return err == nil && res.StatusCode == http.StatusOK
}

func (h *HealthChecker) didStatusChanged(server *models.Server, isHealthy bool) bool {
	defer h.RUnlock()

	h.RLock()

	current, exists := h.serversStatues[server]
	if !exists {
		return true
	}
	return current != isHealthy
}

func (h *HealthChecker) notifyServerUp(server *models.Server) {
	for _, observer := range h.observers {
		observer.ServerUp(server)
	}
}

func (h *HealthChecker) notifyServerDown(server *models.Server) {
	for _, observer := range h.observers {
		observer.ServerDown(server)
	}
}

func (h *HealthChecker) changeServerStatus(server *models.Server, status bool) {
	defer h.Unlock()

	h.Lock()

	h.serversStatues[server] = status
}

func NewHealthChecker(servers []*models.Server) IHealthChecker {
	return &HealthChecker{
		servers:        servers,
		httpClient:     http.DefaultClient,
		serversStatues: make(map[*models.Server]bool),
	}
}
