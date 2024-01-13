package models

import (
	"slices"
	"strings"
)

type Routing struct {
	Rules         RouteRules `mapstructure:"rules" json:"rules" yaml:"rules"`
	DefaultServer string     `mapstructure:"default_server" json:"default_server" yaml:"default_server"`
}

type RequestProps struct {
	Path    string
	Headers map[string][]string
	Method  string
}

type RouteRule struct {
	Conditions []RouteCondition `mapstructure:"conditions" json:"conditions" yaml:"conditions"`
	Action     RouteAction      `mapstructure:"action" json:"action" yaml:"action"`
}

type RouteAction struct {
	RouteTo string `mapstructure:"route_to" json:"route_to" yaml:"route_to"`
}

type RouteCondition struct {
	PathPrefix string            `mapstructure:"path_prefix" json:"path_prefix" yaml:"path_prefix"`
	Headers    map[string]string `mapstructure:"headers" json:"headers" yaml:"headers"`
	Method     string            `mapstructure:"method" json:"method" yaml:"method"`
}

type RouteRules []RouteRule

func (rc *RouteCondition) DoesMatch(req *RequestProps) bool {
	if rc.PathPrefix != "" && !strings.HasPrefix(req.Path, rc.PathPrefix) {
		return false
	}
	if rc.Method != "" && strings.ToLower(rc.Method) != strings.ToLower(req.Method) {
		return false
	}
	for k, v := range rc.Headers {
		if !slices.Contains(req.Headers[k], v) {
			return false
		}
	}
	return true
}

func (rr *RouteRule) DoesMatch(req *RequestProps) bool {
	for _, condition := range rr.Conditions {
		if !condition.DoesMatch(req) {
			return false
		}
	}
	return true
}

func (rr *RouteRules) GetRouteTo(req *RequestProps) string {
	for _, rule := range *rr {
		if rule.DoesMatch(req) {
			return rule.Action.RouteTo
		}
	}

	return ""
}
