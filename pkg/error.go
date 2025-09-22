package pkg

import (
	"fmt"
	"log"
	"runtime"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrorResponse struct {
	Error       ErrorDetails      `json:"error"`
	Validations validation.Errors `json:"validations,omitempty"`
}

type ErrorDetails struct {
	Message string `json:"msg"`
	Details any    `json:"details"`
}

// Error represents an error that could be wrapping another error, it includes a code for determining what
// triggered the error.
type Error struct {
	orig error
	msg  string
	code ErrorCode
}

// ErrorCode defines supported error codes.
type ErrorCode uint

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeNotFound
	ErrorCodeInvalidArgument
	ErrNoDocuments
	ErrNoChange
	ErrNoRows
)

func WrapIfErr[T any](res T, err error, msg string) (T, error) {
	if err != nil {
		return *new(T), WrapErrorf(err, ErrorCodeUnknown, msg)
	}
	return res, nil
}

// WrapErrorf returns a wrapped error.
func WrapErrorf(orig error, code ErrorCode, format string, a ...interface{}) error {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("ERROR: [%s:%d] %v \n", file, line, orig)
	return &Error{
		code: code,
		orig: orig,
		msg:  fmt.Sprintf(format, a...),
	}
}

// NewErrorf instantiates a new error.
func NewErrorf(code ErrorCode, format string, a ...interface{}) error {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("ERROR: [%s:%d] %v \n", file, line, format)
	return WrapErrorf(nil, code, format, a...)
}

// Error returns the message, when wrapping errors the wrapped error is returned.
func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}

	return e.msg
}

// Unwrap returns the wrapped error, if any.
func (e *Error) Unwrap() error {
	return e.orig
}

// Code returns the code representing this error.
func (e *Error) Code() ErrorCode {
	return e.code
}

// func WrapErrLang(c *fiber.Ctx, orig error, code ErrorCode, format string, a ...interface{}) error {
// 	localize, err := fiberi18n.Localize(c, format)
// 	if err != nil {
// 		return err
// 	}
// 	return &Error{
// 		code: code,
// 		orig: orig,
// 		msg:  fmt.Sprintf(localize, a...),
// 	}
// }
// func NewErrLang(c *fiber.Ctx, code ErrorCode, format string, a ...interface{}) error {
// 	localize, err := fiberi18n.Localize(c, format)
// 	if err != nil {
// 		return err
// 	}
// 	return WrapErrorf(nil, code, localize, a...)
// }
