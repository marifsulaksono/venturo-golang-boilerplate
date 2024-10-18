package helpers

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMailGmail(message, subject, target string) error {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	sender := os.Getenv("SMTP_SENDER_NAME")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	senderName := fmt.Sprintf("%s <%s>", sender, user)

	// set message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderName)
	mailer.SetHeader("To", target)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	// set messsage dialer
	dialer := gomail.NewDialer(host, port, user, pass)

	// dial and send message
	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Println("Error send email:", err)
		return err
	}

	return nil
}
