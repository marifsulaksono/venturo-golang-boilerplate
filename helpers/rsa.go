package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func EncryptMessageRSA(plaintext string) (string, error) {
	publicKey, err := base64.StdEncoding.DecodeString(os.Getenv("RSA_PUBLIC_KEY"))
	if err != nil {
		return "", fmt.Errorf("failed to load public key: %s", err)
	}

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("not an RSA public key")
	}

	// encrypt the plaintext using the RSA public key
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, []byte(plaintext), nil)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptMessageRSA(ciphertext string) (string, error) {
	privateKey, err := base64.StdEncoding.DecodeString(os.Getenv("RSA_PRIVATE_KEY"))
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %s", err)
	}

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %v", err)
	}

	// decrypt the ciphertext using the RSA private key
	plaintextBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %v", err)
	}

	return string(plaintextBytes), nil
}
