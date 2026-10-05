package usecase

import (
	"context"
	"errors"
	"testing"

	"go-1/internal/entity"
)

type fakeLocationProvider struct {
	city   string
	err    error
	called bool
}

func (f *fakeLocationProvider) FindCity(ctx context.Context, zipcode string) (string, error) {
	f.called = true
	return f.city, f.err
}

type fakeWeatherProvider struct {
	tempC   float64
	err     error
	called  bool
	gotCity string
}

func (f *fakeWeatherProvider) CurrentTemp(ctx context.Context, city string) (float64, error) {
	f.called = true
	f.gotCity = city
	return f.tempC, f.err
}

func TestExecuteSuccess(t *testing.T) {
	locations := &fakeLocationProvider{city: "São Paulo"}
	weather := &fakeWeatherProvider{tempC: 28.5}
	uc := NewGetWeatherByZipcodeUseCase(locations, weather)

	got, err := uc.Execute(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := GetWeatherByZipcodeOutput{TempC: 28.5, TempF: 83.3, TempK: 301.65}
	if got != want {
		t.Errorf("Execute() = %+v, want %+v", got, want)
	}
	if weather.gotCity != "São Paulo" {
		t.Errorf("weather called with %q, want %q", weather.gotCity, "São Paulo")
	}
}

func TestExecuteInvalidZipcode(t *testing.T) {
	locations := &fakeLocationProvider{}
	weather := &fakeWeatherProvider{}
	uc := NewGetWeatherByZipcodeUseCase(locations, weather)

	_, err := uc.Execute(context.Background(), "123")
	if !errors.Is(err, entity.ErrInvalidZipcode) {
		t.Errorf("err = %v, want %v", err, entity.ErrInvalidZipcode)
	}
	if locations.called || weather.called {
		t.Error("providers should not be called for an invalid zipcode")
	}
}

func TestExecuteZipcodeNotFound(t *testing.T) {
	locations := &fakeLocationProvider{err: entity.ErrZipcodeNotFound}
	weather := &fakeWeatherProvider{}
	uc := NewGetWeatherByZipcodeUseCase(locations, weather)

	_, err := uc.Execute(context.Background(), "99999999")
	if !errors.Is(err, entity.ErrZipcodeNotFound) {
		t.Errorf("err = %v, want %v", err, entity.ErrZipcodeNotFound)
	}
	if weather.called {
		t.Error("weather provider should not be called when zipcode is not found")
	}
}

func TestExecuteWeatherError(t *testing.T) {
	weatherErr := errors.New("weather unavailable")
	uc := NewGetWeatherByZipcodeUseCase(&fakeLocationProvider{city: "São Paulo"}, &fakeWeatherProvider{err: weatherErr})

	_, err := uc.Execute(context.Background(), "01001000")
	if !errors.Is(err, weatherErr) {
		t.Errorf("err = %v, want %v", err, weatherErr)
	}
}
