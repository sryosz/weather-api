package main

import (
	"weatherapi/internal/api"
	"weatherapi/internal/config"
	"weatherapi/internal/sender"
)

func main(){
	cfg := config.MustLoad()

	emailSender := sender.NewEmailSender(cfg.EmailTo, cfg.EmailFrom, cfg.Password, cfg.Subject, cfg.Host, cfg.Port)

	wp := api.NewPoller(emailSender, cfg.PollInterval, cfg.Endpoint, cfg.Latitude, cfg.Longitude)
	go func() {
		wp.Start()
	}()

	select {

	}
}