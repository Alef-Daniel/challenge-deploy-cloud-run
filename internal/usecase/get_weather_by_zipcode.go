package usecase

import (
	"context"
	"fmt"

	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/entity"
)

type LocationProvider interface {
	FindCity(ctx context.Context, zipcode string) (string, error)
}

type WeatherProvider interface {
	CurrentTemp(ctx context.Context, city string) (float64, error)
}

type GetWeatherByZipcodeOutput struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetWeatherByZipcodeUseCase struct {
	locations LocationProvider
	weather   WeatherProvider
}

func NewGetWeatherByZipcodeUseCase(locations LocationProvider, weather WeatherProvider) *GetWeatherByZipcodeUseCase {
	return &GetWeatherByZipcodeUseCase{locations: locations, weather: weather}
}

func (uc *GetWeatherByZipcodeUseCase) Execute(ctx context.Context, zipcode string) (GetWeatherByZipcodeOutput, error) {
	if err := entity.ValidateZipcode(zipcode); err != nil {
		return GetWeatherByZipcodeOutput{}, err
	}

	city, err := uc.locations.FindCity(ctx, zipcode)
	if err != nil {
		return GetWeatherByZipcodeOutput{}, err
	}

	tempC, err := uc.weather.CurrentTemp(ctx, city)
	if err != nil {
		return GetWeatherByZipcodeOutput{}, fmt.Errorf("current temperature for %s: %w", city, err)
	}

	temp := entity.NewTemperatureFromCelsius(tempC)
	return GetWeatherByZipcodeOutput{
		TempC: temp.Celsius,
		TempF: temp.Fahrenheit,
		TempK: temp.Kelvin,
	}, nil
}
