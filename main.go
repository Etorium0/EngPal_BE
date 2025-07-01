package main

import (
	"log"
	"net/http"

	"EngPal/internal"
	"EngPal/router"
	"EngPal/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}

	internal.InitGeminiClient()

	// Khởi tạo kết nối database
	db, err := config.OpenDBConnection()
	if err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}
	defer db.Close()

	r := router.SetupRouter(db)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
