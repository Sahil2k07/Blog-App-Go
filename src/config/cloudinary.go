package config

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	cloudinary "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func Cloudinary(file multipart.File) (string, error) {
	CLOUD_NAME := os.Getenv("CLOUD_NAME")
	API_KEY := os.Getenv("API_KEY")
	API_SECRET := os.Getenv("API_SECRET")
	FOLDER_NAME := os.Getenv("FOLDER_NAME")

	cld, err := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	if err != nil {
		return "", fmt.Errorf("failed to create Cloudinary client: %v", err)
	}

	// Upload the image
	response, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: FOLDER_NAME,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	return response.SecureURL, nil
}
