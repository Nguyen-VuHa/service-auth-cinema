package initializers

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger() {
	// Cấu hình log xoay vòng
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./Public/logs/app.log", // khi triển khai thay đổi đến đường dãn path lưu log trong server .
		MaxSize:    10,                      // MB
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	// Thiết lập encoder cho console và file
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// Thiết lập các mức log
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), zapcore.InfoLevel),
	)

	// Tạo logger
	Logger = zap.New(core)
}

func GetLogger() *zap.Logger {
	return Logger
}
