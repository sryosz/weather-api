package sender

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"weatherapi/internal/models"
)

type EmailSender struct {
	emailTo   string
	emailFrom string
	password  string
	subject   string
	host      string
	port      string
	log       *slog.Logger
}

func NewEmailSender(log *slog.Logger, to, from, pass, subj, host, port string) *EmailSender {
	return &EmailSender{
		emailTo:   to,
		emailFrom: from,
		password:  pass,
		subject:   subj,
		host:      host,
		port:      port,
		log:       log,
	}
}

func (s *EmailSender) Send(data *models.WeatherData) error {
	const op = "sender.Send"
	log := s.log.With("op", op)

	auth := smtp.PlainAuth("", s.emailFrom, s.password, s.host)

	to := []string{s.emailTo}

	parsedData, err := formatMsg(data)
	if err != nil {
		return err
	}

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to[0], s.subject, parsedData))

	err = smtp.SendMail(fmt.Sprintf("%s:%s", s.host, s.port), auth, s.emailFrom, to, msg)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("sent msg from %s to %s", s.emailFrom, s.emailTo))

	return nil
}

func formatMsg(data *models.WeatherData) (string, error) {

	msg := fmt.Sprintf("Weather at %.2fm\n", data.Elevation)

	timeSlice, ok := data.Hourly["time"].([]any)
	if !ok {
		return "", fmt.Errorf("error extracting time from map")
	}
	tempSlice, ok := data.Hourly["temperature_2m"].([]any)
	if !ok {
		return "", fmt.Errorf("error extracting temperature from map")
	}

	for i := 0; i < len(timeSlice) && i < len(tempSlice); i++ {
		time, ok := timeSlice[i].(string)
		if !ok {
			return "", fmt.Errorf("error converting time element to string")
		}
		temp, ok := tempSlice[i].(float64)
		if !ok {
			return "", fmt.Errorf("error converting temperature element to float64")
		}
		msg += fmt.Sprintf("%s: %.1f\n", time, temp)
	}

	return msg, nil
}
