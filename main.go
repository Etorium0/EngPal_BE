package main

import (
	"log"
	"net/http"

	"EngPal/internal"
	"EngPal/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}

	internal.InitGeminiClient()

	r := router.SetupRouter()

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
