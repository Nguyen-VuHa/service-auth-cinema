package constants

// Khai báo const cho các trạng thái trong hệ thống
const (
	// mã 200
	STATUS_SUCCESS = "OK"

	// mã 400
	STATUS_BAD_REQUEST   = "BAD_REQUEST"
	STATUS_INVALID_FIELD = "INVALID_FIELD" // tương ứng với UNPROCESSABLE_ENTITY

	// mã 500
	STATUS_SERVER_INTERNAL_ERROR = "SERVER_INTERNAL_ERROR"
)

// khai báo status code cho các trạng thái trong hệ thống
const (
	CODE_SUCCESS = 200

	CODE_BAD_REQUEST   = 400
	CODE_INVALID_FIELD = 422

	CODE_SERVER_INTERNAL_ERROR = 500
)
