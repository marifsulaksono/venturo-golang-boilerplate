package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		return "", SendTraceErrorToSentry(err)
	}

	// Create the file
	f, err := os.Create(fmt.Sprintf("%s/%s", img.fullPath, filename))
	if err != nil {
		return "", SendTraceErrorToSentry(err)
	}
	defer f.Close()

	// Write the decoded data to the file
	if _, err := f.Write(dec); err != nil {
		return "", SendTraceErrorToSentry(err)
	}
	if err := f.Sync(); err != nil {
		return "", SendTraceErrorToSentry(err)
	}

	// Return the file path
	return fmt.Sprintf("%s/%s/%s", "images", img.category, filename), nil
}

func (img *ImageHelper) Read(filepath string) (string, error) {
	// Open the file
	f, err := os.Open(filepath)
	if err != nil {
		return "", SendTraceErrorToSentry(err)
	}
	defer f.Close()

	// Read the file contents
	fileInfo, err := f.Stat()
	if err != nil {
		return "", SendTraceErrorToSentry(err)
	}

	fileSize := fileInfo.Size()
	fileBytes := make([]byte, fileSize)

	_, err = f.Read(fileBytes)
	if err != nil {
		return "", SendTraceErrorToSentry(err)
	}

	// Convert the file content to Base64 string
	encodedString := base64.StdEncoding.EncodeToString(fileBytes)

	// return the MIME type prefix for image format
	return fmt.Sprintf("data:image/png;base64,%s", encodedString), nil
}

func MoveToTrash(imagePath string) error {
	fmt.Println("imagePath:", imagePath)
	// Tentukan path folder trash
	trashFolder := "./assets/trash/"
	// Dapatkan nama file
	fileName := filepath.Base(imagePath)
	// Tentukan path baru untuk file di folder trash
	destination := filepath.Join(trashFolder, fileName)

	// Cek apakah folder trash ada, jika belum buat
	if _, err := os.Stat(trashFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(trashFolder, os.ModePerm); err != nil {
			return err
		}
	}

	// Pindahkan file ke folder trash
	if err := os.Rename(imagePath, destination); err != nil {
		return err
	}
	return nil
}

// CleanUpTrash menghapus file yang sudah ada di folder trash selama lebih dari 1 bulan
func CleanUpTrash() error {
	trashFolder := "./assets/trash/"

	// Check if the trash folder exists
	if _, err := os.Stat(trashFolder); os.IsNotExist(err) {
		fmt.Println("Trash folder does not exist.")
		return nil
	}

	// Read the files in the trash folder
	files, err := os.ReadDir(trashFolder)
	if err != nil {
		return err
	}

	// Set the threshold to 1 minute ago for testing (use -1 month for production)
	thresholdTime := time.Now().Add(-1 * time.Minute)

	// Loop through each file and check modification time
	for _, file := range files {
		filePath := filepath.Join(trashFolder, file.Name())
		info, err := file.Info()
		if err != nil {
			return err
		}

		// Print file info for debugging
		fmt.Printf("Checking file: %s, ModTime: %s\n", file.Name(), info.ModTime())

		// Check if the file modification time is before the threshold
		if info.ModTime().Before(thresholdTime) {
			// Attempt to delete the file
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Failed to delete file %s: %v\n", filePath, err)
			} else {
				fmt.Printf("Deleted file: %s\n", file.Name())
			}
		} else {
			fmt.Printf("File %s is not old enough to be deleted.\n", file.Name())
		}
	}

	return nil
}
