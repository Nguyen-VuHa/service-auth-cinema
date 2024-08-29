package domains

type ServiceMailRepository interface {
	SendOTPCodeToMail(params map[string]interface{}) error
}
