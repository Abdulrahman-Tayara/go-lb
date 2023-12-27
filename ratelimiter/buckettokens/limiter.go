package buckettokens

import (
	"golang.org/x/exp/slog"
	"time"
)

type Configs struct {
	Rate   time.Duration
	Tokens int
}

func NewConfigs(rate time.Duration, tokens int) *Configs {
	cfg := &Configs{
		Rate:   rate,
		Tokens: tokens,
	}

	if tokens <= 0 {
		cfg.Tokens = defaultConfigs.Tokens
	}
	if rate <= 0 {
		cfg.Rate = defaultConfigs.Rate
	}

	return cfg
}

var (
	defaultConfigs = Configs{
		Rate:   time.Second * 2,
		Tokens: 10,
	}
)

// --------------------

type Limiter struct {
	configs *Configs

	data *buckets
}

func NewLimiter(cfg *Configs) *Limiter {
	limiter := &Limiter{
		configs: cfg,
		data:    newBuckets(cfg),
	}

	slog.Info("rate limiter is running,", "configs", *cfg)

	startResetBackground(cfg, limiter.data)

	return limiter
}

type LimitOutput struct {
	Allowed       bool
	RemainingHits int
}

func (l *Limiter) Limit(key string) *LimitOutput {
	b := l.data.get(key)

	if b.shouldReset(l.configs.Rate) {

		l.data.delete(key)

		return &LimitOutput{
			Allowed:       true,
			RemainingHits: b.tokens,
		}
	}

	if !b.hasTokens() {
		return &LimitOutput{Allowed: false, RemainingHits: 0}
	}

	b.take()

	return &LimitOutput{Allowed: true, RemainingHits: b.tokens}
}
