package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"EngPal/internal"

	"google.golang.org/genai"
)

// Request/Response types
type GenerateQuizzesRequest struct {
	Topic           string   `json:"topic"`
	AssignmentTypes []string `json:"assignment_types"`
	EnglishLevel    string   `json:"english_level"`
	TotalQuestions  int      `json:"total_questions"`
}

type Quiz struct {
	ID           int      `json:"id"`
	Type         string   `json:"type"`
	Question     string   `json:"question"`
	Answer       string   `json:"answer,omitempty"`
	Options      []string `json:"options,omitempty"`
	CorrectIndex int      `json:"correct_index,omitempty"`
	Explanation  string   `json:"explanation,omitempty"`
}

type QuizResponse struct {
	Topic     string `json:"topic"`
	Level     string `json:"level"`
	Total     int    `json:"total"`
	Generated int    `json:"generated"`
	Quizzes   []Quiz `json:"quizzes"`
}

// Gemini API structures
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

type GeminiQuizData struct {
	Quizzes []GeminiQuiz `json:"quizzes"`
}

type GeminiQuiz struct {
	Type         string   `json:"type"`
	Question     string   `json:"question"`
	Answer       string   `json:"answer,omitempty"`
	Options      []string `json:"options,omitempty"`
	CorrectIndex int      `json:"correct_index,omitempty"`
	Explanation  string   `json:"explanation,omitempty"`
}

// Cache struct
type cacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

var cache = make(map[string]cacheItem)

// Gemini API configuration
const GEMINI_API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"

// Thêm dòng này để lấy API key từ biến môi trường
var GEMINI_API_KEY = os.Getenv("GEMINI_API_KEY")

// EnglishLevel enum
var englishLevels = map[int]string{
	1: "A1 - Beginner",
	2: "A2 - Elementary",
	3: "B1 - Intermediate",
	4: "B2 - Upper Intermediate",
	5: "C1 - Advanced",
	6: "C2 - Proficient",
}

// AssignmentType enum
var assignmentTypes = map[int]string{
	1: "Multiple Choice",
	2: "Fill in the Blank",
	3: "Short Answer",
	4: "Essay",
}

// Difficulty mapping for different English levels
var difficultyMapping = map[string]string{
	"A1 - Beginner":           "very basic vocabulary and simple grammar structures",
	"A2 - Elementary":         "basic vocabulary with simple past and present tenses",
	"B1 - Intermediate":       "intermediate vocabulary with complex sentence structures",
	"B2 - Upper Intermediate": "advanced vocabulary with sophisticated grammar",
	"C1 - Advanced":           "complex vocabulary with nuanced language usage",
	"C2 - Proficient":         "expert-level vocabulary with native-like complexity",
}

// --- MAIN HANDLER ---

func GenerateAssignment(w http.ResponseWriter, r *http.Request) {
	var request GenerateQuizzesRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validation
	if err := validateRequest(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check cache
	cacheKey := generateCacheKey(request)
	now := time.Now()
	if item, found := cache[cacheKey]; found && item.ExpiresAt.After(now) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item.Data)
		return
	}

	// Generate quizzes using Gemini API
	quizResponse, err := generateQuizzesWithGemini(request)
	if err != nil {
		log.Printf("Error generating quizzes: %v", err)
		http.Error(w, "Failed to generate quizzes", http.StatusInternalServerError)
		return
	}

	// Cache for 10 minutes
	cache[cacheKey] = cacheItem{Data: quizResponse, ExpiresAt: now.Add(10 * time.Minute)}

	log.Printf("Generated %d quizzes for topic: %s", len(quizResponse.Quizzes), request.Topic)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quizResponse)
}

// Validate request parameters
func validateRequest(request GenerateQuizzesRequest) error {
	request.Topic = strings.TrimSpace(request.Topic)
	if request.Topic == "" {
		return errors.New("tên chủ đề không được để trống")
	}
	if len(strings.Fields(request.Topic)) > 10 {
		return errors.New("chủ đề không được chứa nhiều hơn 10 từ")
	}
	if request.TotalQuestions < 1 || request.TotalQuestions > 50 {
		return errors.New("số lượng câu hỏi phải nằm trong khoảng 1 đến 50")
	}
	if len(request.AssignmentTypes) > request.TotalQuestions {
		return errors.New("số lượng câu hỏi không được nhỏ hơn số dạng câu hỏi mà bạn chọn")
	}
	if len(request.AssignmentTypes) == 0 {
		return errors.New("phải chọn ít nhất một loại câu hỏi")
	}
	return nil
}

