package router

import (
    "net/http"
)

func RegisterAPI() {
    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("API Endpoint"))
    })
}