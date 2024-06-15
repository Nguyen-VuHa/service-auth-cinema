package helpers

import "errors"

func CreateMessageError(message string) error {
	return errors.New(message)
}
