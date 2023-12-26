package strategy

const (
	RoundRobinLoadBalancerStrategy       = "round_robin"
	RandomLoadBalancerStrategy           = "random"
	LeastConnectionsLoadBalancerStrategy = "least_connections"
)

func GetLoadBalancerStrategy(strategy string) ILoadBalancerStrategy {
	switch strategy {
	case RoundRobinLoadBalancerStrategy:
		return NewRoundRobinStrategy()
	case RandomLoadBalancerStrategy:
		return NewRandomStrategy()
	case LeastConnectionsLoadBalancerStrategy:
		return NewLeastConnectionsStrategy()
	default:
		return NewRoundRobinStrategy()
	}
}
