package helpers

import (
	"os"
	"simple-crud-rnd/structs"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(user *structs.User, isRefresh bool) (string, *time.Time, error) {
	var (
		expiredInSecond int
		secretKey       string
	)

	if isRefresh {
		expiredInSecond, _ = strconv.Atoi(os.Getenv("REFRESH_JWT_EXPIRY_IN_SECOND"))
		secretKey = os.Getenv("ACCESS_JWT_SECRET")
	} else {
		expiredInSecond, _ = strconv.Atoi(os.Getenv("ACCESS_JWT_EXPIRY_IN_SECOND"))
		secretKey = os.Getenv("REFRESH_JWT_SECRET")
	}

	expiredAt := time.Now().Add(time.Second * time.Duration(expiredInSecond))
	claims := &structs.JWTUser{
		ID:              user.ID,
		Email:           user.Name,
		UpdatedSecurity: user.UpdatedSecurity,
		Access:          user.Role.Access,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "venturo",
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &expiredAt, nil
}

func VerifyTokenJWT(tokenString string, isRefresh bool) (*structs.User, error) {
	var (
		secretKey string
		user      = new(structs.User)
	)

	if isRefresh {
		secretKey = os.Getenv("REFRESH_JWT_EXPIRY_IN_SECOND")
	} else {
		secretKey = os.Getenv("ACCESS_JWT_EXPIRY_IN_SECOND")
	}

	token, err := jwt.ParseWithClaims(tokenString, &structs.JWTUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
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
