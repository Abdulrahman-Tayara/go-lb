package tests

import (
	"net"
	"net/http"
	"net/http/httptest"
	"tayara/go-lb/models"
)

type testsServers []*httptest.Server

func (s testsServers) Close() {
	for _, server := range s {
		server.Close()
	}
}

func setupServers(serverHandler func(server *models.Server) http.Handler, servers ...*models.Server) testsServers {
	var httpTestServers []*httptest.Server

	for _, server := range servers {
		s := server
		testServer := httptest.NewUnstartedServer(serverHandler(s))
		testServer.Listener, _ = net.Listen("tcp", server.Host())
		testServer.Start()

		httpTestServers = append(httpTestServers, testServer)
	}

	return httpTestServers
}
