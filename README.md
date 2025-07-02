# EngPal_BE

**EngPal_BE** là backend phục vụ cho ứng dụng học tiếng Anh EngPal, được phát triển bằng ngôn ngữ Golang. Dự án này cung cấp API cho các tính năng học tập, kiểm tra và quản lý người dùng của hệ thống EngPal.

## Mục tiêu dự án

- Xây dựng hệ thống backend hiệu quả, dễ mở rộng.
- Cung cấp API RESTful cho ứng dụng EngPal.
- Hỗ trợ quản lý người dùng, bài học, kiểm tra và thống kê tiến độ học tập.

## Cấu trúc dự án

```
EngPal_BE/
├── assets/: Tài nguyên tĩnh, dữ liệu mẫu.
├── entities/: Định nghĩa các kiểu dữ liệu như assignment_type, english_level.
├── handler/: Xử lý các logic API như authentication, assignment, chatbot, review, translate, notification...
├── internal/: Mã nguồn chỉ dùng nội bộ dự án.
├──── config/: Cấu hình hệ thống.
├──── migrations/: Các file migration SQL (tạo user, notify, push token, 
├── router/: Định tuyến API.
└── README.md       # Tài liệu dự án
```

## Hướng dẫn cài đặt & chạy dự án

1. **Yêu cầu:** Cài đặt Golang (>=1.18).

2. **Clone repository:**
   ```bash
   git clone <repo-url>
   cd EngPal_BE
   ```

3. **Cài đặt dependencies (nếu cần):**
   ```bash
   go mod tidy
   ```

4. **Khởi động server:**
   ```bash
   go run cmd/main.go
   ```

5. **API sẽ chạy tại:**  
   `http://localhost:8080`

## Đóng góp

- Fork repository, tạo branch mới từ nhánh chính.
- Commit các thay đổi và gửi pull request.
- Mọi đóng góp đều được hoan nghênh!

## Liên hệ & hỗ trợ

- Vui lòng tạo issue trên GitHub nếu cần hỗ trợ hoặc báo lỗi.
