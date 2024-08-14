package domains

type ValidateRepository interface {
	IsRequireString(data string) error
	IsEmail(data string) error
	IsMaxLengthString(data string, maxLength int) error
	IsRangeLength(data string, minLength, maxLength int) error
	IsPhoneNumber(data string) error
}
