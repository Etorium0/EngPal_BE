package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"EngPal/internal"

	"google.golang.org/genai"
)

// Request/Response types
type GenerateCommentRequest struct {
	Content     string `json:"content"`
	UserLevel   string `json:"user_level"`
	Requirement string `json:"requirement"`
	Category    string `json:"category,omitempty"` // writing, speaking, etc.
	Language    string `json:"language,omitempty"` // en, vi for response language
}

type ReviewCriteria struct {
	Grammar      float64 `json:"grammar"`       // 0-10
	Vocabulary   float64 `json:"vocabulary"`    // 0-10
	Coherence    float64 `json:"coherence"`     // 0-10
	TaskResponse float64 `json:"task_response"` // 0-10
	Overall      float64 `json:"overall"`       // 0-10
}

type ReviewSuggestion struct {
	Category   string `json:"category"`   // Grammar, Vocabulary, etc.
	Issue      string `json:"issue"`      // What's wrong
	Suggestion string `json:"suggestion"` // How to fix
	Example    string `json:"example"`    // Better version
	Priority   string `json:"priority"`   // High, Medium, Low
}

type ReviewResponse struct {
	Content          string             `json:"content"`
	UserLevel        string             `json:"user_level"`
	Requirement      string             `json:"requirement"`
	WordCount        int                `json:"word_count"`
	EstimatedLevel   string             `json:"estimated_level"`
	Scores           ReviewCriteria     `json:"scores"`
	OverallFeedback  string             `json:"overall_feedback"`
	StrengthPoints   []string           `json:"strength_points"`
	ImprovementAreas []string           `json:"improvement_areas"`
	Suggestions      []ReviewSuggestion `json:"suggestions"`
	CorrectedVersion string             `json:"corrected_version,omitempty"`
	GeneratedAt      time.Time          `json:"generated_at"`
	ProcessingTime   float64            `json:"processing_time_ms"`
}

// Gemini API structures for review
type GeminiReviewRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiReviewResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiReviewData struct {
	EstimatedLevel   string             `json:"estimated_level"`
	Scores           ReviewCriteria     `json:"scores"`
	OverallFeedback  string             `json:"overall_feedback"`
	StrengthPoints   []string           `json:"strength_points"`
	ImprovementAreas []string           `json:"improvement_areas"`
	Suggestions      []ReviewSuggestion `json:"suggestions"`
	CorrectedVersion string             `json:"corrected_version,omitempty"`
}

// Cache for reviews
type reviewCacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

var reviewCache = make(map[string]reviewCacheItem)

// Constants
const (
	MIN_TOTAL_WORDS = 10
	MAX_TOTAL_WORDS = 1000
	CACHE_DURATION  = 1 * time.Hour // Cache for 1 hour like C# version
)

// English level mapping
var reviewEnglishLevels = map[string]string{
	"A1": "A1 - Beginner",
	"A2": "A2 - Elementary",
	"B1": "B1 - Intermediate",
	"B2": "B2 - Upper Intermediate",
	"C1": "C1 - Advanced",
	"C2": "C2 - Proficient",
}

// Writing categories
var writingCategories = map[string]string{
	"essay":       "Academic Essay",
	"letter":      "Formal/Informal Letter",
	"report":      "Report Writing",
	"article":     "Article Writing",
	"story":       "Creative Writing",
	"email":       "Email Writing",
	"description": "Descriptive Writing",
	"opinion":     "Opinion Writing",
}

// --- MAIN HANDLER ---

