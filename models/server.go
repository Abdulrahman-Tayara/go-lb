package models

import (
	"net/url"
	"strings"
)

const (
	defaultHealthUrl = "/health"
)

type Server struct {
	Name      string `mapstructure:"name" json:"name" yaml:"name"`
	Url       string `mapstructure:"url" json:"url" yaml:"url"`
	HealthUrl string `mapstructure:"health_url" json:"health_url" yaml:"health_url"`
	Weight    int    `mapstructure:"weight" json:"weight" yaml:"weight"`

	host              string
	absoluteHealthUrl string
}

func (s *Server) GetUrl() string {
	return s.Url
}

func (s *Server) GetHealthUrl() string {
	if s.absoluteHealthUrl != "" {
		return s.absoluteHealthUrl
	}

	var healthUrl string
	if s.HealthUrl == "" {
		healthUrl = defaultHealthUrl
	} else {
		healthUrl = s.HealthUrl
	}

	if strings.HasPrefix(healthUrl, "http") {
		return healthUrl
	}

	healthUrl, _ = url.JoinPath(s.Url, healthUrl)

	s.absoluteHealthUrl = healthUrl

	return s.absoluteHealthUrl
}

func (s *Server) Host() string {
	if s.host != "" {
		return s.host
	}

	u, _ := url.Parse(s.Url)
	s.host = u.Host

	return s.host
}

func (s *Server) Equals(other *Server) bool {
	return s.Url == other.Url
}
