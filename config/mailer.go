package config

import (
	"gopkg.in/gomail.v2"
	"log"
)

func (cfg *Config) InitMailerConn() {
	log.Println("Trying to init mailer dialer . . . .")
	mailerSingleton.Do(func() {
		Mailer = gomail.NewDialer(
			cfg.MailHost,
			cfg.MailPort,
			cfg.MailUsername,
			cfg.MailPassword,
		)
		log.Println("Mailer dialer created . . . .")
	})
}
