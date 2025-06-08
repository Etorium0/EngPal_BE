package router

import (
	"EngPal/handler"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Assignment routes
	r.HandleFunc("/api/assignment/generate", handler.GenerateAssignment).Methods("POST")
	r.HandleFunc("/api/assignment/suggest-topics", handler.SuggestTopics).Methods("GET")

	// Healthcheck routes
	r.HandleFunc("/api/healthcheck", handler.Healthcheck).Methods("GET")

	// Review routes
	r.HandleFunc("/api/review/generate", handler.GenerateReview).Methods("POST")

	// Chatbot routes
	r.HandleFunc("/api/chatbot/generate-answer", handler.GenerateAnswer).Methods("POST")

	return r
}