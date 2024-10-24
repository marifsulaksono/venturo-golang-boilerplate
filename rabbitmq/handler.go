package rabbitmq

import (
	"fmt"
	"log"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"

	"gorm.io/gorm"
)

const (
	SendMailGmail    = "send-mail-gmail"
	SendMailSendgrid = "send-mail-sendgrid"
)

func HandleRequest(data structs.Request, db *gorm.DB) string {
	log.Println("Accepted command -> ", data.Command)

	// Handle perintah berdasarkan command yang diterima
	switch data.Command {
	case SendMailGmail:
		var req structs.Mail
		// Unmarshal data dari request menjadi struct Mail
		err := helpers.MarshalUnmarshal(data.Data, &req)
		if err != nil {
			return err.Error()
		}

		// Kirim email menggunakan Gmail
		err = helpers.SendMailGmail(req.Body, req.Subject, req.TargetEmail)

		if err != nil {
			// Jika terjadi error, update status job menjadi gagal
			log.Printf("Error sending email via Gmail for job ID: %s. Error: %v", data.JobId, err)
			updateErr := db.Model(&structs.Job{}).
				Where("id = ?", data.JobId).
				Update("status", structs.JOB_FAILED).Error
			if updateErr != nil {
				log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", data.JobId, updateErr)
			}
			return err.Error()
		}

		log.Println("Done send email via Gmail... Error:", err)

		log.Println("Job ID:", data.JobId)

		// // Jika sukses, update status job menjadi sukses
		updateErr := db.Model(&structs.Job{}).
			Where("id = ?", data.JobId).
			Update("status", structs.JOB_SUCCESS).Error
		if updateErr != nil {
			log.Printf("Failed to update job status to success for job ID: %s. Error: %v", data.JobId, updateErr)
			return updateErr.Error()
		}

		fmt.Println("Email sent successfully via Gmail for job ID:", data.JobId)
		return "Success"

	case SendMailSendgrid:
		var req structs.Mail
		// Unmarshal data dari request menjadi struct Mail
		err := helpers.MarshalUnmarshal(data.Data, &req)
		if err != nil {
			return err.Error()
		}

		fmt.Println("Sendgrid JOB ID:", data.JobId)

		// Kirim email menggunakan Sendgrid
		err = helpers.SendMailSendgrid(req.Body, req.Subject, req.TargetName, req.TargetEmail)
		if err != nil {
			// Jika terjadi error, update status job menjadi gagal
			log.Printf("Error sending email via Sendgrid for job ID: %s. Error: %v", data.JobId, err)
			updateErr := config.DB.Model(&structs.Job{}).
				Where("id = ?", data.JobId).
				Update("status", structs.JOB_FAILED).Error
			if updateErr != nil {
				log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", data.JobId, updateErr)
			}
			return err.Error()
		}

		// Jika sukses, update status job menjadi sukses
		updateErr := config.DB.Model(&structs.Job{}).
			Where("id = ?", data.JobId).
			Update("status", structs.JOB_SUCCESS).Error
		if updateErr != nil {
			log.Printf("Failed to update job status to success for job ID: %s. Error: %v", data.JobId, updateErr)
			return updateErr.Error()
		}

		fmt.Println("Email sent successfully via Sendgrid for job ID:", data.JobId)
		return "Success"

	default:
		// Jika command tidak dikenali, kembalikan pesan error
		log.Printf("Command not found: %s", data.Command)
		return "command not found"
	}
}
