package helpers

import (
	"log"
)

func HandleError(message string, err interface{}) {
	log.Println("========== Start Error Message ==========")
	log.Println("Message => " + message + ".")
	if err != nil {
		log.Println("Error => ", err)
	}
	log.Println("========== End Of Error Message ==========")
	log.Println()
}
