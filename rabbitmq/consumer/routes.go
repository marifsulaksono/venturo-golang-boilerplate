package consumer

import (
	"encoding/json"
	"log"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"
	"strings"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func ConsumerRouter(data *amqp.Delivery, db *gorm.DB) {
	switch data.MessageId {
	case helpers.KEY_EXPORT_DATA_USER:
		body := strings.Split(string(data.Body), "|||")

		var request structs.Job
		err := json.Unmarshal([]byte(body[1]), &request)
		if err != nil {
			log.Println("ERROR:", err)
		}
		ExportUserDataToCsv(&request, body[0], db)
	}
}
