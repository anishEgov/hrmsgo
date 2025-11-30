package errors

import (
	"fmt"
	"runtime"
)

// Error represents a structured error with stack trace
type Error struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Params      interface{} `json:"params,omitempty"`
	Op          string      `json:"-"` // Operation that caused the error
	Err         error       `json:"-"` // Underlying error
	Stack       string      `json:"-"` // Stack trace
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	if e.Description != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Description)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the Unwrap method for error unwrapping
func (e *Error) Unwrap() error { return e.Err }

// New creates a new error
func New(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Stack:   getStack(),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
		Stack:   getStack(),
	}
}

// WithDescription adds a description to the error
func (e *Error) WithDescription(desc string) *Error {
	e.Description = desc
	return e
}

// WithParams adds parameters to the error
func (e *Error) WithParams(params interface{}) *Error {
	e.Params = params
	return e
}

// WithOperation adds operation context to the error
func (e *Error) WithOperation(op string) *Error {
	e.Op = op
	return e
}

// Common error codes
var (
	// Generic errors
	ErrInternalServer = New("INTERNAL_ERROR", "An internal server error occurred")
	ErrInvalidInput   = New("INVALID_INPUT", "Invalid input provided")
	ErrNotFound       = New("NOT_FOUND", "The requested resource was not found")
	ErrUnauthorized   = New("UNAUTHORIZED", "You are not authorized to perform this action")
	ErrForbidden      = New("FORBIDDEN", "You don't have permission to access this resource")

	// Validation errors
	ErrValidationFailed = New("VALIDATION_ERROR", "Validation failed")

	// Employee errors
	ErrEmployeeNotFound    = New("EMPLOYEE_NOT_FOUND", "Employee not found")
	ErrEmployeeExists      = New("EMPLOYEE_EXISTS", "Employee already exists")
	ErrEmployeeDeactivated = New("EMPLOYEE_DEACTIVATED", "Employee account is deactivated")

	// Jurisdiction errors
	ErrJurisdictionNotFound = New("JURISDICTION_NOT_FOUND", "Jurisdiction not found")
	ErrJurisdictionExists   = New("JURISDICTION_EXISTS", "Jurisdiction already exists")

	// Database errors
	ErrDatabase = New("DATABASE_ERROR", "A database error occurred")
)

// ErrorResponse represents the JSON response for errors
type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Params      interface{} `json:"params,omitempty"`
}

// ToResponse converts an error to an ErrorResponse
func ToResponse(err error) ErrorResponse {
	switch e := err.(type) {
	case *Error:
		return ErrorResponse{
			Error: ErrorDetails{
				Code:        e.Code,
				Message:     e.Message,
				Description: e.Description,
				Params:      e.Params,
			},
		}
	default:
		return ErrorResponse{
			Error: ErrorDetails{
				Code:    ErrInternalServer.Code,
				Message: e.Error(),
			},
		}
	}
}

// Is checks if the error is of a specific error code
func Is(err error, target *Error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == target.Code
}

// getStack returns the current stack trace as a string
func getStack() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
