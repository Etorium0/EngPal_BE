# EngAce-Go

Backend cho ứng dụng học tiếng Anh EngAce, được triển khai bằng Golang.

## Cấu trúc dự án
```
EngAce-Go/
├── cmd/            # Tệp lệnh chính, entry point
├── pkg/            # Logic chính của ứng dụng
├── internal/       # Các gói chỉ dùng trong dự án
├── go.mod          # Module Go
├── go.sum          # Phụ thuộc
└── README.md
```

## Cách chạy dự án
1. Cài đặt Golang (>=1.18) trên máy của bạn.
2. Clone repo này:
   ```bash
   git clone <repo-url>
   cd EngAce-Go
   ```
3. Chạy lệnh sau để khởi động server:
   ```bash
   go run cmd/main.go
   ```
4. API sẽ hoạt động tại `http://localhost:8080`.