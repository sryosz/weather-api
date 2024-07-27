package api

import (
	"encoding/json"
	"fmt"
	"log"
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
}

func NewPoller(sender Sender, interval time.Duration, endpoint string, latitude, longitude float64) *WPoller {
	return &WPoller{
		endpoint:     endpoint,
		pollInterval: interval,
		latitude:     latitude,
		longitude:    longitude,
		closeCh:      make(chan struct{}),
		sender:       sender,
	}
}

func (wp *WPoller) Start() {
	ticker := time.NewTicker(wp.pollInterval)

	for {
		select {
		case <-ticker.C:
			data, err := getWeatherResults(wp.endpoint, wp.latitude, wp.longitude)
			if err != nil {
				log.Fatal(err)
			}
			err = wp.sender.Send(data)
			if err != nil {
				return
			}
		case <-wp.closeCh:
			fmt.Println("shutdown gracefully")
			return
		}
	}

}

func (wp *WPoller) Close() {
	close(wp.closeCh)
}

func getWeatherResults(endpoint string, lat, long float64) (*models.WeatherData, error) {
	uri := fmt.Sprintf("%s?latitude=%f&longitude=%f&hourly=temperature_2m", endpoint, lat, long)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var data models.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	return &data, nil
}
