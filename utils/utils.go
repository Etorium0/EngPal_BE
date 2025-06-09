package utils

import (
	"os"
	"strings"
	"unicode"
)

// GetTotalWords counts the total number of words in a string.
func GetTotalWords(input string) int {
	return len(strings.Fields(input))
}

// IsEnglish checks if the input contains only English letters, digits, or common punctuation.
func IsEnglish(input string) bool {
	for _, c := range input {
		if unicode.IsLetter(c) && c > unicode.MaxASCII {
			return false
		}
		if c > unicode.MaxASCII && !unicode.IsSpace(c) && !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func getGeminiAPIKey() string {
	return os.Getenv("GEMINI_API_KEY")
}
