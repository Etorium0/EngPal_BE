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

	// Translate route
	r.HandleFunc("/api/translate", handler.TranslateHandler).Methods("POST")

	// Word routes
	r.HandleFunc("/api/word/suggest", handler.SuggestWordsHandler).Methods("GET")

	// Auth routes
	userRepo := repo_impl.NewUserRepo(db)
	authHandler := &handler.AuthenticationHandler{UserRepo: userRepo}
	r.HandleFunc("/api/login", authHandler.ServeLogin).Methods("POST")
	r.HandleFunc("/api/register", authHandler.ServeRegister).Methods("POST")
	r.HandleFunc("/api/forgot-password", authHandler.ServeForgotPassword).Methods("POST")

	// Notification routes
	r.HandleFunc("/api/user/notification-setting", handler.UpdateNotificationSettingHandler(db)).Methods("POST")
	r.HandleFunc("/api/user/notification-setting", handler.GetNotificationSettingHandler(db)).Methods("GET")
	r.HandleFunc("/api/user/notification-history", handler.GetNotificationHistoryHandler(db)).Methods("GET")

	// New route
	r.HandleFunc("/api/notify/test", handler.SendTestNotification).Methods("POST")

	// Push token route
	r.HandleFunc("/api/user/push-token", handler.SavePushTokenHandler(db)).Methods("POST")

	// MangaDex proxy route (API và ảnh uploads)
	r.PathPrefix("/api/mangadex/").HandlerFunc(handler.MangaDexProxyHandler)

	return r
}