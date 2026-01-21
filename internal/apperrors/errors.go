package apperrors

import (
	"errors"
	"fmt"
)

// Сентинельные ошибки — можно сравнивать через errors.Is
var (
	ErrNotFound   = errors.New("not found")
	ErrPermission = errors.New("permission denied")
	ErrTimeout    = errors.New("timeout")
)

// ValidationError — конкретная ошибка валидации поля
type ValidationError struct {
	Field   string
	Message string
}

func (v *ValidationError) Error() string {
	if v.Field != "" {
		return fmt.Sprintf("validation failed: %s: %s", v.Field, v.Message)
	}
	return fmt.Sprintf("validation failed: %s", v.Message)
}
