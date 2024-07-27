package sender

import (
	"fmt"
	"log"
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
}

func NewEmailSender(to, from, pass, subj, host, port string) *EmailSender {
	return &EmailSender{
		emailTo:   to,
		emailFrom: from,
		password:  pass,
		subject:   subj,
		host:      host,
		port:      port,
	}
}

func (s *EmailSender) Send(data *models.WeatherData) error {
	auth := smtp.PlainAuth("", s.emailFrom, s.password, s.host)

	to := []string{s.emailTo}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to[0], s.subject, data.ToString()))

	err := smtp.SendMail(fmt.Sprintf("%s:%s", s.host, s.port), auth, s.emailFrom, to, msg)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
