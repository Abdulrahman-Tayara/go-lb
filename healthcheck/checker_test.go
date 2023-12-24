package healthcheck

import (
	"github.com/stretchr/testify/assert"
	"tayara/go-lb/models"
	"testing"
)

func TestHealthChecker_didStatusChanged(t *testing.T) {
	server1 := &models.Server{
		Url: "http://localhost:8080",
	}
	server2 := &models.Server{
		Url: "http://localhost:8081",
	}

	checker := HealthChecker{
		serversStatues: map[*models.Server]bool{
			server1: true,
		},
	}

	assert.False(t, checker.didStatusChanged(server1, true))
	assert.True(t, checker.didStatusChanged(server1, false))

	assert.True(t, checker.didStatusChanged(server2, false))
}

func TestHealthChecker_changeServerStatus(t *testing.T) {
	server1 := &models.Server{
		Url: "http://localhost:8080",
	}
	server2 := &models.Server{Url: "http://localhost:8081"}

	checker := HealthChecker{
		serversStatues: map[*models.Server]bool{
			server1: true,
		},
	}

	checker.changeServerStatus(server1, false)

	assert.False(t, checker.serversStatues[server1])

	checker.changeServerStatus(server2, true)

	assert.True(t, checker.serversStatues[server2])
}

func TestHealthChecker_checkAndNotify_serverDown(t *testing.T) {
	server1 := &models.Server{
		Url: "http://localhost:8080",
	}
	server2 := &models.Server{Url: "http://localhost:8081"}

	observer := &mockObserver{}
	checker := HealthChecker{
		serversStatues: map[*models.Server]bool{
			server1: true,
			server2: true,
		},
	}
	checker.Attach(observer)

	alwaysDown := func(s *models.Server) bool {
		return false
	}

	checker.checkAndNotify(
		server1,
		alwaysDown,
		func(server *models.Server, b bool) bool {
			return true
		},
	)

	assert.Equal(t, 1, observer.serversDownCount)

	checker.checkAndNotify(
		server1,
		alwaysDown,
		func(server *models.Server, b bool) bool {
			return false
		},
	)

	assert.Equal(t, 1, observer.serversDownCount)
}

func TestHealthChecker_checkAndNotify_serverUp(t *testing.T) {
	server1 := &models.Server{
		Url: "http://localhost:8080",
	}
	server2 := &models.Server{Url: "http://localhost:8081"}

	observer := &mockObserver{}
	checker := HealthChecker{
		serversStatues: map[*models.Server]bool{
			server1: false,
			server2: true,
		},
	}
	checker.Attach(observer)

	alwaysUp := func(s *models.Server) bool {
		return true
	}

	checker.checkAndNotify(
		server1,
		alwaysUp,
		func(server *models.Server, b bool) bool {
			return true
		},
	)

	assert.Equal(t, 1, observer.serversUpCount)

	checker.checkAndNotify(
		server1,
		alwaysUp,
		func(server *models.Server, b bool) bool {
			return false
		},
	)

	assert.Equal(t, 1, observer.serversUpCount)
}

func TestHealthChecker_checkAndNotify_serverUpThenDownTheUp(t *testing.T) {
	server1 := &models.Server{
		Url: "http://localhost:8080",
	}
	server2 := &models.Server{Url: "http://localhost:8081"}

	observer := &mockObserver{}
	checker := HealthChecker{
		serversStatues: map[*models.Server]bool{
			server1: true,
			server2: true,
		},
	}
	checker.Attach(observer)

	alwaysDown := func(s *models.Server) bool {
		return false
	}
	alwaysUp := func(s *models.Server) bool {
		return true
	}

	checker.checkAndNotify(
		server1,
		alwaysDown,
		checker.didStatusChanged,
	)

	assert.Equal(t, 1, observer.serversDownCount)
	assert.Equal(t, 0, observer.serversUpCount)

	checker.checkAndNotify(
		server1,
		alwaysUp,
		checker.didStatusChanged,
	)

	assert.Equal(t, 1, observer.serversDownCount)
	assert.Equal(t, 1, observer.serversUpCount)
}
