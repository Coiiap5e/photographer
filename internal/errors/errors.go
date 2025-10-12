package errors

import "fmt"

type ErrorCode string

const (
	// Database errors

	ErrCodeDBConnection ErrorCode = "DB_CONNECTION_ERROR"
	ErrCodeDBQuery      ErrorCode = "DB_QUERY_ERROR"
	ErrCodeDBInsert     ErrorCode = "DB_INSERT_ERROR"
	ErrCodeDBDelete     ErrorCode = "DB_DELETE_ERROR"
	ErrCodeDBSelect     ErrorCode = "DB_SELECT_ERROR"

	// Client operations

	ErrCodeClientNotFound ErrorCode = "CLIENT_NOT_FOUND"
	ErrCodeClientCreate   ErrorCode = "CLIENT_CREATE_ERROR"
	ErrCodeClientDelete   ErrorCode = "CLIENT_DELETE_ERROR"
	ErrCodeClientList     ErrorCode = "CLIENT_LIST_ERROR"

	// Shoot operations

	ErrCodeShootNotFound ErrorCode = "SHOOT_NOT_FOUND"
	ErrCodeShootCreate   ErrorCode = "SHOOT_CREATE_ERROR"
	ErrCodeShootDelete   ErrorCode = "SHOOT_DELETE_ERROR"
	ErrCodeShootList     ErrorCode = "SHOOT_LIST_ERROR"

	// Validation

	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT_ERROR"

	//Configuration

	ErrCodeConfig ErrorCode = "CONFIG_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
