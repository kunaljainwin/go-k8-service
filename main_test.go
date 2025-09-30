package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupServer() http.Handler {
	mux := http.NewServeMux()
	base := "/api/v1"

	mux.HandleFunc(base+"/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	mux.HandleFunc(base+"/healthcheck", HealthHandler)
	mux.HandleFunc(base+"/healthz", HealthHandler)
	mux.HandleFunc(base+"/readyz", HealthHandler)
	mux.HandleFunc(base+"/echo", EchoHandler)

	return enableCORS(mux)
}

// -------------------- Tests --------------------

func TestRootHandler(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/")
	if err != nil {
		t.Fatalf("could not send GET: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello World!" {
		t.Errorf("expected Hello World!, got %s", string(body))
	}
}

func TestHealthHandler(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	endpoints := []string{"/api/v1/healthcheck", "/api/v1/healthz", "/api/v1/readyz"}
	for _, ep := range endpoints {
		resp, err := http.Get(ts.URL + ep)
		if err != nil {
			t.Fatalf("GET %s failed: %v", ep, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s returned %d", ep, resp.StatusCode)
		}

		var hr HealthResponse
		if err := json.NewDecoder(resp.Body).Decode(&hr); err != nil {
			t.Errorf("%s response not valid JSON: %v", ep, err)
		}
		if hr.Status != "running" {
			t.Errorf("%s expected status=running, got %s", ep, hr.Status)
		}
		if time.Since(hr.Timestamp) > time.Second {
			t.Errorf("%s timestamp too old", ep)
		}
	}
}

func TestEchoHandler(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	msg := `{"message":"hello"}`
	resp, err := http.Post(ts.URL+"/api/v1/echo", "application/json", bytes.NewBuffer([]byte(msg)))
	if err != nil {
		t.Fatalf("POST echo failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != msg {
		t.Errorf("expected %s, got %s", msg, string(body))
	}
}

func TestCORSPreflight(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	req, _ := http.NewRequest("OPTIONS", ts.URL+"/api/v1/echo", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("preflight request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if origin := resp.Header.Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("expected CORS *, got %s", origin)
	}
}
