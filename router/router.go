package router

import (
	"EngPal/handler"
	"EngPal/repository/repo_impl"

	"database/sql"

	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Assignment routes
	r.HandleFunc("/api/assignment/generate", handler.GenerateAssignment).Methods("POST")
	r.HandleFunc("/api/assignment/suggest-topics", handler.SuggestTopics).Methods("GET")

	// Review routes
	r.HandleFunc("/api/review/generate", handler.GenerateReview).Methods("POST")

	// Chatbot routes
	r.HandleFunc("/api/chatbot/answer", handler.GenerateAnswer).Methods("POST")

	// Auth routes
	userRepo := repo_impl.NewUserRepo(db)
	authHandler := &handler.AuthenticationHandler{UserRepo: userRepo}
	r.HandleFunc("/api/login", authHandler.ServeLogin).Methods("POST")
	r.HandleFunc("/api/register", authHandler.ServeRegister).Methods("POST")
	r.HandleFunc("/api/forgot-password", authHandler.ServeForgotPassword).Methods("POST")

	return r
}