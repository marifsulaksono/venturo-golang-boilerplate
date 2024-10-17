package helpers

import (
	"net/http"
	"simple-crud-rnd/structs"

	"github.com/labstack/echo/v4"
)

func getMessage(status int) string {
	switch status {
	case http.StatusOK:
		return "Success"
	case http.StatusCreated:
		return "Created"
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}

func Response(c echo.Context, status int, data interface{}, message string) error {
	response := structs.JSONResponse{
		ResponseCode:    status,
		ResponseMessage: getMessage(status),
		Message:         message,
		Data:            data,
	}

	return c.JSON(status, response)
}

func PageData(data interface{}, total int64) *structs.PagedData {
	return &structs.PagedData{
		List: data,
		Meta: structs.MetaData{
			Total: int(total),
		},
	}
}
