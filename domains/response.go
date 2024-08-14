package domains

type ResponseBasic struct {
	Code    int    `json:"code"`    // mã lỗi khi chương trình chạy
	Status  string `json:"status"`  // trạng thái request (thất bại hoặc thành công)
	Message string `json:"message"` // Message thông báo lỗi hoặc thành công
}
