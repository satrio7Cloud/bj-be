package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Gunakan satu instance validator saja untuk cache struct info agar lebih cepat
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// FormatValidationError mengubah error asli dari validator menjadi array string
// yang lebih mudah dibaca oleh user/frontend
func FormatValidationError(err error) []string {
	var errorsMsg []string
	if err != nil {
		// Pengecekan apakah error merupakan ValidationErrors
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, err := range validationErrors {
				errorsMsg = append(errorsMsg, fmt.Sprintf("Field '%s' tidak memenuhi aturan validasi: '%s'", err.Field(), err.Tag()))
			}
		} else {
			// Jika error lain (misal tipe data salah saat unmarshal)
			errorsMsg = append(errorsMsg, err.Error())
		}
	}
	return errorsMsg
}

// ValidateStruct memeriksa apakah struct DTO memenuhi aturan tag `validate`
func ValidateStruct(data interface{}) []string {
	err := validate.Struct(data)
	if err != nil {
		return FormatValidationError(err)
	}
	return nil
}
