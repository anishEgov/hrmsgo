package models

import (
	"time"
)

// Employee represents an employee in the system
type Employee struct {
	ID                string        `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Code              string        `json:"code,omitempty" gorm:"index"`
	UserID            string        `json:"userId,omitempty"`
	IndividualID      string        `json:"individualId,omitempty"`
	Status            string        `json:"status,omitempty"`
	EmployeeType      string        `json:"employeeType,omitempty" gorm:"not null"`
	DateOfAppointment *time.Time    `json:"dateOfAppointment,omitempty"`
	Department        string        `json:"department,omitempty" gorm:"not null"`
	Designation       string        `json:"designation,omitempty" gorm:"not null"`
	IsActive          bool          `json:"isActive,omitempty" gorm:"default:true"`
	Jurisdictions     []*Jurisdiction `json:"jurisdictions,omitempty" gorm:"foreignKey:EmployeeID"`
	TenantID          string        `json:"-"`
	CreatedBy         string        `json:"-"`
	LastModifiedBy    *string       `json:"-"`
	CreatedTime       int64         `json:"-"`
	LastModifiedTime  *int64        `json:"-"`
}

// CreateEmployeeRequest represents the request payload for creating an employee
type CreateEmployeeRequest struct {
	Code              string        `json:"code,omitempty"`
	UserID            string        `json:"userId,omitempty"`
	IndividualID      string        `json:"individualId,omitempty"`
	Status            string        `json:"status,omitempty"`
	EmployeeType      string        `json:"employeeType,omitempty"`
	DateOfAppointment *time.Time    `json:"dateOfAppointment,omitempty"`
	Department        string        `json:"department,omitempty"`
	Designation       string        `json:"designation,omitempty"`
	IsActive          *bool         `json:"isActive,omitempty"`
	Jurisdictions     []*Jurisdiction `json:"jurisdictions,omitempty"`
}

// UpdateEmployeeRequest represents the request payload for updating an employee
type UpdateEmployeeRequest struct {
	EmployeeStatus *string `json:"employeeStatus,omitempty"`
	EmployeeType   *string `json:"employeeType,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	EmailId        *string `json:"emailId,omitempty"`
	IsActive       *bool   `json:"isActive,omitempty"`
}

// EmployeeResponse represents the response payload for employee operations
type EmployeeResponse struct {
	ID                string                `json:"id"`
	Code              string                `json:"code,omitempty"`
	UserID            string                `json:"userId,omitempty"`
	IndividualID      string                `json:"individualId,omitempty"`
	Status            string                `json:"status,omitempty"`
	EmployeeType      string                `json:"employeeType,omitempty"`
	DateOfAppointment *time.Time            `json:"dateOfAppointment,omitempty"`
	Department        string                `json:"department,omitempty"`
	Designation       string                `json:"designation,omitempty"`
	IsActive          bool                  `json:"isActive"`
	Jurisdictions     []*JurisdictionResponse `json:"jurisdictions,omitempty"`
}

// EmployeeSearchCriteria represents the search criteria for employees
type EmployeeSearchCriteria struct {
	UUIDs        []string `form:"uuids"`
	Codes        []string `form:"codes"`
	Departments  []string `form:"departments"`
	Designations []string `form:"designations"`
	Phone        string   `form:"phone"`
	IsActive     *bool    `form:"isActive"`
	Limit        int      `form:"limit,default=10"`
	Offset       int      `form:"offset,default=0"`
	SortBy       string   `form:"sortBy,default=createdAt"`
	SortOrder    string   `form:"sortOrder,default=desc"`
	TenantID     string
}

// TableName specifies the table name for the Employee model
func (Employee) TableName() string {
	return "eg_hrms_employee_v3"
}
