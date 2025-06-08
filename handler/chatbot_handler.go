package handler
 
import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"EngPal/utils"
)

// Placeholder types for demonstration.
type Conversation struct {
	Question string `json:"question"`
}

type ChatResponse struct {
	MessageInMarkdown string `json:"message_in_markdown"`
}

// GenerateAnswer handles chatbot question processing and response generation.
func GenerateAnswer(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request into `Conversation`.
	var request Conversation
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Placeholder for additional parameters
	username := r.URL.Query().Get("username")
	gender := r.URL.Query().Get("gender")
	age := r.URL.Query().Get("age")
	englishLevel := r.URL.Query().Get("english_level")
	enableReasoning := r.URL.Query().Get("enable_reasoning") == "true"
	enableSearching := r.URL.Query().Get("enable_searching") == "true"

	// Validate the question.
	request.Question = strings.TrimSpace(request.Question)
	if request.Question == "" {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Gửi vội vậy bé yêu! Chưa nhập câu hỏi kìa.",
		})
		return
	}

	if utils.GetTotalWords(request.Question) > 30 {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Hỏi ngắn thôi bé yêu, bộ mắc hỏi quá hay gì 💢\nHỏi câu nào dưới 30 từ thôi, để thời gian cho anh suy nghĩ với chứ.",
		})
		return
	}

	// Generate chatbot response.
	result, err := generateChatbotResponse(request, username, gender, age, englishLevel, enableReasoning, enableSearching)
	if err != nil {
		log.Printf("Error generating answer: %v", err)
		json.NewEncoder(w).Encode(ChatResponse{
			MessageInMarkdown: "Nhắn từ từ thôi bé yêu, bộ mắc đi đẻ quá hay gì 💢\nNgồi đợi 1 phút cho anh đi uống ly cà phê đã. Sau 1 phút mà vẫn lỗi thì xóa lịch sử trò chuyện rồi thử lại nha!",
		})
		return
	}

	// Log the successful response.
	log.Printf("%s (%s) asked (Reasoning: %v - Grounding: %v): %s", "access-key", username, enableReasoning, enableSearching, request.Question)

	// Send the result back to the client.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// Simulate chatbot response generation.
func generateChatbotResponse(request Conversation, username, gender, age, englishLevel string, enableReasoning, enableSearching bool) (ChatResponse, error) {
	// Placeholder logic for generating chatbot response.
	if strings.Contains(request.Question, "error") {
		return ChatResponse{}, errors.New("error generating response")
	}
	return ChatResponse{
		MessageInMarkdown: "Đây là câu trả lời mẫu từ chatbot! 🚀",
	}, nil
}

