package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"net/http"
	"slices"
	"tayara/go-lb/configs"
	"tayara/go-lb/healthcheck"
	"tayara/go-lb/lb"
	"tayara/go-lb/strategy"
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

	selectedStrategy := strategy.GetLoadBalancerStrategy(cfg.LoadBalancerStrategy)

	loadBalancer := lb.NewLoadBalancer(
		slices.Clone(cfg.Servers),
		selectedStrategy,
	)

	healthChecker := healthcheck.NewHealthChecker(slices.Clone(cfg.Servers))
	healthChecker.Attach(loadBalancer)

	healthChecker.Start(cfg.HealthCheckIntervalSeconds)

	runHTTPServer(cfg, loadBalancer)
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

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
