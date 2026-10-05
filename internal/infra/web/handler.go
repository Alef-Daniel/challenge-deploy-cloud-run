package web

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/entity"
	"github.com/Alef-Daniel/challenge-deploy-cloud-run/internal/usecase"
)

type GetWeatherByZipcode interface {
	Execute(ctx context.Context, zipcode string) (usecase.GetWeatherByZipcodeOutput, error)
}

type WeatherHandler struct {
	getWeather GetWeatherByZipcode
}

func NewWeatherHandler(getWeather GetWeatherByZipcode) *WeatherHandler {
	return &WeatherHandler{getWeather: getWeather}
}

func (h *WeatherHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{cep}", h.GetWeather)
	return mux
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	output, err := h.getWeather.Execute(r.Context(), r.PathValue("cep"))
	switch {
	case errors.Is(err, entity.ErrInvalidZipcode):
		writeMessage(w, http.StatusBadRequest, entity.ErrInvalidZipcode.Error())
		return
	case errors.Is(err, entity.ErrZipcodeNotFound):
		writeMessage(w, http.StatusNotFound, entity.ErrZipcodeNotFound.Error())
		return
	case err != nil:
		log.Printf("get weather: %v", err)
		writeMessage(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func writeMessage(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(message))
}
