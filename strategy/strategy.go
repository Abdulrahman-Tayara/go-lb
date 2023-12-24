package strategy

import (
	"net/http"
	"tayara/go-lb/models"
)

type ILoadBalancerStrategy interface {
	Next(request *http.Request) *models.Server

	UpdateServers(servers []*models.Server)
}
