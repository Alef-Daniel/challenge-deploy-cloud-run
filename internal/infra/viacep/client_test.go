package viacep

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/entity"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	c := NewClient()
	c.BaseURL = srv.URL
	return c
}

func TestFindCity(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json/" {
			t.Errorf("path = %s, want /ws/01001000/json/", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"01001-000","localidade":"São Paulo","uf":"SP"}`))
	})

	city, err := c.FindCity(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if city != "São Paulo" {
		t.Errorf("city = %q, want %q", city, "São Paulo")
	}
}

func TestFindCityNotFound(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name: "erro in body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"erro":"true"}`))
			},
		},
		{
			name: "bad request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, tt.handler)
			_, err := c.FindCity(context.Background(), "99999999")
			if !errors.Is(err, entity.ErrZipcodeNotFound) {
				t.Errorf("err = %v, want %v", err, entity.ErrZipcodeNotFound)
			}
		})
	}
}

func TestFindCityUnexpectedStatus(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := c.FindCity(context.Background(), "01001000")
	if err == nil || errors.Is(err, entity.ErrZipcodeNotFound) {
		t.Errorf("err = %v, want unexpected status error", err)
	}
}
