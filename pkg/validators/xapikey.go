package validators

import (
	"errors"
	"unicode"
)

func ValidateXApiKey(token string) error {
	if len(token) != 64 {
		return errors.New("token must be of 64 characters long")
	}
	for _, char := range token {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			continue
		}
		return errors.New("token must contain only letters")
	}
	return nil
}
