package DTO

// Type request sign up account
type SignUpRequest struct {
	Email       string `json:"email"` // trả ra JSON với định dạng email thay vì Email
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	BirthDay    string `json:"birth_day"`
	PhoneNumber string `json:"phone_number"`
}
