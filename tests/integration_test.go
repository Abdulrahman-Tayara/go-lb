package tests

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"slices"
	"tayara/go-lb/healthcheck"
	lb2 "tayara/go-lb/lb"
	"tayara/go-lb/models"
	"tayara/go-lb/strategy"
	"testing"
	"time"
)

func serverHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/health" {
		w.WriteHeader(http.StatusOK)
	} else if r.RequestURI == "/endpoint1" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("welcome"))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func doRequest(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("expected 200 ok")
	}

	return nil
}

func TestE2EIntegration(t *testing.T) {
	server1 := &models.Server{
		Url:       "http://localhost:9000",
		HealthUrl: "/health",
	}
	server2 := &models.Server{
		Url:       "http://localhost:9001",
		HealthUrl: "/health",
	}

	server1Http := setupServers(func(server *models.Server) http.Handler {
		return http.HandlerFunc(serverHandler)
	}, server1)
	server2Http := setupServers(func(server *models.Server) http.Handler {
		return http.HandlerFunc(serverHandler)
	}, server2)

	defer server1Http.Close()
	defer server2Http.Close()

	serversList := []*models.Server{server1, server2}

	lb := lb2.NewLoadBalancer(slices.Clone(serversList), nil, strategy.NewRoundRobinStrategy())
	lbServer := &models.Server{
		Url: "http://localhost:9065",
	}

	_ = setupServers(func(server *models.Server) http.Handler {
		return lb
	}, lbServer)

	healthCheckInterval := 1
	healthChecker := healthcheck.NewHealthChecker(slices.Clone(serversList))
	healthChecker.Attach(lb)
	healthChecker.Start(int(healthCheckInterval))

	time.Sleep(time.Second * time.Duration(healthCheckInterval))

	endpointUrl, _ := url.JoinPath(lbServer.Url, "endpoint1")

	// Try 100 requests on the load balancer with two servers
	for i := 0; i < 100; i++ {
		err := doRequest(endpointUrl)
		assert.NoError(t, err)
	}

	// Stop the first server
	server1Http.Close()

	// Wait the health checker to re-check
	time.Sleep(time.Duration(healthCheckInterval) * time.Second)

	// Try 100 requests on the load balancer with one server
	for i := 0; i < 100; i++ {
		err := doRequest(endpointUrl)
		assert.NoError(t, err)
	}

	// Stop the second server
	server2Http.Close()

	// Wait the health checker to re-check
	time.Sleep(time.Duration(healthCheckInterval) * time.Second)

	// Try to make request to the load balancer with no servers
	err := doRequest(endpointUrl)

	assert.Error(t, err)

	// Re-setup the first server
	server1Http = setupServers(func(server *models.Server) http.Handler {
		return http.HandlerFunc(serverHandler)
	}, server1)

	// Wait the health checker to re-check
	time.Sleep(time.Duration(healthCheckInterval) * time.Second)

	// Try 100 requests on the load balancer with one server
	for i := 0; i < 100; i++ {
		err := doRequest(endpointUrl)
		assert.NoError(t, err)
	}
}
