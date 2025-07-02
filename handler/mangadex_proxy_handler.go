package handler

import (
	"io"
	"net/http"
	"strings"
)

func MangaDexProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Nếu path bắt đầu bằng /api/mangadex/uploads thì proxy sang uploads.mangadex.org
	if strings.HasPrefix(r.URL.Path, "/api/mangadex/uploads") {
		proxyPath := strings.TrimPrefix(r.URL.Path, "/api/mangadex/uploads")
		targetURL := "https://uploads.mangadex.org" + proxyPath
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}

		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to reach MangaDex Uploads", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	// Mặc định: proxy sang api.mangadex.org như cũ
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
