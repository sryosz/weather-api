package models

import "fmt"

type WeatherData struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

func (d *WeatherData) ToString() string{
	return fmt.Sprintf("%v%v", d.Elevation, d.Hourly)
}
