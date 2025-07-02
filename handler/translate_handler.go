package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type TranslateRequest struct {
	Text       string `json:"text"`
	TargetLang string `json:"target_lang"` // "vi" hoặc "en"
}

type TranslateResponse struct {
	Translated string `json:"translated"`
}

func validateTranslateRequest(req TranslateRequest) error {
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		return errors.New("Bạn chưa nhập câu cần dịch!")
	}
	if req.TargetLang != "vi" && req.TargetLang != "en" {
		return errors.New("Ngôn ngữ đích phải là 'vi' hoặc 'en'.")
	}
	if len([]rune(req.Text)) > 500 {
		return errors.New("Câu quá dài, chỉ hỗ trợ tối đa 500 ký tự.")
	}
	return nil
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	if err := validateTranslateRequest(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	prompt := buildTranslatePrompt(req.Text, req.TargetLang)
	translated, err := callGeminiAPI(prompt)
	if err != nil {
		log.Printf("Error translating: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dịch thất bại, thử lại sau!"})
		return
	}
	translated = strings.TrimSpace(translated)
	json.NewEncoder(w).Encode(TranslateResponse{Translated: translated})
}

func buildTranslatePrompt(text, targetLang string) string {
	if targetLang == "vi" {
		return fmt.Sprintf("Dịch câu sau sang tiếng Việt, chỉ trả về bản dịch, không giải thích: %s", text)
	}
	return fmt.Sprintf("Dịch câu sau sang tiếng Anh, chỉ trả về bản dịch, không giải thích: %s", text)
} 