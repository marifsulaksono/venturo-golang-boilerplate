package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

type (
	ImageHelper struct {
		fullPath    string
		storagePath string
		category    string
	}
)

func NewImageHelper(storagePath, category string) (*ImageHelper, error) {
	fullPath := fmt.Sprintf("%s/%s/%s", storagePath, "images", category)
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return &ImageHelper{}, nil
	}
	return &ImageHelper{fullPath, storagePath, category}, nil
}

func (img *ImageHelper) Writer(imageString string, filename string) (string, error) {
	// Check if imageString contains the MIME type prefix
	if strings.Contains(imageString, "base64,") {
		// Remove the MIME type prefix
		parts := strings.Split(imageString, "base64,")
		if len(parts) > 1 {
			imageString = parts[1] // Get the actual Base64 string
		}
	}

	// Decode the Base64 string
	dec, err := base64.StdEncoding.DecodeString(imageString)
	if err != nil {
		return "", err
	}

	// Create the file
	f, err := os.Create(fmt.Sprintf("%s/%s", img.fullPath, filename))
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Write the decoded data to the file
	if _, err := f.Write(dec); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}

	// Return the file path
	return fmt.Sprintf("%s/%s/%s", "images", img.category, filename), nil
}

func (img *ImageHelper) Read(filepath string) (string, error) {
	// Open the file
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the file contents
	fileInfo, err := f.Stat()
	if err != nil {
		return "", err
	}

	fileSize := fileInfo.Size()
	fileBytes := make([]byte, fileSize)

	_, err = f.Read(fileBytes)
	if err != nil {
		return "", err
	}

	// Convert the file content to Base64 string
	encodedString := base64.StdEncoding.EncodeToString(fileBytes)

	// return the MIME type prefix for image format
	return fmt.Sprintf("data:image/png;base64,%s", encodedString), nil
}
