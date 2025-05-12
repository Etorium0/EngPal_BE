package main

import (
	"log"
	"net/http"

	"course-golang/router"
)

func main() {
	r := router.SetupRouter()

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}