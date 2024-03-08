package main

import (
	"flag"
	"fmt"
	"net/http"
	"slices"
	"tayara/go-lb/configs"
	"tayara/go-lb/healthcheck"
	"tayara/go-lb/lb"
	"tayara/go-lb/ratelimiter/buckettokens"
	"tayara/go-lb/strategy"
	"tayara/go-lb/utils"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

var (
	configFilePath = ""
)

func init() {
	flag.StringVar(&configFilePath, "config", "config.json", "config file path")
	flag.Parse()
}

func main() {
	if configFilePath == "" {
		panic("config file path is empty")
	}

	slog.Info("loading the configs from", "filepath", configFilePath)

	cfg, err := configs.LoadConfigs(configFilePath)
	if err != nil {
		panic(errors.Wrap(err, "error while loading the config file"))
	}

	slog.Info("configs were loaded", "configs", *cfg)

	selectedStrategy := strategy.GetLoadBalancerStrategy(cfg.LoadBalancerStrategy, cfg.StrategyConfigs)

	loadBalancer := lb.NewLoadBalancer(
		slices.Clone(cfg.Servers),
		&cfg.Routing,
		selectedStrategy,
	)

	healthChecker := healthcheck.NewHealthChecker(slices.Clone(cfg.Servers))
	healthChecker.Attach(loadBalancer)
	healthChecker.Start(cfg.HealthCheckIntervalSeconds)

	var httpHandler http.Handler = loadBalancer

	if cfg.RateLimiterEnabled {
		httpHandler = buckettokens.LimiterMiddleware(
			loadBalancer,
			buckettokens.NewConfigs(
				time.Second*time.Duration(cfg.RateLimitIntervalSeconds),
				cfg.RateLimitTokens,
			),
		)
	}

	runHTTPServer(cfg, httpHandler)
}

func runHTTPServer(cfg *configs.Configs, handler http.Handler) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: handler,
	}

	slog.Info("Load Balancer is running",
		"address", server.Addr,
		"configs", *cfg,
	)

	if cfg.TLSEnabled {
		if !utils.IsFileExist(cfg.TLSCertPath) {
			panic(fmt.Errorf("TLS cert filepath doesn't exist %v", cfg.TLSCertPath))
		}
		if !utils.IsFileExist(cfg.TLSKeyPath) {
			panic(fmt.Errorf("TLS key filepath doesn't exist %v", cfg.TLSKeyPath))
		}

		if err := server.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath); err != nil {
			panic(err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}
}
