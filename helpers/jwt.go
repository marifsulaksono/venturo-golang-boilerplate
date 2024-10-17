package helpers

import (
	"os"
	"simple-crud-rnd/structs"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(user structs.User) (string, *time.Time, error) {
	expiredInSecond, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_IN_SECOND"))
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
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &expiredAt, nil
}
