# Sử dụng hình ảnh Golang chính thức làm hình ảnh cơ sở
FROM golang:1.22.4-alpine AS builder

# Đặt thư mục làm việc bên trong container
WORKDIR /app

# Sao chép các tệp go.mod và go.sum và tải xuống các phụ thuộc cần thiết
COPY go.mod go.sum ./
RUN go mod download

# Sao chép mã nguồn còn lại vào container
COPY . .

# Biên dịch ứng dụng Go
RUN go build -o main .

# Sử dụng hình ảnh nhẹ để chứa ứng dụng đã biên dịch
FROM alpine:latest

# Thiết lập thư mục làm việc
WORKDIR /root/

# Sao chép tệp nhị phân đã biên dịch từ giai đoạn xây dựng trước
COPY --from=builder /app/main .

# Mở cổng ứng dụng (ví dụ: 5000)
EXPOSE 5000

# Chạy ứng dụng
CMD ["./main"]
