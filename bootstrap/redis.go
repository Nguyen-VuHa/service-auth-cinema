package bootstrap

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisAuth *redis.Client // Khai báo biến global cho Redis client liên quan đến xác thực
var RedisUser *redis.Client // Khai báo biến global cho Redis client liên quan đến người dùng

// Hàm khởi tạo kết nối tới Redis
func InitRedis(addr, userName, password string, db int) *redis.Client {
	var Redis *redis.Client // Khai báo biến Redis client

	var ctx = context.Background() // Tạo context mặc định

	// Khởi tạo Redis client với các thông tin kết nối
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,     // Địa chỉ của Redis server
		Username: userName, // Username (nếu có)
		Password: password, // Mật khẩu (nếu có)
		DB:       db,       // Chỉ định database sử dụng (0 là mặc định)
	})

	// Thử kết nối đến Redis bằng lệnh PING
	err := Redis.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err) // Log lỗi và thoát nếu không thể kết nối
	}
	log.Println("Connected to Redis successfully") // Log thành công nếu kết nối được

	return Redis // Trả về Redis client đã khởi tạo
}

// Hàm dọn dẹp để đóng kết nối Redis
func CleanupRedis(Redis *redis.Client) {
	log.Println("Closing Redis connection...") // Log thông báo đóng kết nối
	Redis.Close()                              // Đóng kết nối Redis
}
