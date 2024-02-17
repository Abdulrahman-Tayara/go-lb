package strategy

import (
	"math/rand"
	"net/http"
	"tayara/go-lb/models"
	"time"
)

var (
	defaultSessionName = "lbsession"
	defaultTTLSeconds  = 300
)

type StickySessionStrategy struct {
	cfg Configs

	servers []*models.Server
}

func NewStickySessionStrategy(cfg Configs) ILoadBalancerStrategy {
	if cfg.StickySessionCookieName == "" {
		cfg.StickySessionCookieName = defaultSessionName
	}
	if cfg.StickySessionTTLSeconds <= 0 {
		cfg.StickySessionTTLSeconds = defaultTTLSeconds
	}

	return &StickySessionStrategy{
		cfg: cfg,
	}
}

func (s *StickySessionStrategy) Next(request *http.Request) *models.Server {
	cookie, err := request.Cookie(s.cfg.StickySessionCookieName)
	if err != nil || cookie.Value == "" {
		cookieValue := generateSessionID()
		cookie.Value = cookieValue
		request.AddCookie(&http.Cookie{
			Name:     s.cfg.StickySessionCookieName,
			Value:    cookieValue,
			Expires:  time.Now().Add(time.Second * time.Duration(s.cfg.StickySessionTTLSeconds)),
			HttpOnly: true,
		})
	}

	return s.getServer(cookie.Value)
}

func (s *StickySessionStrategy) getServer(sessionId string) *models.Server {
	hash := hashSessionToInt(sessionId)
	index := hash % len(s.servers)
	return s.servers[index]
}

func hashSessionToInt(sessionId string) int {
	hash := 0
	for _, char := range sessionId {
		hash += int(char)
	}
	return hash
}

func (*StickySessionStrategy) RequestServed(server *models.Server, request *http.Request) {
}

func (s *StickySessionStrategy) UpdateServers(servers []*models.Server) {
	s.servers = servers
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateSessionID() string {
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
