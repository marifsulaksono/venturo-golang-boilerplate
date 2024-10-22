package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

func RandomByte(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}

func JSONEncode(data interface{}) string {
	jsonResult, _ := json.Marshal(data)

	return string(jsonResult)
}

func MarshalUnmarshal(param interface{}, result interface{}) error {
	paramByte, err := json.Marshal(param)
	if err != nil {
		log.Println("Error marshal", err.Error())
		return err
	}

	err = json.Unmarshal(paramByte, &result)
	if err != nil {
		log.Println("Error unmarshal", err.Error())
		return err
	}

	return nil
}
