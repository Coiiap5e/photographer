package errors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	// Database errors

	ErrCodeDBConnection  ErrorCode = "DB_CONNECTION_ERROR"
	ErrCodeDBQuery       ErrorCode = "DB_QUERY_ERROR"
	ErrCodeDBInsert      ErrorCode = "DB_INSERT_ERROR"
	ErrCodeDBDelete      ErrorCode = "DB_DELETE_ERROR"
	ErrCodeDBSelect      ErrorCode = "DB_SELECT_ERROR"
	ErrCodeDBTransaction ErrorCode = "DB_TRANSACTION_ERROR"
	ErrCodeDBScan        ErrorCode = "DB_SCAN_ERROR"
	ErrCodeDBConfig      ErrorCode = "DB_CONFIG_ERROR"

	// Client operations

	ErrCodeClientNotFound ErrorCode = "CLIENT_NOT_FOUND"
	ErrCodeClientCreate   ErrorCode = "CLIENT_CREATE_ERROR"
	ErrCodeClientDelete   ErrorCode = "CLIENT_DELETE_ERROR"
	ErrCodeClientUpdate   ErrorCode = "CLIENT_UPDATE_ERROR"

	// Shoot operations

	ErrCodeShootNotFound ErrorCode = "SHOOT_NOT_FOUND"
	ErrCodeShootCreate   ErrorCode = "SHOOT_CREATE_ERROR"
	ErrCodeShootDelete   ErrorCode = "SHOOT_DELETE_ERROR"
	ErrCodeShootUpdate   ErrorCode = "SHOOT_UPDATE_ERROR"

	// Validation

	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT_ERROR"

	// External Services
	ErrCodeCurrencyAPIRequest  ErrorCode = "CURRENCY_API_REQUEST_ERROR"
	ErrCodeCurrencyAPIResponse ErrorCode = "CURRENCY_API_RESPONSE_ERROR"
	ErrCodeCurrencyAPIParsing  ErrorCode = "CURRENCY_API_PARSING_ERROR"
	ErrCodeCurrencyNotFound    ErrorCode = "CURRENCY_NOT_FOUND"

	ErrCodeKafkaProduce ErrorCode = "KAFKA_SEND_MSG_ERROR"

	//Configuration

	ErrCodeConfig ErrorCode = "CONFIG_ERROR"

	// Jobs
	ErrCodeJobError ErrorCode = "JOB_ERROR"
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

func IsErrorCode(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
