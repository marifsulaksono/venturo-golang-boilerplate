package helpers

import (
	"fmt"
	"net/http"
	"runtime"
	"simple-crud-rnd/structs"

	"github.com/getsentry/sentry-go"
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

	fmt.Println(status, response)
	return c.JSON(status, response)
}

func ResponseError(c echo.Context, status int, err error, message string) error {
	response := structs.JSONResponse{
		ResponseCode:    status,
		ResponseMessage: getMessage(status),
		Message:         message,
	}

	buf := make([]byte, 1<<16)            // Allocate a buffer to hold the stack trace
	stackSize := runtime.Stack(buf, true) // Capture the stack trace

	// Kirim ke Sentry
	sentry.CaptureException(fmt.Errorf("panic: %v\nStack trace:\n%s", err, buf[:stackSize]))
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
