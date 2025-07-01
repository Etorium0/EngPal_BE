package handler
 
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"EngPal/utils"
)

// Request/Response types
// Đặt lại tên cho rõ ràng
// Chatbot request
type ChatbotRequest struct {
	Question string `json:"question"`
}

// Chatbot response
type ChatbotResponse struct {
	Message string `json:"message"`
}

// Validate chatbot request
func validateChatbotRequest(request ChatbotRequest) error {
	request.Question = strings.TrimSpace(request.Question)
	if request.Question == "" {
		return errors.New("Gửi vội vậy bé yêu! Chưa nhập câu hỏi kìa.")
	}
	if utils.GetTotalWords(request.Question) > 30 {
		return errors.New("Hỏi ngắn thôi bé yêu, bộ mắc hỏi quá hay gì 💢\nHỏi câu nào dưới 30 từ thôi, để thời gian cho anh suy nghĩ với chứ.")
	}
	return nil
}

// GenerateAnswer handles chatbot question processing and response generation.
func GenerateAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the incoming JSON request into `ChatbotRequest`.
	var request ChatbotRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validateChatbotRequest(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Lấy các tham số phụ nếu cần
	username := r.URL.Query().Get("username")
	gender := r.URL.Query().Get("gender")
	age := r.URL.Query().Get("age")
	englishLevel := r.URL.Query().Get("english_level")
	enableReasoning := r.URL.Query().Get("enable_reasoning") == "true"
	enableSearching := r.URL.Query().Get("enable_searching") == "true"

	// Generate chatbot response.
	result, err := generateChatbotResponse(request, username, gender, age, englishLevel, enableReasoning, enableSearching)
	if err != nil {
		log.Printf("Error generating answer: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "service_unavailable",
			"message": "Nhắn từ từ thôi bé yêu, bộ mắc đi đẻ quá hay gì 💢\nNgồi đợi 1 phút cho anh đi uống ly cà phê đã. Sau 1 phút mà vẫn lỗi thì xóa lịch sử trò chuyện rồi thử lại nha!",
		})
		return
	}

	// Log the successful response.
	log.Printf("[Chatbot] %s (%s) asked (Reasoning: %v - Grounding: %v): %s", "access-key", username, enableReasoning, enableSearching, request.Question)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Simulate chatbot response generation.
func generateChatbotResponse(request ChatbotRequest, username, gender, age, englishLevel string, enableReasoning, enableSearching bool) (ChatbotResponse, error) {
	answer, err := callGeminiForChatbot(request.Question)
	if err != nil {
		return ChatbotResponse{}, err
	}
	return ChatbotResponse{Message: answer}, nil
}

func callGeminiForChatbot(question string) (string, error) {
	// Tạo prompt cho Gemini
	prompt := fmt.Sprintf("Answer this question in a friendly, helpful way: %s", question)
	// Gọi Gemini API (tái sử dụng logic từ handler khác)
	// Ví dụ:
	response, err := callGeminiAPI(prompt) // callGeminiAPI là hàm đã có sẵn
	if err != nil {
		return "", err
	}
	// Xử lý response nếu cần
	return response, nil
}