// Generate quizzes using Gemini API
func generateQuizzesWithGemini(req GenerateQuizzesRequest) (*QuizResponse, error) {
	// Build prompt for Gemini
	prompt := buildGeminiPrompt(req)

	// Call Gemini API
	geminiResp, err := callGeminiAPI(prompt)
	if err != nil {
		return nil, fmt.Errorf("gemini API call failed: %w", err)
	}

	// Parse response
	quizzes, err := parseGeminiResponse(geminiResp, req.AssignmentTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}

	// Ensure we have the right number of questions
	if len(quizzes) < req.TotalQuestions {
		// If we don't have enough, try to generate more
		additionalQuizzes, err := generateAdditionalQuizzes(req, len(quizzes))
		if err == nil {
			quizzes = append(quizzes, additionalQuizzes...)
		}
	}

	// Limit to requested number
	if len(quizzes) > req.TotalQuestions {
		quizzes = quizzes[:req.TotalQuestions]
	}

	// Add IDs to quizzes
	for i := range quizzes {
		quizzes[i].ID = i + 1
	}

	response := &QuizResponse{
		Topic:     req.Topic,
		Level:     req.EnglishLevel,
		Total:     req.TotalQuestions,
		Generated: len(quizzes),
		Quizzes:   quizzes,
	}

	return response, nil
}

// Build comprehensive prompt for Gemini
func buildGeminiPrompt(req GenerateQuizzesRequest) string {
	difficulty, exists := difficultyMapping[req.EnglishLevel]
	if !exists {
		difficulty = "intermediate level"
	}

	typeDistribution := distributeQuestionTypes(req.AssignmentTypes, req.TotalQuestions)

	prompt := fmt.Sprintf(`Create %d high-quality quiz questions about "%s" for %s English level students.

REQUIREMENTS:
- English Level: %s (%s)
- Topic: %s
- Total Questions: %d
- Each question must be unique and non-repetitive
- Questions should be similar in style to IELTS/TOEIC exams
- Include detailed explanations for answers

QUESTION DISTRIBUTION:
%s

FORMATTING RULES:
- Return ONLY valid JSON without any markdown formatting or code blocks
- Use this exact JSON structure:
{
  "quizzes": [
    {
      "type": "Multiple Choice",
      "question": "question text here",
      "options": ["A", "B", "C", "D"],
      "correct_index": 0,
      "explanation": "detailed explanation"
    },
    {
      "type": "Fill in the Blank",
      "question": "Complete this sentence: The weather today is _____ than yesterday.",
      "answer": "better",
      "explanation": "explanation here"
    },
    {
      "type": "Short Answer",
      "question": "question text here",
      "answer": "expected answer",
      "explanation": "explanation here" 
    },
    {
      "type": "Essay",
      "question": "essay question here",
      "answer": "sample key points or structure",
      "explanation": "grading criteria and expectations"
    }
  ]
}

QUALITY STANDARDS:
- Multiple Choice: 4 options, only one correct, plausible distractors
- Fill in the Blank: Clear context, single correct answer
- Short Answer: Specific, measurable expected responses
- Essay: Clear prompts with specific requirements
- All questions must test different aspects of the topic
- Vary sentence structures and vocabulary within the appropriate level
- Include practical, real-world applications when possible

Generate exactly %d questions now:`,
		req.TotalQuestions, req.Topic, req.EnglishLevel, req.EnglishLevel, difficulty, req.Topic, req.TotalQuestions,
		formatTypeDistribution(typeDistribution), req.TotalQuestions)

	return prompt
}

// Distribute question types evenly
func distributeQuestionTypes(types []string, total int) map[string]int {
	distribution := make(map[string]int)
	baseCount := total / len(types)
	remainder := total % len(types)

	for i, questionType := range types {
		distribution[questionType] = baseCount
		if i < remainder {
			distribution[questionType]++
		}
	}

	return distribution
}

