package constants

// const user status account
const (
	USER_STATUS_PENDING = "pending" // trạng thái chưa xác nhận tài khoản của user
	USER_STATUS_ACTIVE  = "active"  // trạng thái tài khoản có thể sử dụng
	USER_STATUS_HIDDEN  = "hidden"  // trạng thái tài khoản đã xoá
	USER_STATUS_BLOCKED = "blocked" // trạng thái tài khoản bị khoá tạm thời
)

// const user type account
const (
	USER_TYPE_NORMAL  = "normal"  // trạng thái đăng ký trực tiếp trên website hệ thống
	USER_TYPE_ANOTHER = "another" // trạng thái đăng ký ở các platform khác
)

// const key value trong bảng userProfile
const (
	USER_PROFILE_FULLNAME    = "full_name"
	USER_PROFILE_BIRTHDAY    = "birth_day"
	USER_PROFILE_PHONENUMBER = "phone_number"
)
