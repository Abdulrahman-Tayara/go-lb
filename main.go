package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
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
	logger         *slog.Logger
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

	logger = getLogger(cfg)

	logger.Info("configs were loaded", "configs", *cfg)

	selectedStrategy := strategy.GetLoadBalancerStrategy(cfg.LoadBalancerStrategy, cfg.StrategyConfigs)

	loadBalancer := lb.NewLoadBalancer(
		slices.Clone(cfg.Servers),
		&cfg.Routing,
		selectedStrategy,
		&lb.Options{
			Logger: getLogger(cfg),
		},
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

func getLogger(cfg *configs.Configs) *slog.Logger {
	mapLogLevel := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	logLevel := mapLogLevel[strings.ToLower(cfg.LogLevel)]
	if cfg.LogFile != "" {
		file, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: logLevel}))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

func runHTTPServer(cfg *configs.Configs, handler http.Handler) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: handler,
	}

	logger.Info("Load Balancer is running",
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
