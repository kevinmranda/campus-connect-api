package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SmsRecipient struct {
	Number int64 `json:"number"`
}

type SmsRequest struct {
	SenderID   int            `json:"sender_id"`
	SMS        string         `json:"sms"`
	Schedule   string         `json:"schedule"`
	Recipients []SmsRecipient `json:"recipients"`
}

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

func SendRegistrationSMS(message string, recipients []string) error {
	// Step 1: Authenticate to get the token
	authBody := map[string]interface{}{
		"user_id":  255787504956,
		"password": "adventist145",
	}
	authJson, _ := json.Marshal(authBody)

	authResp, err := http.Post("https://dev.hudumasms.com/api/create-token", "application/json", bytes.NewBuffer(authJson))
	if err != nil {
		return fmt.Errorf("failed to authenticate to SMS API: %w", err)
	}
	defer authResp.Body.Close()

	var authResult struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	bodyBytes, _ := io.ReadAll(authResp.Body)
	if err := json.Unmarshal(bodyBytes, &authResult); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}
	if authResult.Data.Token == "" {
		return errors.New("token not received from SMS API")
	}

	// Step 2: Prepare recipients
	var formattedRecipients []SmsRecipient
	for _, rec := range recipients {
		var number int64
		fmt.Sscanf(rec, "%d", &number)
		formattedRecipients = append(formattedRecipients, SmsRecipient{Number: number})
	}

	// Step 3: Send the SMS
	smsBody := SmsRequest{
		SenderID:   20,
		SMS:        message,
		Schedule:   "None",
		Recipients: formattedRecipients,
	}
	smsJson, _ := json.Marshal(smsBody)

	req, _ := http.NewRequest("POST", "https://dev.hudumasms.com/api/send-sms", bytes.NewBuffer(smsJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResult.Data.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("SMS API returned non-200 status: %d, %s", resp.StatusCode, string(respBody))
	}

	return nil
}
