package models

import "fmt"

// AuditDetails holds audit information for database records
type AuditDetails struct {
	CreatedBy    string `json:"createdBy,omitempty" gorm:"column:createdby;type:varchar(250)"`
	LastModified string `json:"lastModified,omitempty" gorm:"column:lastmodified;type:varchar(64)"`
	CreatedTime  int64  `json:"createdTime,omitempty" gorm:"column:createdtime;type:bigint"`
	UpdatedTime  int64  `json:"updatedTime,omitempty" gorm:"column:lastmodifiedtime;type:bigint"`
}

// ErrorResponse represents an error response in the API
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error represents a single error in the API response
// Error implements the error interface for the Error type
func (e *Error) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Description)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

type Error struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Params      interface{} `json:"params,omitempty"`
}
