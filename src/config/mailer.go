package config

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/Sahil2k07/Blog-App-Go/src/utils"
)

func Mailer(toEmail string, otp string) error {

	mailHost := os.Getenv("MAIL_HOST")
	mailUser := os.Getenv("MAIL_USER")
	mailPass := os.Getenv("MAIL_PASS")

	auth := smtp.PlainAuth("", mailUser, mailPass, mailHost)

	subject := "Subject: Email Verification\n"
	contentType := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := utils.AuthEmail(otp)

	message := []byte(subject + contentType + body)

	smtpAddr := mailHost + ":587"

	err := smtp.SendMail(smtpAddr, auth, mailUser, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
