package utils

// Hàm kiểm tra ID
func ItemIsArrayString(id string, ids []string) bool {
	// Tạo map để lưu trữ ID
	idMap := make(map[string]struct{})

	// Thêm tất cả ID vào map
	for _, v := range ids {
		idMap[v] = struct{}{} // struct{} là một kiểu không chứa giá trị
	}

	// Kiểm tra xem ID có tồn tại trong map không
	_, exists := idMap[id]

	return exists // Nếu không tồn tại thì trả về true
}
