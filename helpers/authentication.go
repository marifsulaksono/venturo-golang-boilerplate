package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/url"

	"github.com/dgryski/dgoogauth"
)

func Generate2FAQRCodeURL(username, key string) string {
	accountName := "VenturoApp:" + username
	issuer := "Venturo"
	authURL := url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   fmt.Sprintf("/%s", accountName),
		RawQuery: url.Values{
			"secret": {key},
			"issuer": {issuer},
		}.Encode(),
	}
	return authURL.String()
}

func VerifyOTP(secret, code string) (bool, error) {
	otpConfig := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3, // jumlah percobaan kode yang valid dalam 90 detik (disesuaikan)
		HotpCounter: 0,
	}
	verified, err := otpConfig.Authenticate(code)
	return verified, err
}

func Generate2FASecretKey(ctx context.Context) (string, error) {
	randomBytes := make([]byte, 10)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes), nil
}
