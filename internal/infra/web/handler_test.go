package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-1/internal/entity"
	"go-1/internal/usecase"
)

type fakeGetWeather struct {
	output     usecase.GetWeatherByZipcodeOutput
	err        error
	gotZipcode string
}

func (f *fakeGetWeather) Execute(ctx context.Context, zipcode string) (usecase.GetWeatherByZipcodeOutput, error) {
	f.gotZipcode = zipcode
	return f.output, f.err
}

func TestGetWeatherSuccess(t *testing.T) {
	uc := &fakeGetWeather{output: usecase.GetWeatherByZipcodeOutput{TempC: 28.5, TempF: 83.3, TempK: 301.65}}
	rec := httptest.NewRecorder()
	NewWeatherHandler(uc).Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/weather/01001000", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if uc.gotZipcode != "01001000" {
		t.Errorf("zipcode = %q, want %q", uc.gotZipcode, "01001000")
	}

	var got usecase.GetWeatherByZipcodeOutput
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got != uc.output {
		t.Errorf("response = %+v, want %+v", got, uc.output)
	}
}

func TestGetWeatherErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "invalid zipcode", err: entity.ErrInvalidZipcode, wantStatus: http.StatusBadRequest, wantBody: "invalid zipcode"},
		{name: "zipcode not found", err: entity.ErrZipcodeNotFound, wantStatus: http.StatusNotFound, wantBody: "can not find zipcode"},
		{name: "unexpected error", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantBody: "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			NewWeatherHandler(&fakeGetWeather{err: tt.err}).Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/weather/123", nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
