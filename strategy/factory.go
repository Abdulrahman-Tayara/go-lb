package strategy

const (
	RoundRobinLoadBalancerStrategy         = "round_robin"
	RandomLoadBalancerStrategy             = "random"
	LeastConnectionsLoadBalancerStrategy   = "least_connections"
	StickySessionLoadBalancerStrategy      = "stick_session"
	WeightedRoundRobinLoadBalancerStrategy = "weighted_round_robin"
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
	case WeightedRoundRobinLoadBalancerStrategy:
		return NewWeightedRoundRobinStrategy()
	default:
		return NewRoundRobinStrategy()
	}
}
