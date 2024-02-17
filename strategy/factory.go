package strategy

const (
	RoundRobinLoadBalancerStrategy       = "round_robin"
	RandomLoadBalancerStrategy           = "random"
	LeastConnectionsLoadBalancerStrategy = "least_connections"
	StickySessionLoadBalancerStrategy    = "stick_session"
)

func GetLoadBalancerStrategy(strategy string, cfg Configs) ILoadBalancerStrategy {
	switch strategy {
	case RoundRobinLoadBalancerStrategy:
		return NewRoundRobinStrategy()
	case RandomLoadBalancerStrategy:
		return NewRandomStrategy()
	case LeastConnectionsLoadBalancerStrategy:
		return NewLeastConnectionsStrategy()
	case StickySessionLoadBalancerStrategy:
		return NewStickySessionStrategy(cfg)
	default:
		return NewRoundRobinStrategy()
	}
}
