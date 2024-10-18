package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/gomail.v2"
)

func SendMailGmail(message, subject, target string) error {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	sender := os.Getenv("SMTP_SENDER_NAME")
	email := os.Getenv("SMTP_SENDER_EMAIL")
	pass := os.Getenv("SMTP_PASS_KEY")

	senderName := fmt.Sprintf("%s <%s>", sender, email)

	// set message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderName)
	mailer.SetHeader("To", target)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	// set messsage dialer
	dialer := gomail.NewDialer(host, port, email, pass)

	// dial and send message
	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Println("Error send email:", err)
		return err
	}

	return nil
}

func SendMailSendgrid(content, subject, targetName, targetEmail string) error {
	from := mail.NewEmail(os.Getenv("SENDGRID_SENDER_NAME"), os.Getenv("SENDGRID_SENDER_EMAIL"))
	to := mail.NewEmail(targetName, targetEmail)
	message := mail.NewSingleEmail(from, subject, to, content, fmt.Sprintf("<strong>%s</strong>", content))
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil || response.StatusCode == http.StatusInternalServerError {
		log.Println("Error send email:", err)
		return err
	}

	return nil
}
