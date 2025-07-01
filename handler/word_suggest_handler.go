package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"fmt"
	"math/rand"
)

type VocabSense struct {
	Definition string `json:"definition"`
	Examples   []struct {
		Cf string `json:"cf"`
		X  string `json:"x"`
	} `json:"examples"`
}

type VocabEntry struct {
	Word           string        `json:"word"`
	Pos            string        `json:"pos"`
	Phonetic       string        `json:"phonetic"`
	PhoneticText   string        `json:"phonetic_text"`
	PhoneticAm     string        `json:"phonetic_am"`
	PhoneticAmText string        `json:"phonetic_am_text"`
	Senses         []VocabSense  `json:"senses"`
}

// Đường dẫn thư mục json từ vựng (tạo nếu chưa có)
const vocabDir = "assets/json/oxford_words"

// Handler gợi ý 5 từ vựng ngẫu nhiên mỗi ngày, không cần query đầu vào
func SuggestWordsHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy danh sách tất cả file json trong vocabDir
	files, err := ioutil.ReadDir(vocabDir)
	if err != nil || len(files) == 0 {
		// Nếu không có file json nào, gọi Gemini lấy 5 từ mới
		words := FetchWordsFromGemini("", 5)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(words)
		return
	}

	var allEntries []VocabEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			data, err := ioutil.ReadFile(filepath.Join(vocabDir, file.Name()))
			if err == nil {
				var entries []VocabEntry
				if err := json.Unmarshal(data, &entries); err == nil {
					allEntries = append(allEntries, entries...)
				}
			}
		}
	}

	// Lấy ngẫu nhiên 5 từ
	suggestions := []VocabEntry{}
	if len(allEntries) > 0 {
		randomIdx := rand.Perm(len(allEntries))
		for _, idx := range randomIdx {
			suggestions = append(suggestions, allEntries[idx])
			if len(suggestions) == 5 {
				break
			}
		}
	}

	// Nếu chưa đủ 5 từ, gọi Gemini bổ sung
	if len(suggestions) < 5 {
		missing := 5 - len(suggestions)
		newWords := FetchWordsFromGemini("", missing)
		suggestions = append(suggestions, newWords...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// Hàm gọi Gemini API để lấy từ mới (thực tế)
func FetchWordsFromGemini(q string, count int) []VocabEntry {
	prompt := buildGeminiSuggestPrompt(q, count)
	response, err := callGeminiAPI(prompt)
	if err != nil {
		// Nếu lỗi, trả về rỗng hoặc log lỗi
		return []VocabEntry{}
	}
	// Xử lý response: loại bỏ markdown nếu có
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var entries []VocabEntry
	if err := json.Unmarshal([]byte(response), &entries); err != nil {
		// Nếu lỗi parse, trả về rỗng hoặc log lỗi
		return []VocabEntry{}
	}
	return entries
}

// Hàm tạo prompt cho Gemini
func buildGeminiSuggestPrompt(q string, count int) string {
	return fmt.Sprintf(`Suggest exactly %d English words related to or similar to "%s". For each word, return a JSON object with the following fields: word, pos, phonetic, phonetic_text, phonetic_am, phonetic_am_text, senses (array of objects with definition and examples (array of {cf, x})).

Example output:
[
  {
    "word": "example",
    "pos": "noun",
    "phonetic": "https://...",
    "phonetic_text": "/.../",
    "phonetic_am": "https://...",
    "phonetic_am_text": "/.../",
    "senses": [
      {
        "definition": "...",
        "examples": [
          {"cf": "", "x": "..."}
        ]
      }
    ]
  }
]
Return only valid JSON array, no explanation.`, count, q)
} 