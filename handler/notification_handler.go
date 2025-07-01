package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type NotificationSettingRequest struct {
	Enable bool `json:"enable"`
}

type NotificationSettingResponse struct {
	Enable bool `json:"enable"`
}

// Giả lập lấy user_id từ context (bạn thay bằng auth thực tế)
func getUserIDFromContext(r *http.Request) int {
	return 1 // test cứng user_id = 1
}

// POST /api/user/notification-setting
func UpdateNotificationSettingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r)
		var req NotificationSettingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`INSERT INTO user_notification_setting (user_id, enable_notify, updated_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE SET enable_notify = $2, updated_at = $3`,
			userID, req.Enable, time.Now())
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// GET /api/user/notification-setting
func GetNotificationSettingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r)
		var enable bool
		err := db.QueryRow(`SELECT enable_notify FROM user_notification_setting WHERE user_id = $1`, userID).Scan(&enable)
		if err == sql.ErrNoRows {
			enable = true // mặc định bật nếu chưa có
		} else if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(NotificationSettingResponse{Enable: enable})
	}
}

// Handler gửi thông báo test (giả lập)
type NotifyRequest struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func SendTestNotification(w http.ResponseWriter, r *http.Request) {
	var req NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// TODO: Lấy device token từ DB, gửi push notification qua FCM/OneSignal/Expo...
	// Hiện tại chỉ log ra hoặc trả về JSON thành công
	// log.Printf("Send notify to user %d: %s - %s", req.UserID, req.Title, req.Body)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
}

// Handler lưu Expo Push Token cho user
// POST /api/user/push-token

type PushTokenRequest struct {
	Token string `json:"token"`
}

func SavePushTokenHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromContext(r)
		var req PushTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		_, err := db.Exec(`INSERT INTO user_push_token (user_id, token, updated_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE SET token = $2, updated_at = $3`,
			userID, req.Token, time.Now())
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
} 