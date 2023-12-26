package strategy

import (
	"github.com/stretchr/testify/assert"
	"tayara/go-lb/models"
	"testing"
)

func TestRandom_Next(t *testing.T) {

	servers := []*models.Server{
		{
			Url: "http://localhost:8080",
		},
		{
			Url: "http://localhost:8010",
		},
	}

	s := NewRandomStrategy()
	s.UpdateServers(servers)

	for i := 0; i < 100; i++ {
		server := s.Next(nil)
		assert.NotNil(t, server)
	}
}

func TestRandom_Next_EmptyList(t *testing.T) {

	var servers []*models.Server

	s := NewRandomStrategy()
	s.UpdateServers(servers)

	for i := 0; i < 100; i++ {
		server := s.Next(nil)
		assert.Nil(t, server)
	}
}