// Format type distribution for prompt
func formatTypeDistribution(dist map[string]int) string {
	var parts []string
	for qType, count := range dist {
		parts = append(parts, fmt.Sprintf("- %s: %d questions", qType, count))
	}
	return strings.Join(parts, "\n")
}

// Call Gemini API using SDK
func callGeminiAPI(prompt string) (string, error) {
	client := internal.GeminiClient
	if client == nil {
		return "", errors.New("Gemini client not initialized")
	}
	ctx := context.Background()
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash", // hoặc "gemini-1.5-pro" nếu bạn muốn
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

// Parse Gemini response into Quiz structures
func parseGeminiResponse(response string, requestedTypes []string) ([]Quiz, error) {
	// Clean the response - remove any markdown formatting
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var geminiData GeminiQuizData
	if err := json.Unmarshal([]byte(response), &geminiData); err != nil {
		log.Printf("Failed to parse JSON response: %s", response)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var quizzes []Quiz
	for _, gQuiz := range geminiData.Quizzes {
		quiz := Quiz{
			Type:         gQuiz.Type,
			Question:     strings.TrimSpace(gQuiz.Question),
			Answer:       strings.TrimSpace(gQuiz.Answer),
			Options:      gQuiz.Options,
			CorrectIndex: gQuiz.CorrectIndex,
			Explanation:  strings.TrimSpace(gQuiz.Explanation),
		}

		// Validate question type
		if !contains(requestedTypes, quiz.Type) {
			continue
		}

		// Validate required fields based on type
		if !isValidQuiz(quiz) {
			continue
		}

		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

// Validate quiz based on its type
func isValidQuiz(quiz Quiz) bool {
	if quiz.Question == "" {
		return false
	}

	switch quiz.Type {
	case "Multiple Choice":
		return len(quiz.Options) >= 2 && quiz.CorrectIndex >= 0 && quiz.CorrectIndex < len(quiz.Options)
	case "Fill in the Blank":
		return quiz.Answer != ""
	case "Short Answer":
		return quiz.Answer != ""
	case "Essay":
		return true // Essay questions just need a question
	default:
		return false
	}
}

// Generate additional quizzes if needed
func generateAdditionalQuizzes(req GenerateQuizzesRequest, currentCount int) ([]Quiz, error) {
	needed := req.TotalQuestions - currentCount
	if needed <= 0 {
		return nil, nil
	}

	// Create a new request for the additional questions
	additionalReq := req
	additionalReq.TotalQuestions = needed

	prompt := fmt.Sprintf(`Generate %d additional unique quiz questions about "%s" for %s level. 
Make sure these questions are completely different from any previous questions about this topic.
Focus on different aspects, use different vocabulary, and vary the question formats.

Use the same JSON format as before and ensure high quality, IELTS/TOEIC-style questions.`,
		needed, req.Topic, req.EnglishLevel)

	response, err := callGeminiAPI(prompt)
	if err != nil {
		return nil, err
	}

	return parseGeminiResponse(response, req.AssignmentTypes)
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SuggestTopics suggests random topics for quizzes.
func SuggestTopics(w http.ResponseWriter, r *http.Request) {
	topics := []string{
		"Business Communication", "Environmental Science", "Technology Innovation",
		"Global Economics", "Cultural Diversity", "Health and Wellness",
		"Digital Marketing", "Sustainable Development", "Artificial Intelligence",
		"International Relations", "Climate Change", "Social Media Impact",
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(topics), func(i, j int) { topics[i], topics[j] = topics[j], topics[i] })

	suggestedTopics := topics[:5] // Return 5 random topics

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"topics": suggestedTopics,
	})
}

// GET /api/assignment/get-english-levels
func GetEnglishLevels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(englishLevels)
}

// GET /api/assignment/get-assignment-types
func GetAssignmentTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assignmentTypes)
}

// Helper function to generate cache key.
func generateCacheKey(req GenerateQuizzesRequest) string {
	return strings.ToLower(req.Topic) + "-" + strings.Join(req.AssignmentTypes, "-") + "-" + req.EnglishLevel + "-" + strconv.Itoa(req.TotalQuestions)
}
