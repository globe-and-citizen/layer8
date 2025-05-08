package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HealthCheckResponse struct {
	ForwardProxy *Message `json:"forward_proxy,omitempty"`
	ReverseProxy *Message `json:"reverse_proxy,omitempty"`
}

type Message struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message,omitempty"`
	BodyDump   []byte `json:"body_dump,omitempty"`
}

// Fixme: this is the first iteration of the health check handler, needs a more standardized approach like <https://github.com/alexliesenfeld/health>
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// get the backend url from query params "backend_url"
	backendURL := r.URL.Query().Get("backend_url")
	if backendURL == "" {
		health := HealthCheckResponse{
			ForwardProxy: &Message{
				StatusCode: http.StatusBadRequest,
				Message:    "backend_url is required",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
		return
	}

	backendUrlParsed, err := url.JoinPath(backendURL, "l8_health_check")
	if err != nil {
		health := HealthCheckResponse{
			ForwardProxy: &Message{
				StatusCode: http.StatusBadRequest,
				Message:    "backend_url is malformed",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
		return
	}

	req, err := http.NewRequest("GET", backendUrlParsed, nil)
	if err != nil {
		fmt.Println("Useing url", backendURL)
		fmt.Println("Error", err.Error())

		health := HealthCheckResponse{
			ForwardProxy: &Message{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to create request",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
		return
	}

	req.Header.Set("x-tunnel", "true")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		health := HealthCheckResponse{
			ForwardProxy: &Message{StatusCode: http.StatusOK},
			ReverseProxy: &Message{
				StatusCode: http.StatusBadGateway,
				Message:    "failed to send request to reverse proxy",
				BodyDump:   []byte(err.Error()),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error is", err.Error())

		health := HealthCheckResponse{
			ForwardProxy: &Message{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to read response body",
				BodyDump:   []byte(err.Error()),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
		return
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		health := HealthCheckResponse{
			ForwardProxy: &Message{StatusCode: http.StatusOK},
			ReverseProxy: &Message{
				StatusCode: resp.StatusCode,
				Message:    "reverse proxy is not healthy",
				BodyDump:   data,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
		return
	}

	health := &HealthCheckResponse{
		ForwardProxy: &Message{StatusCode: http.StatusOK},
		ReverseProxy: &Message{
			StatusCode: resp.StatusCode,
			Message:    "reverse proxy code is < 500",
			BodyDump:   data,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
