package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-crud-rnd/config"
	"simple-crud-rnd/structs"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/gomail.v2"
)

func SendMailGmail(body, subject, to string) error {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	sender := os.Getenv("SMTP_SENDER_NAME")
	email := os.Getenv("SMTP_SENDER_EMAIL")
	pass := os.Getenv("SMTP_PASS_KEY")

	senderName := fmt.Sprintf("%s <%s>", sender, email)

	// set message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	// set messsage dialer
	dialer := gomail.NewDialer(host, port, email, pass)

	// dial and send message
	err := dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func SendMailGmailWithAsync(body, subject, to, jobId string) {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	sender := os.Getenv("SMTP_SENDER_NAME")
	email := os.Getenv("SMTP_SENDER_EMAIL")
	pass := os.Getenv("SMTP_PASS_KEY")

	senderName := fmt.Sprintf("%s <%s>", sender, email)

	// set message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	// set messsage dialer
	dialer := gomail.NewDialer(host, port, email, pass)

	// dial and send message
	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Printf("=> send email to %s is failure. error: %v", to, err.Error())
		updateErr := config.DB.Model(&structs.Job{}).
			Where("id = ?", jobId).
			Update("status", structs.JOB_FAILED).Error
		if updateErr != nil {
			log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", jobId, updateErr)
		}
	}

	updateErr := config.DB.Model(&structs.Job{}).
		Where("id = ?", jobId).
		Update("status", structs.JOB_SUCCESS).Error
	if updateErr != nil {
		log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", jobId, updateErr)
	}
	log.Printf("send email to %s is successfully", to)
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

func SendMailSendgridWithAsync(content, subject, targetName, targetEmail, jobId string) {
	from := mail.NewEmail(os.Getenv("SENDGRID_SENDER_NAME"), os.Getenv("SENDGRID_SENDER_EMAIL"))
	to := mail.NewEmail(targetName, targetEmail)
	message := mail.NewSingleEmail(from, subject, to, content, fmt.Sprintf("<strong>%s</strong>", content))
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil || response.StatusCode == http.StatusInternalServerError {
		log.Printf("=> send email to %s is failure. error: %v", to, err.Error())
		updateErr := config.DB.Model(&structs.Job{}).
			Where("id = ?", jobId).
			Update("status", structs.JOB_FAILED).Error
		if updateErr != nil {
			log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", jobId, updateErr)
		}
	}

	updateErr := config.DB.Model(&structs.Job{}).
		Where("id = ?", jobId).
		Update("status", structs.JOB_SUCCESS).Error
	if updateErr != nil {
		log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", jobId, updateErr)
	}
	log.Printf("send email to %s is successfully", to)
}
