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
  "strategy": "round_robin | random | least_connections",
  "health_check_interval_seconds": 2,
  "rate_limiter_enabled": true,
  "rate_limit_tokens": 10, // default 10
  "rate_limit_interval_seconds": 1 // default 2
  "servers": [
    {
      "url": "http://localhost:8080",
      "health_url": "/health"
    },
    {
      "url": "http://localhost:8082",
      "health_url": "/health-check"
    }
  ]
}
```

### Yaml
```yaml
port: "load balancer port"
strategy: "round_robin | random | least_connections"
health_check_interval_seconds: 2
rate_limiter_enabled: True
rate_limit_tokens: 10 # default 10
rate_limit_interval_seconds: 10 # default 2
servers:
  - url: "http://localhost:8080"
    health_url: "/health"
  - url: "http://localhost:8082"
    health_url: "/health-check"
```

## Getting Started

### Docker

```bash
git clone https://github.com/Abdulrahman-Tayara/go-lb.git

docker build . -t go-lb:latest

docker run -p <port>:<port> -e CONFIG_FILE=/path/on/container/config.json -v /path/on/host/config.json:/path/on/container/config.json go-lb:latest
```

## Contributing

We welcome contributions to Go-LB! If you have an idea for a new feature or have found a bug, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License, making it open and accessible for a wide range of use cases.