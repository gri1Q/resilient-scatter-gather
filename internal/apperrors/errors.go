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
	ErrInternal   = errors.New("internal error")
	ErrBadRequest = errors.New("bad request")
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

// AppError — обёртка для ошибок приложения (без Kind)
type AppError struct {
	// Message — краткое описание ошибки в контексте места, где она возникла
	Message string
	// Err — вложенная (оригинальная) ошибка; может быть nil
	Err error
}

func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Wrap оборачиваем ошибку
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	ae, ok := err.(*AppError)
	if ok {
		return &AppError{
			Message: msg,
			Err:     ae,
		}
	}
	return &AppError{
		Message: msg,
		Err:     err,
	}
}

func New(msg string) error {
	return &AppError{Message: msg}
}
