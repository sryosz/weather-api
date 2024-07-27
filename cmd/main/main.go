package main

import (
	"log/slog"
	"os"
	"weatherapi/internal/api"
	"weatherapi/internal/config"
	"weatherapi/internal/sender"
)

func main(){
	cfg := config.MustLoad()

	log := setupLogger()

	emailSender := sender.NewEmailSender(log, cfg.EmailTo, cfg.EmailFrom, cfg.Password, cfg.Subject, cfg.Host, cfg.Port)

	wp := api.NewPoller(log, emailSender, cfg.PollInterval, cfg.Endpoint, cfg.Latitude, cfg.Longitude)
	go func() {
		wp.Start()
	}()

	select {

	}
}

func setupLogger() *slog.Logger{
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}