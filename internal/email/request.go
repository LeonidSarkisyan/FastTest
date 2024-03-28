package email

import (
	"App/pkg/systems"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/smtp"
	"os"
)

type EmailClient struct {
	cfg *systems.AppConfig
	smtp.Auth
}

func NewEmailClient(auth smtp.Auth, cfg *systems.AppConfig) *EmailClient {
	return &EmailClient{Auth: auth, cfg: cfg}
}

func (c *EmailClient) SendCodeToEmail(emailTo string, code int64) error {
	from := os.Getenv("SMTP_EMAIL")
	to := []string{emailTo}

	subject := "Регистрация аккаунта"
	body := fmt.Sprintf(`
	<h1>Здравствуйте!</h1>
	<p>Спасибо за регистрацию на Фаст Тест РФ.</p>
	<p>Чтобы завершить регистрацию перейдите по ссылке ниже:</p>
	`)

	if c.cfg.Debug {
		body += fmt.Sprintf(`<a href="http://localhost/auth/confirm/%d">Подтвердить свою учётную запись</a>`, code)
	} else {
		body += fmt.Sprintf(`<a href="https://фаст-тест.рф/auth/confirm/%d">Подтвердить свою учётную запись</a>`, code)
	}

	message := []byte("From: " + from + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body + "\r\n")

	log.Info().Msg("отправляем письмо...")

	log.Print(c.cfg.SMTP.Host+":"+c.cfg.SMTP.Port, c.Auth, from, to, message)

	err := smtp.SendMail(
		c.cfg.SMTP.Host+":"+c.cfg.SMTP.Port, c.Auth, from, to, message,
	)

	if err != nil {
		log.Err(err).Send()
		return err
	}

	log.Info().Msg("письмо отправлено!")

	return nil
}
