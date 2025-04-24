package helpers

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"regexp"
)

// SaveImage saves a base64-encoded image to the Images folder and returns the relative file path
func SaveImage(base64Image, filename string) (string, error) {
	imageDir := "./Images"
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create Images directory: %v", err)
	}

	data := base64Image
	if strings.Contains(base64Image, ",") {
		data = strings.SplitN(base64Image, ",", 2)[1]
	}

	imgData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}

	ext := ".jpg"
	if strings.Contains(base64Image, "image/png") {
		ext = ".png"
	} else if strings.Contains(base64Image, "image/jpeg") {
		ext = ".jpeg"
	}

	filePath := filepath.Join(imageDir, filename+ext)

	if err := os.WriteFile(filePath, imgData, 0644); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	return filePath, nil
}

// UpdateImageName renames an existing image file to a new filename
func UpdateImageName(oldPath, newFilename string) (string, error) {

	ext := filepath.Ext(oldPath)

	newPath := filepath.Join("./Images", newFilename+ext)

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", fmt.Errorf("failed to rename image: %v", err)
	}

	return newPath, nil
}

// DeleteImage deletes an image file from the Images folder
func DeleteImage(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}
	return nil
}

// SanitizeFilename converts a string into a safe filename
func SanitizeFilename(title string) string {

	result := strings.ToLower(title)
	
	reg, _ := regexp.Compile("[^a-z0-9]+")
	result = reg.ReplaceAllString(result, "-")
	
	result = strings.Trim(result, "-")

	if len(result) > 50 {
		result = result[:50]
	}

	if result == "" {
		result = "untitled"
	}
	return result
}

// SanitizeUserIDFilename creates a safe filename for profile images using the user ID
func SanitizeUserIDFilename(prefix string, userID uint) string {
	
	filename := fmt.Sprintf("%s%d", prefix, userID)
	
	reg, _ := regexp.Compile("[^a-zA-Z0-9-]+")
	filename = reg.ReplaceAllString(filename, "")
	
	if len(filename) > 50 {
		filename = filename[:50]
	}
	
	if filename == "" {
		filename = "profile-picture-UID0"
	}
	return filename
}