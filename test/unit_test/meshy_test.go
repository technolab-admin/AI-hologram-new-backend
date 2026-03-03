package unit_test

import (
	"AI-HOLOGRAM-NEW-BACKEND/internal/http/handlers"
	"AI-HOLOGRAM-NEW-BACKEND/internal/logger"
	"AI-HOLOGRAM-NEW-BACKEND/internal/meshy"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/// Init ///

func init() {
	logger.Warn = log.New(io.Discard, "WARN: ", log.LstdFlags)
}

/// Fakes and Mocks ///

type fakeJobRunner struct {
	called bool
}

func (f *fakeJobRunner) Run(prompt string, req *meshy.TextTo3DRequest) {
	f.called = true
}

type fakeRoundTripper struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.fn(req)
}

/// Tests ///

func Test_HttpHandler_Generate(t *testing.T) {
	tests := []struct {
		name           string
		jsonTest       string
		expectedStatus int
	}{
		{name: "invalid json syntax", jsonTest: `{"name: "test"`, expectedStatus: http.StatusBadRequest},
		{name: "not a json", jsonTest: `not a json`, expectedStatus: http.StatusBadRequest},
		{name: "valid json", jsonTest: `{"name": "test"}`, expectedStatus: http.StatusAccepted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRunner := &fakeJobRunner{}
			h := handlers.NewHttpHandler(fakeRunner)
			json := strings.NewReader(tt.jsonTest)
			req := httptest.NewRequest(http.MethodPost, "/generate", json)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			h.Generate(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusBadRequest && fakeRunner.called {
				t.Fatalf("jobRunner.Run should not be called on bad JSON")
			}
		})
	}

}

func Test_MeshyClient_CreateJob(t *testing.T) {
	tests := []struct {
		name          string
		roundTripFunc func(*http.Request) (*http.Response, error)
		apiKey        string
		url           string
		expectedError bool
		expectedID    string
	}{{
		name: "success response",
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"result":"1234"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		},
		apiKey:        "correctKey",
		url:           "http://fake.url/jobs",
		expectedError: false,
		expectedID:    "1234",
	},
		{
			name: "failed network request",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
			apiKey:        "correctKey",
			url:           "http://fake.url/jobs",
			expectedError: true,
			expectedID:    "",
		},
		{
			name: "invalid JSON response",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				body := `not a json`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			},
			apiKey:        "correctKey",
			url:           "http://fake.url/jobs",
			expectedError: true,
			expectedID:    "",
		}, {
			name: "invalid API key",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				body := `{"error":"invalid api key}`
				return &http.Response{
					StatusCode: 401,
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			},
			apiKey:        "wrongKey",
			url:           "http://fake.url/jobs",
			expectedError: true,
			expectedID:    "",
		}, {
			name: "invalid URL",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("dial tcp: lookup wrongUrl: No such host")
			},
			apiKey:        "correctKey",
			url:           "http://wrong.url/jobs",
			expectedError: true,
			expectedID:    "",
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{
				Transport: &fakeRoundTripper{fn: tt.roundTripFunc},
			}
			c := meshy.NewClient(tt.apiKey, tt.url, httpClient)

			resp, err := c.CreateJob(tt.url, map[string]string{"test": "data"})

			if tt.expectedError && err == nil {
				t.Fatalf("expected error, got none")
			}
			if !tt.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.expectedError && resp.ResultID != tt.expectedID {
				t.Fatalf("expected ID %s, got %s", tt.expectedID, resp.ResultID)
			}
		})
	}
}
