package email

import (
	"backend/seed-savers/config"
	"fmt"
	"net/smtp"
)

func initAuth() smtp.Auth {
	return smtp.PlainAuth(
		"",
		config.Envs.Email,
		config.Envs.EmailPassword,
		config.Envs.Hostsmtp,
	)
}

func SendMail(reciver string,  msg []byte) error {
	auth := initAuth()

	return smtp.SendMail(fmt.Sprintf("%v:%v", config.Envs.Hostsmtp, "587"), auth, config.Envs.Email, []string{reciver},msg)
}
