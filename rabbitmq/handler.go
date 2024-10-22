package rabbitmq

import (
	"fmt"
	"log"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"
)

const (
	SendMailGmail    = "send-mail-gmail"
	SendMailSendgrid = "send-mail-sendgrid"
)

func HandleRequest(data structs.Request) string {
	log.Println("Accepted command -> ", data.Command)

	// sample command
	switch data.Command {
	case SendMailGmail:
		var req structs.Mail
		err := helpers.MarshalUnmarshal(data.Data, &req)
		if err != nil {
			return err.Error()
		}
		fmt.Println("data:", req)

		err = helpers.SendMailGmail(req.Body, req.Subject, req.TargetEmail)
		if err != nil {
			return err.Error()
		}
		return ""
	case SendMailSendgrid:
		var req structs.Mail
		err := helpers.MarshalUnmarshal(data.Data, &req)
		if err != nil {
			return err.Error()
		}

		err = helpers.SendMailSendgrid(req.Body, req.Subject, req.TargetName, req.TargetEmail)
		if err != nil {
			return err.Error()
		}
		return ""
	default:
		return "command not found"
	}

}
