package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"weatherapi/internal/models"
)

type Sender interface {
	Send(data *models.WeatherData) error
}

type WPoller struct {
	endpoint     string
	pollInterval time.Duration
	latitude     float64
	longitude    float64
	closeCh      chan struct{}
	sender       Sender
	log          *slog.Logger
}

func NewPoller(log *slog.Logger, sender Sender, interval time.Duration, endpoint string, latitude, longitude float64) *WPoller {
	return &WPoller{
		endpoint:     endpoint,
		pollInterval: interval,
		latitude:     latitude,
		longitude:    longitude,
		closeCh:      make(chan struct{}),
		sender:       sender,
		log:          log,
	}
}

func (wp *WPoller) Start() {
	const op = "sender.Start"
	log := wp.log.With("op", op)

	ticker := time.NewTicker(wp.pollInterval)

	log.Info("starting Weather Poller")

	for {
		select {
		case <-ticker.C:
			data, err := getWeatherResults(wp.endpoint, wp.latitude, wp.longitude)
			if err != nil {
				log.Error(err.Error())
				return
			}
			err = wp.sender.Send(data)
			if err != nil {
				log.Error(err.Error())
				return
			}
		case <-wp.closeCh:
			log.Info("shutdown gracefully")
			return
		}
	}

}

func (wp *WPoller) Close() {
	const op = "sender.Close"
	log := wp.log.With("op", op)

	close(wp.closeCh)

	log.Info("closed channel")
}

func getWeatherResults(endpoint string, lat, long float64) (*models.WeatherData, error) {
	uri := fmt.Sprintf("%s?latitude=%f&longitude=%f&hourly=temperature_2m", endpoint, lat, long)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var data models.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
