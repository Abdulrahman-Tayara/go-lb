package tests

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	url2 "net/url"
	"reflect"
	"strings"
	lb2 "tayara/go-lb/lb"
	"tayara/go-lb/models"
	"tayara/go-lb/strategy"
	"testing"
)

const (
	successMessage = "the request was forwarded successfully"
)

func TestLoadBalancerRequestForwarding(t *testing.T) {
	destServer := models.Server{
		Url: "http://localhost:8087",
	}

	requests := map[string]http.Request{
		"/endpoint1": {
			Header: map[string][]string{
				"Custom-Header": {"Value1"},
			},
			Method: http.MethodPost,
			Body:   io.NopCloser(strings.NewReader("hello world")),
		},
	}
	// I'm storing the bodies in separate objects, because the original request's body is read while sending the request to the load balancer
	requestsBodies := map[string]io.Reader{
		"/endpoint1": strings.NewReader("hello world"),
	}

	destServerHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		wantedRequest := requests[request.RequestURI]
		wantedRequestBody := requestsBodies[request.RequestURI]

		if request.Method != wantedRequest.Method {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		for k, v := range wantedRequest.Header {
			if !reflect.DeepEqual(request.Header.Values(k), v) {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		requestBody, _ := io.ReadAll(request.Body)
		wantedRequestBodyBytes, _ := io.ReadAll(wantedRequestBody)

		if string(requestBody) != string(wantedRequestBodyBytes) {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(successMessage))
	})

	destHttpServer := setupServers(func(server *models.Server) http.Handler {
		return destServerHandler
	}, &destServer)

	defer destHttpServer.Close()

	lbServer := models.Server{
		Url: "http://localhost:8076",
	}
	lb := lb2.NewLoadBalancer(
		[]*models.Server{&destServer},
		nil,
		strategy.NewRoundRobinStrategy(),
	)

	lbHttpServer := setupServers(func(server *models.Server) http.Handler {
		return lb
	}, &lbServer)

	defer lbHttpServer.Close()

	for k, r := range requests {
		// Setup
		url, _ := url2.JoinPath(lbServer.Url, k)
		request, _ := http.NewRequest(r.Method, url, r.Body)
		request.Header = r.Header

		res, err := http.DefaultClient.Do(request)

		// Validate
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		bytes, _ := io.ReadAll(res.Body)

		assert.Equal(t, successMessage, string(bytes))
	}
}
