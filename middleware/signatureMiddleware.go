package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"simple-crud-rnd/helpers"
	"strings"

	"github.com/labstack/echo/v4"
)

func SignatureMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut || c.Request().Method == http.MethodDelete {
				signature := c.Request().Header.Get("signature")

				if signature == "" {
					return helpers.Response(c, http.StatusUnauthorized, nil, "Signature is not provided. Access Denied.")
				}

				if signature == os.Getenv("SIGNATURE_BYPASS") {
					return next(c)
				}

				payload := getPayload(c)
				if payload == nil {
					return helpers.Response(c, http.StatusBadRequest, nil, "Failed to retrieve payload.")
				}

				decryptedSignatureString, err := helpers.DecryptMessageRSA(signature)
				if err != nil {
					return helpers.Response(c, http.StatusUnauthorized, nil, "Signature mismatch. Operation cannot be completed.")
				}

				if c.Request().Method == http.MethodDelete {
					if payload == decryptedSignatureString {
						return next(c)
					} else {
						return helpers.Response(c, http.StatusUnauthorized, nil, "Signature mismatch. Operation cannot be completed.")
					}
				}

				var decryptedSignatureMap map[string]interface{}
				err = json.Unmarshal([]byte(decryptedSignatureString), &decryptedSignatureMap)
				if err != nil {
					return helpers.Response(c, http.StatusInternalServerError, nil, "Error unmarshalling decrypted signature.")
				}

				payloadMap, ok := payload.(map[string]interface{})
				if !ok {
					return helpers.Response(c, http.StatusBadRequest, nil, "Invalid payload format.")
				}

				if !compareMaps(decryptedSignatureMap, payloadMap) {
					return helpers.Response(c, http.StatusUnauthorized, nil, "Signature mismatch. Operation cannot be completed.")
				}
			}

			return next(c)
		}
	}
}

func getPayload(c echo.Context) interface{} {
	if c.Request().Method == http.MethodDelete {
		pathSegments := strings.Split(c.Request().URL.Path, "/")
		return pathSegments[len(pathSegments)-1]
	}

	// ini sama seperti c.Bind
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, "Failed to read request body.")
	}

	// isi ulang body agar bisa dibaca lagi di controller
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, "Failed to parse payload.")
	}
	return payload
}

func compareMaps(map1, map2 map[string]interface{}) bool {
	for key, val1 := range map1 {
		val2, exists := map2[key]
		if !exists {
			return false
		}

		switch v1 := val1.(type) {
		case float64:
			v2, ok := val2.(float64)
			if !ok || v1 != v2 {
				return false
			}
		case string:
			v2, ok := val2.(string)
			if !ok || v1 != v2 {
				return false
			}
		default:
			if !reflect.DeepEqual(val1, val2) {
				return false
			}
		}
	}

	return true
}