func GenerateReview(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(CACHE_DURATION.Seconds())))

	var request GenerateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validation
	if err := validateReviewRequest(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check cache
	cacheKey := generateReviewCacheKey(request)
	now := time.Now()
	if item, found := reviewCache[cacheKey]; found && item.ExpiresAt.After(now) {
		log.Printf("Serving cached review for content hash: %s", cacheKey[:10])
		json.NewEncoder(w).Encode(item.Data)
		return
	}

	// Generate review using Gemini API
	reviewResponse, err := generateReviewWithGemini(request, startTime)
	if err != nil {
		log.Printf("Error generating review: %v", err)
		// Return friendly error message like C# version
		errorResponse := map[string]string{
			"error":   "service_unavailable",
			"message": "## CẢNH BÁO\nEngPal đang bận đi pha cà phê nên tạm thời vắng mặt. bé yêu vui lòng ngồi chơi 3 phút rồi gửi lại cho EngPal nhận xét nha.\nYêu bé yêu nhiều lắm luôn á!",
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Cache the response
	reviewCache[cacheKey] = reviewCacheItem{
		Data:      reviewResponse,
		ExpiresAt: now.Add(CACHE_DURATION),
	}

	log.Printf("Generated review for %d words, processing time: %.2fms",
		reviewResponse.WordCount, reviewResponse.ProcessingTime)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviewResponse)
}

// Validate review request
func validateReviewRequest(request GenerateCommentRequest) error {
	request.Content = strings.TrimSpace(request.Content)
	if request.Content == "" {
		return errors.New("nội dung bài viết không được để trống")
	}

	wordCount := getTotalWords(request.Content)
	if wordCount < MIN_TOTAL_WORDS {
		return fmt.Errorf("bài viết phải dài tối thiểu %d từ", MIN_TOTAL_WORDS)
	}

	if wordCount > MAX_TOTAL_WORDS {
		return fmt.Errorf("bài viết không được dài hơn %d từ", MAX_TOTAL_WORDS)
	}

	if request.UserLevel != "" {
		if _, exists := reviewEnglishLevels[strings.ToUpper(request.UserLevel)]; !exists {
			return errors.New("trình độ tiếng Anh không hợp lệ (A1, A2, B1, B2, C1, C2)")
		}
	}

	return nil
}

// Generate review using Gemini API
func generateReviewWithGemini(req GenerateCommentRequest, startTime time.Time) (*ReviewResponse, error) {
	// Build comprehensive prompt
	prompt := buildReviewPrompt(req)

	// Call Gemini API
	geminiResp, err := callGeminiForReview(prompt)
	if err != nil {
		return nil, fmt.Errorf("gemini API call failed: %w", err)
	}

	// Parse response
	reviewData, err := parseGeminiReviewResponse(geminiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}

	// Build final response
	processingTime := float64(time.Since(startTime).Nanoseconds()) / 1e6 // Convert to milliseconds

	response := &ReviewResponse{
		Content:          req.Content,
		UserLevel:        req.UserLevel,
		Requirement:      req.Requirement,
		WordCount:        getTotalWords(req.Content),
		EstimatedLevel:   reviewData.EstimatedLevel,
		Scores:           reviewData.Scores,
		OverallFeedback:  reviewData.OverallFeedback,
		StrengthPoints:   reviewData.StrengthPoints,
		ImprovementAreas: reviewData.ImprovementAreas,
		Suggestions:      reviewData.Suggestions,
		CorrectedVersion: reviewData.CorrectedVersion,
		GeneratedAt:      time.Now(),
		ProcessingTime:   processingTime,
	}

	return response, nil
}

// Build comprehensive review prompt for Gemini
func buildReviewPrompt(req GenerateCommentRequest) string {
	userLevelDesc := "intermediate"
	if req.UserLevel != "" {
		if level, exists := reviewEnglishLevels[strings.ToUpper(req.UserLevel)]; exists {
			userLevelDesc = level
		}
	}

	category := "general writing"
	if req.Category != "" {
		if cat, exists := writingCategories[strings.ToLower(req.Category)]; exists {
			category = cat
		}
	}

	responseLanguagePrompt := "English"
	if req.Language == "vi" {
		responseLanguagePrompt = "Tiếng Việt"
	}

	wordCount := getTotalWords(req.Content)

	prompt := fmt.Sprintf(`You are an expert English teacher and IELTS examiner. Analyze the following English writing sample and provide a comprehensive review.

WRITING SAMPLE TO ANALYZE:
"%s"

CONTEXT INFORMATION:
- Student's declared level: %s
- Writing category: %s
- Specific requirement: %s
- Word count: %d

ANALYSIS REQUIREMENTS:
1. Estimate the actual English level (A1-C2) based on the writing quality
2. Score each criterion from 0-10:
   - Grammar: Accuracy, complexity, range of structures
   - Vocabulary: Range, accuracy, appropriateness
   - Coherence: Logical flow, linking, organization
   - Task Response: Meeting requirements, completeness
   - Overall: Holistic impression

3. Provide specific feedback covering:
   - 3-5 strength points (what the student does well)
   - 3-5 improvement areas (what needs work)
   - 5-8 detailed suggestions with examples
   - overall_feedback: Tổng nhận xét chung về bài viết (bắt buộc)

4. If there are significant errors, provide a corrected version

FORMATTING REQUIREMENTS:
Return ONLY valid JSON without markdown formatting.
JSON phải có các trường sau (bắt buộc):
- "estimated_level"
- "scores" (bao gồm: "grammar", "vocabulary", "coherence", "task_response", "overall", tất cả đều là số từ 0 đến 10)
- "overall_feedback"
- "strength_points"
- "improvement_areas"
- "suggestions" (mảng các object, mỗi object gồm: "category", "issue", "suggestion", "example", "priority")
- "corrected_version" (nếu có)

Ví dụ trường "suggestions":
"suggestions": [
  {
    "category": "Grammar",
    "issue": "Subject-verb agreement",
    "suggestion": "Kiểm tra sự hòa hợp giữa chủ ngữ và động từ.",
    "example": "Incorrect: 'She go to school.' Correct: 'She goes to school.'",
    "priority": "High"
  }
]

Nếu không có thông tin cho trường nào, vẫn phải trả về trường đó với giá trị hợp lệ (ví dụ: 0 cho điểm số, chuỗi rỗng cho text).

IMPORTANT: Tất cả phản hồi (bao gồm nhận xét, điểm số, gợi ý, bản sửa lỗi) PHẢI được viết hoàn toàn bằng %s.

Analyze the writing sample now:`, req.Content, userLevelDesc, category, req.Requirement, wordCount, responseLanguagePrompt)

	return prompt
}

// Call Gemini API for review
func callGeminiForReview(prompt string) (string, error) {
	client := internal.GeminiClient
	if client == nil {
		return "", errors.New("Gemini client not initialized")
	}

	ctx := context.Background()
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp", // Use experimental model for better analysis
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

// Parse Gemini response for review
func parseGeminiReviewResponse(response string) (*GeminiReviewData, error) {
	// Clean the response
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var reviewData GeminiReviewData
	err := json.Unmarshal([]byte(response), &reviewData)
	if err != nil {
		// Try fallback: parse suggestions as []string
		var fallback struct {
			EstimatedLevel   string         `json:"estimated_level"`
			Scores           ReviewCriteria `json:"scores"`
			OverallFeedback  string         `json:"overall_feedback"`
			StrengthPoints   []string       `json:"strength_points"`
			ImprovementAreas []string       `json:"improvement_areas"`
			Suggestions      []string       `json:"suggestions"`
			CorrectedVersion string         `json:"corrected_version,omitempty"`
		}
		if err2 := json.Unmarshal([]byte(response), &fallback); err2 == nil {
			// Convert []string to []ReviewSuggestion
			sugs := make([]ReviewSuggestion, len(fallback.Suggestions))
			for i, s := range fallback.Suggestions {
				sugs[i] = ReviewSuggestion{
					Category:   "",
					Issue:      "",
					Suggestion: s,
					Example:    "",
					Priority:   "",
				}
			}
			return &GeminiReviewData{
				EstimatedLevel:   fallback.EstimatedLevel,
				Scores:           fallback.Scores,
				OverallFeedback:  fallback.OverallFeedback,
				StrengthPoints:   fallback.StrengthPoints,
				ImprovementAreas: fallback.ImprovementAreas,
				Suggestions:      sugs,
				CorrectedVersion: fallback.CorrectedVersion,
			}, nil
		}
		log.Printf("Failed to parse review JSON response: %s", response)
		return nil, fmt.Errorf("failed to parse review JSON: %w", err)
	}

	// Validate required fields
	if reviewData.EstimatedLevel == "" {
		reviewData.EstimatedLevel = "B1" // Default
	}
	if reviewData.OverallFeedback == "" {
		return nil, errors.New("missing overall feedback in API response")
	}

	// Ensure we have some suggestions
	if len(reviewData.Suggestions) == 0 {
		reviewData.Suggestions = []ReviewSuggestion{
			{
				Category:   "General",
				Issue:      "Continue practicing",
				Suggestion: "Keep writing regularly to improve your skills",
				Example:    "Practice different types of writing",
				Priority:   "Medium",
			},
		}
	}

	return &reviewData, nil
}

// Helper function to count words
func getTotalWords(input string) int {
	if strings.TrimSpace(input) == "" {
		return 0
	}
	return len(strings.Fields(input))
}

// Generate cache key for reviews
func generateReviewCacheKey(req GenerateCommentRequest) string {
	// Create a hash-like key based on content and parameters
	key := strings.ToLower(req.Content) + "-" + req.UserLevel + "-" + req.Requirement + "-" + req.Category
	// In production, you might want to use actual hashing
	return fmt.Sprintf("%x", len(key)) + "-" + strconv.Itoa(getTotalWords(req.Content))
}

// --- ADDITIONAL ENDPOINTS ---

// Get available English levels
func GetReviewLevels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviewEnglishLevels)
}

// Get writing categories
func GetWritingCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(writingCategories)
}

// Get review statistics (for admin/monitoring)
func GetReviewStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"cache_entries":    len(reviewCache),
		"min_words":        MIN_TOTAL_WORDS,
		"max_words":        MAX_TOTAL_WORDS,
		"cache_duration":   CACHE_DURATION.String(),
		"available_levels": len(reviewEnglishLevels),
		"categories":       len(writingCategories),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Clear review cache (for admin)
func ClearReviewCache(w http.ResponseWriter, r *http.Request) {
	reviewCache = make(map[string]reviewCacheItem)

	response := map[string]string{
		"status":  "success",
		"message": "Review cache cleared successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
