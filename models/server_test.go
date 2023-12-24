package models

import "testing"

func TestGetHealthUrl(t *testing.T) {
	tests := []struct {
		name  string
		input *Server
		want  string
	}{
		{
			name: "Empty health url",
			input: &Server{
				Url:       "http://localhost:8080",
				HealthUrl: "",
			},
			want: "http://localhost:8080/health",
		},
		{
			name: "Absolute health url",
			input: &Server{
				Url:       "http://localhost:8080",
				HealthUrl: "http://localhost:8080/health",
			},
			want: "http://localhost:8080/health",
		},
		{
			name: "Relative health url",
			input: &Server{
				Url:       "http://localhost:8080",
				HealthUrl: "/health-endpoint",
			},
			want: "http://localhost:8080/health-endpoint",
		},
		{
			name: "Relative health url without / prefix",
			input: &Server{
				Url:       "http://localhost:8080",
				HealthUrl: "health-endpoint",
			},
			want: "http://localhost:8080/health-endpoint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.GetHealthUrl(); got != tt.want {
				t.Errorf("Server.GetHealthUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHost(t *testing.T) {
	tests := []struct {
		name  string
		input *Server
		want  string
	}{
		{
			name: "Http url",
			input: &Server{
				Url: "http://localhost:8080",
			},
			want: "localhost:8080",
		},
		{
			name: "Https url",
			input: &Server{
				Url: "https://localhost:8080",
			},
			want: "localhost:8080",
		},
		{
			name: "Url with endpoint and query params",
			input: &Server{
				Url: "https://localhost:8080/hello?test=one",
			},
			want: "localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.Host(); got != tt.want {
				t.Errorf("Server.Hostname() = %v, want %v", got, tt.want)
			}
		})
	}
}
