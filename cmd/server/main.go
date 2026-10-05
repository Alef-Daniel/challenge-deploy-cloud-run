package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/infra/viacep"
	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/infra/weatherapi"
	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/infra/web"
	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/usecase"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	getWeather := usecase.NewGetWeatherByZipcodeUseCase(viacep.NewClient(), weatherapi.NewClient(apiKey))
	handler := web.NewWeatherHandler(getWeather)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler.Routes()); err != nil {
		log.Fatal(err)
	}
}
