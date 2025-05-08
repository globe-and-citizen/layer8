package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		scenario                    string
		reverseProxyUrl             func() (string, func()) // url and Close fn
		expectedCodeForwardProxy    int
		expectedMessageForwardProxy string
		expectedCodeReverseProxy    int
		expectedMessageReverseProxy string
	}{
		{
			scenario: "backend_url is empty",
			reverseProxyUrl: func() (string, func()) {
				return "", nil
			},
			expectedCodeForwardProxy:    http.StatusBadRequest,
			expectedMessageForwardProxy: "backend_url is required",
		},
		{
			scenario: "failed to send request",
			reverseProxyUrl: func() (string, func()) {
				return uuid.NewString(), nil
			},
			expectedCodeForwardProxy:    http.StatusOK,
			expectedCodeReverseProxy:    http.StatusBadGateway,
			expectedMessageReverseProxy: "failed to send request to reverse proxy",
		},
		{
			scenario: "internal server error from reverse proxy",
			reverseProxyUrl: func() (string, func()) {
				internalServerStatusSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))

				return internalServerStatusSrv.URL, internalServerStatusSrv.Close
			},
			expectedCodeForwardProxy:    http.StatusOK,
			expectedCodeReverseProxy:    http.StatusInternalServerError,
			expectedMessageReverseProxy: "reverse proxy is not healthy",
		},
		{
			scenario: "ok from reverse proxy",
			reverseProxyUrl: func() (string, func()) {
				okServerStatusSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))

				return okServerStatusSrv.URL, okServerStatusSrv.Close
			},
			expectedCodeForwardProxy: http.StatusOK,
			expectedCodeReverseProxy: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			var (
				responseRecorder    = httptest.NewRecorder()
				healthCheckResponse HealthCheckResponse
			)

			request := httptest.NewRequest("GET", "/health_check", nil)
			reverseProxyUrl, closeFn := tt.reverseProxyUrl()
			if reverseProxyUrl != "" {
				values := request.URL.Query()
				values.Add("backend_url", reverseProxyUrl)
				request.URL.RawQuery = values.Encode()
			}

			// cleanup test server
			defer func() {
				if closeFn != nil {
					closeFn()
				}
			}()

			HealthCheck(responseRecorder, request)
			assert.NoError(t, json.NewDecoder(responseRecorder.Body).Decode(&healthCheckResponse), "Failed to decode response")

			if tt.expectedCodeForwardProxy > 0 && healthCheckResponse.ForwardProxy != nil {
				assert.Equal(t, tt.expectedCodeForwardProxy, healthCheckResponse.ForwardProxy.StatusCode, "Expected status code %d", tt.expectedCodeForwardProxy)
			}

			if tt.expectedMessageForwardProxy != "" && healthCheckResponse.ForwardProxy != nil {
				assert.Equal(t, tt.expectedMessageForwardProxy, healthCheckResponse.ForwardProxy.Message, "Expected message %s", tt.expectedMessageForwardProxy)
			}

			if tt.expectedCodeReverseProxy > 0 && healthCheckResponse.ReverseProxy != nil {
				assert.Equal(t, tt.expectedCodeReverseProxy, healthCheckResponse.ReverseProxy.StatusCode, "Expected status code %d", tt.expectedCodeReverseProxy)
			}

			if tt.expectedMessageReverseProxy != "" && healthCheckResponse.ReverseProxy != nil {
				assert.Equal(t, tt.expectedMessageReverseProxy, healthCheckResponse.ReverseProxy.Message, "Expected message %s", tt.expectedMessageReverseProxy)
			}
		})
	}
}
