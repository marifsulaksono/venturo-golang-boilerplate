package middleware

import (
	"encoding/json"
	"net/http"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"
	"strings"

	"github.com/labstack/echo/v4"
)

func RoleMiddleware(requiredRoles string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return helpers.Response(c, http.StatusUnauthorized, nil, "Invalid Token")
			}
			user, err := helpers.VerifyTokenJWT(tokenString, false)
			if err != nil {
				return helpers.Response(c, http.StatusUnauthorized, err.Error(), "Invalid Token")
			}

			c.Set("user_id", user.ID) // set saves data in the context

			// check if user permissions
			if requiredRoles == "" {
				// skip validation if no require
				return next(c)
			} else if !hasRole(user, requiredRoles) {
				return helpers.Response(c, http.StatusUnauthorized, nil, "Access denied")
			}
			return next(c)
		}
	}
}

// check if the user has the required permissions
func hasRole(user *structs.User, requiredRoles string) bool {
	var accessMap map[string]map[string]bool
	err := json.Unmarshal([]byte(user.Role.Access), &accessMap)
	if err != nil {
		return false
	}

	requiredPermissions := strings.Split(requiredRoles, "|")
	for _, val := range requiredPermissions {
		permissionParts := strings.Split(val, ".")
		if len(permissionParts) < 2 {
			continue
		}

		feature, activity := permissionParts[0], permissionParts[1]

		// check if the user has access to the specified feature and activity
		if !hasAccess(accessMap, feature, activity) {
			return false
		}
	}
	return true // no required permissions or permission were valid
}

func hasAccess(accessMap map[string]map[string]bool, feature, activity string) bool {
	if permissions, exists := accessMap[feature]; exists {
		return permissions[activity] // check if the specific activity is true
	}
	return false
}
