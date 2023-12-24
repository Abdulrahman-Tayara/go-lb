package strategy

import (
	"tayara/go-lb/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobin_GetSelectedServer(t *testing.T) {
	rr := NewRoundRobinStrategy()
	rr.UpdateServers([]*models.Server{
		{
			Url: "http://localhost:8080",
		},
		{
			Url: "http://localhost:8081",
		},
		{
			Url: "http://localhost:8082",
		},
	})

	expectedUrls := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8080",
	}

	for i := 0; i < len(expectedUrls); i++ {
		server := rr.Next(nil)

		assert.Equal(t, expectedUrls[i], server.Url)
	}
}

func TestRoundRobin_GetSelectedServer_EmptyServersList(t *testing.T) {
	rr := NewRoundRobinStrategy()
	rr.UpdateServers([]*models.Server{})

	server := rr.Next(nil)

	assert.Nil(t, server)
}

func TestRoundRobin_GetSelectedServer_UpdateServers(t *testing.T) {
	rr := NewRoundRobinStrategy()

	rr.UpdateServers([]*models.Server{
		{
			Url: "http://localhost:8080",
		},
		{
			Url: "http://localhost:8081",
		},
		{
			Url: "http://localhost:8082",
		},
	})

	runTests := func(expectedUrls []string) {
		if len(expectedUrls) == 0 {
			assert.Nil(t, rr.Next(nil))
			return
		}

		for i := 0; i < len(expectedUrls); i++ {
			server := rr.Next(nil)

			assert.Equal(t, expectedUrls[i], server.Url)
		}
	}

	expectedUrls := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8080",
	}

	runTests(expectedUrls)

	rr.UpdateServers([]*models.Server{})

	runTests([]string{})

	rr.UpdateServers(
		[]*models.Server{
			{
				Url: "http://localhost:8089",
			},
			{
				Url: "http://localhost:8090",
			},
		},
	)

	runTests([]string{
		"http://localhost:8089",
		"http://localhost:8090",
		"http://localhost:8089",
	})
}
