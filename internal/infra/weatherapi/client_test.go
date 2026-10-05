package weatherapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	c := NewClient("test-key")
	c.BaseURL = srv.URL
	return c
}

func TestCurrentTemp(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/current.json" {
			t.Errorf("path = %s, want /v1/current.json", r.URL.Path)
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Errorf("key = %q, want %q", got, "test-key")
		}
		if got := r.URL.Query().Get("q"); got != "São Paulo, Brazil" {
			t.Errorf("q = %q, want %q", got, "São Paulo, Brazil")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"location":{"name":"Sao Paulo"},"current":{"temp_c":28.5,"temp_f":83.3}}`))
	})

	got, err := c.CurrentTemp(context.Background(), "São Paulo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 28.5 {
		t.Errorf("temp = %v, want 28.5", got)
	}
}

func TestCurrentTempUnexpectedStatus(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	if _, err := c.CurrentTemp(context.Background(), "São Paulo"); err == nil {
		t.Error("expected error, got nil")
	}
}
