package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteCondition_DoesMatch(t *testing.T) {
	type input struct {
		routeCondition RouteCondition
		req            RequestProps
	}

	tests := []struct {
		name  string
		input input
		want  bool
	}{
		{
			name: "match",
			want: true,
			input: input{
				req: RequestProps{
					Path:   "/api/v1",
					Method: "GET",
					Headers: map[string][]string{
						"UserAgent": {"Mobile"},
					},
				},
				routeCondition: RouteCondition{
					PathPrefix: "/api/v1",
					Method:     "GET",
					Headers: map[string]string{
						"UserAgent": "Mobile",
					},
				},
			},
		},
		{
			name: "not match - path",
			want: false,
			input: input{
				req: RequestProps{
					Path:   "/api/v1",
					Method: "GET",
					Headers: map[string][]string{
						"UserAgent": {"Mobile"},
					},
				},
				routeCondition: RouteCondition{
					PathPrefix: "/api/v2",
				},
			},
		},
		{
			name: "not match - missed header",
			want: false,
			input: input{
				req: RequestProps{
					Path:   "/api/v1",
					Method: "GET",
					Headers: map[string][]string{
						"UserAgent": {"Mobile"},
					},
				},
				routeCondition: RouteCondition{
					Headers: map[string]string{
						"UserAgent":    "Mobile",
						"CustomHeader": "custom",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.routeCondition.DoesMatch(&tt.input.req)

			assert.Equal(t, tt.want, got)

		})
	}
}
