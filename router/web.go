package router

import (
    "net/http"
)

func RegisterWeb() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Web Endpoint"))
    })
}