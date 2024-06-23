package helpers

import (
	"time"
)

func StringToDate(dateStr string) (time.Time, error) {
	// Định dạng thời gian phải theo chuẩn thời gian "YYYY-MM-DD"
	layoutFormat := "2006-01-02"

	// Chuyển đổi chuỗi sang time.Time
	timeConvert, err := time.Parse(layoutFormat, dateStr)

	if err != nil { // trường hợp convert failed
		return timeConvert, err
	}

	return timeConvert, nil
}
