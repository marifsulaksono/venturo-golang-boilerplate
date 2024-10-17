package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(requiredRoles string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			secretKey := os.Getenv("JWT_SECRET")
			user, err := parseJWT(c, []byte(secretKey))
			if err != nil {
				return helpers.Response(c, http.StatusUnauthorized, err.Error(), "Invalid Token")
			}

			c.Set("user", user) // set saves data in the context

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

// parse jwt token
func parseJWT(c echo.Context, secretKey []byte) (*structs.User, error) {
	user := new(structs.User)
	authHeader := c.Request().Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	token, err := jwt.ParseWithClaims(tokenString, &structs.JWTUser{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// extract user claims if the token is valid
	if claims, ok := token.Claims.(*structs.JWTUser); ok && token.Valid {
		// set user properties
		if user.Role == nil {
			user.Role = &structs.Role{}
		}

		user.ID = claims.ID
		user.Name = claims.Email
		user.Role.Access = claims.Access

		return user, nil
	}

	return nil, err
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
