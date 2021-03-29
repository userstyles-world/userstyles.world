package utils

import (
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"userstyles.world/config"
)

var auth = sasl.NewPlainClient("", config.EMAIL_ADDRESS, config.EMAIL_PWD)

func SendEmail(to, message string) error {
	r := strings.NewReader(message)
	return smtp.SendMail("mail.userstyles.world:587", auth, config.EMAIL_ADDRESS, []string{to}, r)
}
