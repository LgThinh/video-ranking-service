# Video Ranking Microservice

Dự án này là một microservice để xếp hạng video dựa trên các tương tác như lượt xem, lượt thích, bình luận, chia sẻ và thời gian xem. Hệ thống hỗ trợ xếp hạng toàn cầu và cá nhân hóa cho người dùng và người sáng tạo.

## Tính năng
- Cập nhật điểm số video real-time dựa trên các tương tác người dùng.
- Xếp hạng video toàn cầu.
- Xếp hạng video cá nhân hóa cho từng người dùng hoặc người sáng tạo.
- Cung cấp API để lấy danh sách video xếp hạng.

## Cài đặt

1. Clone repository về máy:
    ```bash
    git clone https://github.com/username/video-ranking-microservice.git
    cd video-ranking-microservice
    ```

2. Cài đặt các phụ thuộc:
    ```bash
    go mod tidy
    ```

3. Cấu hình các biến môi trường:
    - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
    - `REDIS_HOST`, `REDIS_PORT`

4. Chạy dịch vụ:
    ```bash
    go run main.go
    ```

Dịch vụ sẽ chạy tại `http://localhost:8080`.

## Cách sử dụng

### Cập nhật điểm số video

Sử dụng API `PUT /score/update` để cập nhật điểm số video.

Ví dụ:
```bash
curl -X PUT http://localhost:8080/score/update -d '{
  "video_id": "12345",
  "views": 100,
  "likes": 50
}'
