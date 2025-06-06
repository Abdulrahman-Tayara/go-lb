# Go-LB: Golang Load Balancer

<div align="center">
    <img src="./img.png" alt="Go-LB Logo" width="200"/>
</div>

Go-LB is a simple load balancer solution written in Go, offering a range of features to enhance the distribution of network or application traffic. With support for multiple algorithms, a robust health checker, and convenient configuration options through JSON or YAML files, Go-LB provides a flexible and efficient solution for load balancing in your projects.

## Features

### 1. Multiple Algorithms

Go-LB supports various load balancing algorithms, allowing you to choose the one that best suits your specific requirements. Whether you prefer Round Robin, Least Connections, or a custom algorithm, Go-LB has you covered.

### 2. Health Checker

Ensure the reliability of your backend servers with Go-LB's built-in health checker. It monitors the health of each server and intelligently distributes traffic only to healthy instances, preventing disruptions and improving overall system stability.

### 3. Easy Configuration

Configuring Go-LB is a breeze, it supports both JSON and YAML configuration files. Choose the format that fits your preference and effortlessly define the settings for your load balancing setup.

## Configuration

### JSON
```json
{
  "port": "load balancer port",
  "strategy": "round_robin | weighted_round_robin | sticky_session | least_connections | random", // default round_robin
  "strategy_configs": {
    // Sticky session configs
    "sticky_session_cookie_name": "example",
    "sticky_session_ttl_seconds": 100
  },
  "health_check_interval_seconds": 2,
  "rate_limiter_enabled": true,
  "rate_limit_tokens": 10, // default 10
  "rate_limit_interval_seconds": 1 // default 2
  "servers": [
    {
      "name": "server1",
      "url": "http://localhost:8080",
      "health_url": "/health",
      "weight": 1,
    },
    {
      "name": "server2",
      "url": "http://localhost:8082",
      "health_url": "/health-check",
      "weight": 2
    }
  ],
  "tls_enabled": true, // default false
  "tls_cert_file": "/path/to/cert.pem",
  "tls_key_file": "/path/to/key.pem",
  "log_file": "/path/to/log",
  "log_level": "info | debug | error | warn"
}
```

### Yaml
```yaml
port: "load balancer port"
strategy: "round_robin | weighted_round_robin | sticky_session | least_connections | random" # default round_robin
strategy_configs:
  # Sticky session configs
  sticky_session_cookie_name: "example"
  sticky_session_ttl_seconds: 100
health_check_interval_seconds: 2
rate_limiter_enabled: True
rate_limit_tokens: 10 # default 10
rate_limit_interval_seconds: 10 # default 2
servers:
  - name: "server1"
    url: "http://localhost:8080"
    health_url: "/health"
    weight: 1
  - name: "server2"
    url: "http://localhost:8082"
    health_url: "/health-check"
    weight: 2
tls_enabled: true # default false
tls_cert_file: "/path/on/container/cert.pem"
tls_key_file: "/path/on/container/key.pem"
log_file: "/path/to/log"
log_level: "info | debug | error | warn"
```

## Content Based Routing (CBR) Configuration

### JSON
```json
{
  "port": "load balancer port",
  "strategy": "round_robin | weighted_round_robin | sticky_session | least_connections | random",
  "health_check_interval_seconds": 2,
  "rate_limiter_enabled": true,
  "rate_limit_tokens": 10,
  "rate_limit_interval_seconds": 1,
  "servers": [
    {
      "name": "server1",
      "url": "http://localhost:8080",
      "health_url": "/health",
      "weight": 1
    },
    {
      "name": "server2",
      "url": "http://localhost:8082",
      "health_url": "/health-check",
      "weight": 2
    }
  ],
  "tls_enabled": true,
  "tls_cert_file": "/path/to/cert.pem",
  "tls_key_file": "/path/to/key.pem",
  "routing": {
    "default_server": "server1",
    "rules": [
      {
        "conditions": [
          {
            "path_prefix": "/api/v1 (Optional)",
            "method": "GET | post | Put (Optional)",
            "headers": { // Optional
              "MyHeader": "my-value"
            }
          }
        ],
        "action": {
          "route_to": "server2"
        }
      }
    ]
  }
}
```

### Yaml
```yaml
port: "load balancer port"
strategy: "round_robin | weighted_round_robin | sticky_session | least_connections | random"
health_check_interval_seconds: 2
rate_limiter_enabled: True
rate_limit_tokens: 10
rate_limit_interval_seconds: 10
servers:
  - url: "http://localhost:8080"
    health_url: "/health"
    weight: 1
  - url: "http://localhost:8082"
    health_url: "/health-check"
    weight: 2
tls_enabled: true
tls_cert_file: "/path/on/container/cert.pem"
tls_key_file: "/path/on/container/key.pem"
routing:
  default_server: "server1"
  rules:
    - conditions:
        - path_prefix: "/api/v1"
          method: "GET"
          headers:
            MyHeader: "my-value"
      action:
        route_to: "server2"
```





## Getting Started

### Docker

```bash
git clone https://github.com/Abdulrahman-Tayara/go-lb.git

docker build . -t go-lb:latest

docker run -p <port>:<port> -e CONFIG_FILE=/path/on/container/config.json -v /path/on/host/config.json:/path/on/container/config.json go-lb:latest
```

If the TLS is enabled:

```bash
git clone https://github.com/Abdulrahman-Tayara/go-lb.git

docker build . -t go-lb:latest

docker run -p <port>:<port> -e CONFIG_FILE=/path/on/container/config.json \
  -v /path/on/host/config.json:/path/on/container/config.json \
  -v /path/on/host/cert.pem:/path/on/container/cert.pem \
  -v /path/on/host/key.pem:/path/on/container/key.pem \
  go-lb:latest
```

## Contributing

We welcome contributions to Go-LB! If you have an idea for a new feature or have found a bug, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License, making it open and accessible for a wide range of use cases.