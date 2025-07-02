package handler

import (
	"io"
	"net/http"
	"strings"
)

func MangaDexProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy phần path sau /api/mangadex
	proxyPath := strings.TrimPrefix(r.URL.Path, "/api/mangadex")
	targetURL := "https://api.mangadex.org" + proxyPath
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	// Tạo request mới
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header

	// Gửi request đến MangaDex
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach MangaDex", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy header và body về cho client
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
