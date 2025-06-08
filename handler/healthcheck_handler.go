package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io" // Thay thế ioutil bằng io
	"log"
	"net/http"
	"time"
)

// Placeholder types for demonstration.
type CommitInfo struct {
	ShaCode string    `json:"sha_code"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

// Healthcheck verifies the validity of an API key.
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	accessKey := "some-access-key" // Replace with actual logic to retrieve the access key
	if accessKey == "" {
		http.Error(w, "Invalid Access Key", http.StatusUnauthorized)
		return
	}

	isValidApiKey, err := validateApiKey(accessKey)
	if err != nil || !isValidApiKey {
		http.Error(w, "Invalid Access Key", http.StatusUnauthorized)
		return
	}

	log.Printf("Gemini API Key: %s", accessKey)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(isValidApiKey)
}

// SendFeedback logs user feedback.
func SendFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback struct {
		UserName     string `json:"user_name"`
		UserFeedback string `json:"user_feedback"`
	}

	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	log.Printf("%s's feedback: %s", feedback.UserName, feedback.UserFeedback)
	w.WriteHeader(http.StatusNoContent)
}

// ExtractTextFromImage decodes a base64 image and extracts text.
func ExtractTextFromImage(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Base64Image string `json:"base64_image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if request.Base64Image == "" {
		http.Error(w, "Base64 image is required", http.StatusBadRequest)
		return
	}

	accessKey := "some-access-key" // Replace with actual logic to retrieve the access key
	if accessKey == "" {
		http.Error(w, "Invalid Access Key", http.StatusUnauthorized)
		return
	}

	content, err := extractTextFromBase64Image(accessKey, request.Base64Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(content)
}

// GetLatestGithubCommit fetches the latest commit from the GitHub repository.
func GetLatestGithubCommit(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/phanxuanquang/EngPal/commits/master", nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "request")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch commit", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Thay thế ioutil.ReadAll bằng io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var commitData map[string]interface{}
	if err := json.Unmarshal(body, &commitData); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	commitInfo := CommitInfo{
		ShaCode: commitData["sha"].(string),
		Message: commitData["commit"].(map[string]interface{})["message"].(string),
		Date:    time.Now(), // Replace with actual parsing logic for date
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commitInfo)
}

// Helper function to validate the API key (placeholder).
func validateApiKey(apiKey string) (bool, error) {
	// Simulate API key validation.
	return true, nil
}

// Helper function to extract text from a base64 image (placeholder).
func extractTextFromBase64Image(apiKey, base64Image string) (string, error) {
	decodedImage, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", errors.New("failed to decode base64 image")
	}

	// Simulate text extraction from the image.
	// You can replace this with actual OCR logic.
	if len(decodedImage) == 0 {
		return "", errors.New("image content is empty")
	}

	return "Extracted text from image", nil
}
