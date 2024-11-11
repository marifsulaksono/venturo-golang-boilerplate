package helpers

import (
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
