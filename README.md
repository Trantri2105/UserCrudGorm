# User CRUD

## Overview
- Dự án cung cấp API để quản lý người dùng (user management). Project sử dụng:
    - Ngôn ngữ `Go`.
    - Thư viện `GORM` để tương tác với `PostgreSQL`.
    - `JWT` để xác thực.

## Project structure

```
UserCrud/
├── dto/
│   ├── request/ # Request DTOs
│   └── response/ # Response DTOs
├── handler/ # HTTP handlers
├── middleware/ # Middleware components
├── model/ # Data models
├── repository/ # Database operations
├── service/ # Business logic
├── util/ # Utility functions
├── .env # Environment variables
├── Dockerfile # Docker build configuration
├── docker-compose.yml # Docker compose file
├── go.mod # Go modules file
├── go.sum # Go modules checksums
└── main.go # Application entry point
```

## Installation and Setup

### Local
1. Clone project
   ```bash
   git clone https://github.com/Trantri2105/UserCrudGorm.git
   cd UserCrudGorm
   ```
2. Tải các dependencies cần thiết
   ```bash
   go mod tidy
   ```
3. Chạy ứng dụng
   ```bash
   go run .
   ```

### Docker 

- Sử dụng câu lệnh `docker compose up -d` để chạy ứng dụng, server sẽ chạy tại `localhost:8080`

## API Documentation

### Đăng ký một người dùng mới

```
POST /user/register
```

Request body:
```json
{
    "first_name":"Tran Tri",
    "last_name":"Nguyen",
    "email":"trint@gmail.com",
    "password":"123456",
    "phone_number":"1234567890",
    "gender":"male"
}
```

Response
- Nếu đăng ký thành công, server sẽ trả về
```json
201 Created
{
    "message": "User register successfully"
}
```

- Các trường `first_name`, `last_name`, `emai`, `password`, `phone_number`, `gender` đều phải có trong `response body`. Nếu không server sẽ trả về lỗi.
```json
400 Bad Request
{
    "error": "Field Gender is required"
}
```

- Trường `email` phải là email hợp lệ, có dạng `...@abcdxyz.xyz`. Nếu không server sẽ trả về lỗi.
```json
400 Bad Request
{
    "error": "Field Email must be a valid email"
}
```

### Đăng nhập
```
POST /user/login
```
Request Body
```json
{
    "email":"trint@gmail.com",
    "password":"123456"
}
```
Response
- Nếu đăng nhập thành công, server sẽ trả về `access token`
```json
200 Ok
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI..."
}
```
- Các trường `first_name`, `last_name`, `emai`, `password`, `phone_number`, `gender` đều phải có trong `response body`. Nếu không server sẽ trả về lỗi
```json
400 Bad Request
{
    "error": "Field Password is required"
}
```
- Nếu tài khoản chưa tồn tại, server sẽ trả về
```json
404 Not Found
{
    "error": "user not found"
}
```
- Nếu sai password, server sẽ trả về
```json
400 Bad Request
{
    "error": "wrong password"
}
```
### Lấy thông tin người dùng
- Yêu cầu xác thực bằng cách gửi kèm `access token` (lấy ở endpoint đăng nhập) trong header `Authorization`.

```
GET /user/profile
```

Response
```json
200 Ok
{
    "id": 1,
    "first_name": "Tran Tri",
    "last_name": "Nguyen",
    "email": "trint@gmail.com",
    "phone_number": "1234567890",
    "gender": "male",
    "created_at": "2025-03-14T10:57:20.482049Z",
    "updated_at": "2025-03-14T10:57:20.482049Z"
}
```
### Cập nhật thông tin người dùng
- Yêu cầu xác thực bằng cách gửi kèm `access token` (lấy ở endpoint đăng nhập) trong header `Authorization`.

```
PATCH /user
```

Request Body
```json
{
    "phone_number":"12345678910"
}
```

Response
```json
{
    "message": "User update successfully"
}
```

### Xóa người dùng

- Yêu cầu xác thực bằng cách gửi kèm `access token` (lấy ở endpoint đăng nhập) trong header `Authorization`.

```
DELETE /user
```

Response
```json
200 Ok
{
    "message": "User delete successfully"
}
```




