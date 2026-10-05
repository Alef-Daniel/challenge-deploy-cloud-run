package entity

import "math"

type Temperature struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func NewTemperatureFromCelsius(c float64) Temperature {
	return Temperature{
		Celsius:    round(c),
		Fahrenheit: round(CelsiusToFahrenheit(c)),
		Kelvin:     round(CelsiusToKelvin(c)),
	}
}

func CelsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func CelsiusToKelvin(c float64) float64 {
	return c + 273.15
}

func round(v float64) float64 {
	return math.Round(v*100) / 100
}
