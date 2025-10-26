package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ==========================
// 🎯 ENUM KODE ERROR
// ==========================
type ErrorCode uint

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeNotFound
	ErrorCodeInvalidArgument
	ErrorCodeUnauthorized
	ErrorCodeConflict
	ErrorCodeInternal
	ErrorCodeBadRequest
)

// ==========================
// ⚙️ STRUKTUR ERROR
// ==========================
type AppError struct {
	Code        ErrorCode
	Message     string
	Orig        error
	Validations validation.Errors
	Expose      bool
}

// ==========================
// 🧱 PEMBUAT ERROR
// ==========================

// WrapError membungkus error internal, tidak diekspos ke client
func WrapError(err error, code ErrorCode, msg string) error {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf("%s (%s:%d)", msg, file, line),
		Orig:    err,
		Expose:  false,
	}
}

// ExposeError membuat error yang bisa ditampilkan ke client
func ExposeError(code ErrorCode, msg string) error {
	return &AppError{
		Code:    code,
		Message: msg,
		Expose:  true,
	}
}

func WrapValidationError(err validation.Errors, message string) *AppError {
	return &AppError{
		Code:        ErrorCodeInvalidArgument,
		Message:     message,
		Validations: err,
	}
}

// ==========================
// 🔍 UTILITAS
// ==========================
func (e *AppError) Error() string {
	if e.Orig != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Orig)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Orig }

// HTTPStatus menentukan kode HTTP dari ErrorCode
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeInvalidArgument:
		return http.StatusBadRequest
	case ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrorCodeConflict:
		return http.StatusConflict
	case ErrorCodeBadRequest:
		return http.StatusBadRequest
	case ErrorCodeUnknown, ErrorCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// AsAppError mencoba cast error biasa menjadi *AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}

// ==========================
// 🧩 STRUKTUR RESPON JSON
// ==========================
type ErrorResponse struct {
	Message     string            `json:"message"`
	Details     any               `json:"details,omitempty"`
	Validations validation.Errors `json:"validations,omitempty"`
}

type ErrorDetails struct {
	Message string `json:"msg"`
	Details any    `json:"details"`
}
