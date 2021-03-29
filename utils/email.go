package utils

import (
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"userstyles.world/config"
)

var auth = sasl.NewPlainClient("", config.EMAIL_ADDRESS, config.EMAIL_PWD)

func SendEmail(to, subject, message string) error {
	r := strings.NewReader(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, config.EMAIL_ADDRESS, subject, message))
	return smtp.SendMail("mail.userstyles.world:587", auth, config.EMAIL_ADDRESS, []string{to}, r)
}
