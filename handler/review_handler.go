package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Placeholder types for demonstration.
type GenerateCommentRequest struct {
	Content     string `json:"content"`
	UserLevel   string `json:"user_level"`
	Requirement string `json:"requirement"`
}

// GenerateReview handles the generation of a review based on user input.
func GenerateReview(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request.
	var request GenerateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	accessKey := "some-access-key" // Placeholder for actual access key retrieval logic.
	if accessKey == "" {
		http.Error(w, "Invalid Access Key", http.StatusUnauthorized)
		return
	}

	content := strings.TrimSpace(request.Content)

	// Validate content length.
	minTotalWords := 10  // Placeholder for ReviewScope.MinTotalWords
	maxTotalWords := 500 // Placeholder for ReviewScope.MaxTotalWords

	wordCount := getTotalWords(content)
	if wordCount < minTotalWords {
		http.Error(w, fmt.Sprintf("Bài viết phải dài tối thiểu %d từ.", minTotalWords), http.StatusBadRequest)
		return
	}

	if wordCount > maxTotalWords {
		http.Error(w, fmt.Sprintf("Bài viết không được dài hơn %d từ.", maxTotalWords), http.StatusBadRequest)
		return
	}

	// Generate review.
	result, err := generateReview(accessKey, request.UserLevel, request.Requirement, content)
	if err != nil {
		log.Printf("Error generating review: %v", err)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("## CẢNH BÁO\n EngPal đang bận đi pha cà phê nên tạm thời vắng mặt. bé yêu vui lòng ngồi chơi 3 phút rồi gửi lại cho EngPal nhận xét nha.\nYêu bé yêu nhiều lắm luôn á!")
		return
	}

	// Return the result.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Helper function to count total words in a string.
func getTotalWords(input string) int {
	return len(strings.Fields(input))
}

// Helper function to simulate review generation (placeholder).
func generateReview(accessKey, userLevel, requirement, content string) (string, error) {
	// Simulate review generation logic.
	if strings.Contains(content, "error") {
		return "", errors.New("failed to generate review")
	}
	return "Đây là nhận xét mẫu từ hệ thống EngPal!", nil
}
