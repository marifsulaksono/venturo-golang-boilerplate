package structs

import (
	"time"

	"github.com/go-playground/validator"
)

type JSONResponse struct {
	ResponseCode    int         `json:"response_code"`
	ResponseMessage string      `json:"response_message"`
	Message         string      `json:"message,omitempty"`
	Data            interface{} `json:"data"`
}

type Delete struct {
	DeletedAt time.Time `json:"deleted_at"`
	DeletedBy int       `json:"deleted_by"`
}

type Id struct {
	ID string `json:"id" param:"id"`
}

type Tabler interface {
	TableName() string
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

type PagedData struct {
	List interface{} `json:"lists"`
	Meta MetaData    `json:"metadata"`
}

type MetaData struct {
	Total int `json:"total"`
}
