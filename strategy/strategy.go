package strategy

import (
	"net/http"
	"tayara/go-lb/models"
)

type Configs struct {
	StickySessionCookieName string `mapstructure:"sticky_session_cookie_name" json:"sticky_session_cookie_name" yaml:"sticky_session_cookie_name"`
	StickySessionTTLSeconds int    `mapstructure:"sticky_session_ttl_seconds" json:"sticky_session_ttl_seconds" yaml:"sticky_session_ttl_seconds"`
}

type ILoadBalancerStrategy interface {
	Next(request *http.Request) *models.Server

	UpdateServers(servers []*models.Server)

	RequestServed(server *models.Server, request *http.Request)
}
