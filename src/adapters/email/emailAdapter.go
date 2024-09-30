package email

import (
	"fmt"
	"net/smtp"

	"github.com/Dialosoft/src/app/config"
)

func SendEmail(to []string, subject, body string, config config.GeneralConfig) error {
	headers := make(map[string]string)
	headers["From"] = config.FromAddress
	headers["To"] = to[0]
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", config.MailUsername, config.MailPassword, config.SMTPHost)

	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)

	err := smtp.SendMail(addr, auth, config.FromAddress, to, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
