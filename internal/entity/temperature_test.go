package entity

import "testing"

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius float64
		want    float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
		{28.5, 83.3},
	}

	for _, tt := range tests {
		got := round(CelsiusToFahrenheit(tt.celsius))
		if got != tt.want {
			t.Errorf("CelsiusToFahrenheit(%v) = %v, want %v", tt.celsius, got, tt.want)
		}
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius float64
		want    float64
	}{
		{0, 273.15},
		{-273.15, 0},
		{28.5, 301.65},
	}

	for _, tt := range tests {
		got := round(CelsiusToKelvin(tt.celsius))
		if got != tt.want {
			t.Errorf("CelsiusToKelvin(%v) = %v, want %v", tt.celsius, got, tt.want)
		}
	}
}

func TestNewTemperatureFromCelsius(t *testing.T) {
	got := NewTemperatureFromCelsius(28.5)
	want := Temperature{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.65}

	if got != want {
		t.Errorf("NewTemperatureFromCelsius(28.5) = %+v, want %+v", got, want)
	}
}
