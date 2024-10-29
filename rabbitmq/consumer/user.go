package consumer

import (
	"fmt"
	"log"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"

	"gorm.io/gorm"
)

func ExportUserDataToCsv(job *structs.Job, to string, db *gorm.DB) {
	fmt.Printf("data job: %v\nto: %s\n", job, to)

	err := helpers.SendMailGmail("Berikut adalah file export data users", job.JobName, to, job.Payload)
	if err != nil {
		log.Printf("=> send email to %s is failure. error: %v", to, err.Error())
		updateErr := config.DB.Model(&structs.Job{}).
			Where("id = ?", job.ID).
			Update("status", structs.JOB_FAILED).Error
		if updateErr != nil {
			log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", job.ID, updateErr)
		}
	}

	updateErr := config.DB.Model(&structs.Job{}).
		Where("id = ?", job.ID).
		Update("status", structs.JOB_SUCCESS).Error
	if updateErr != nil {
		log.Printf("Failed to update job status to failed for job ID: %s. Error: %v", job.ID, updateErr)
	}
	log.Printf("send email to %s is successfully", to)
}
