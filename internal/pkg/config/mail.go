package config

import (
	xtremecore "github.com/globalxtreme/go-core/v2"
	"os"
)

var (
	SMTPMail xtremecore.MailConf
)

func InitMail() {
	SMTPMail = xtremecore.MailConf{
		Host:     os.Getenv("MAIL_HOST"),
		Port:     os.Getenv("MAIL_PORT"),
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
	}
}
