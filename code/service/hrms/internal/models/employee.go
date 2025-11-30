package models

import (
	"time"
)

// Employee represents an employee in the system
type Employee struct {
	ID                string        `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Code              string        `json:"code" gorm:"not null;index"`
	Name              string        `json:"name" gorm:"not null"`
	DateOfBirth       *time.Time    `json:"dateOfBirth,omitempty"`
	Gender            string        `json:"gender,omitempty"`
	MobileNumber      string        `json:"mobileNumber,omitempty"`
	Email             string        `json:"email,omitempty"`
	PanNumber         string        `json:"panNumber,omitempty"`
	AadhaarNumber     string        `json:"aadhaarNumber,omitempty"`
	DateOfAppointment *time.Time    `json:"dateOfAppointment,omitempty"`
	DateOfRetirement  *time.Time    `json:"dateOfRetirement,omitempty"`
	DepartmentID      string        `json:"departmentId,omitempty"`
	DesignationID     string        `json:"designationId,omitempty"`
	EmployeeType      string        `json:"employeeType,omitempty"`
	Status            string        `json:"status" gorm:"default:'ACTIVE'"`
	IsActive          bool          `json:"isActive" gorm:"default:true"`
	TenantID          string        `json:"tenantId" gorm:"not null;index"`
	Jurisdictions     []*Jurisdiction `json:"jurisdictions,omitempty" gorm:"foreignKey:EmployeeID"`
	CreatedBy         string        `json:"-" gorm:"not null"`
	LastModifiedBy    *string       `json:"-"`
	CreatedTime       int64         `json:"createdAt" gorm:"column:created_time;not null"`
	LastModifiedTime  *int64        `json:"updatedAt,omitempty" gorm:"column:last_modified_time"`
	DeletedAt         *time.Time    `json:"-" gorm:"index"`
}

// CreateEmployeeRequest represents the request payload for creating an employee
type CreateEmployeeRequest struct {
	Code              string        `json:"code" binding:"required"`
	Name              string        `json:"name" binding:"required"`
	DateOfBirth       *time.Time    `json:"dateOfBirth,omitempty"`
	Gender            string        `json:"gender,omitempty"`
	MobileNumber      string        `json:"mobileNumber,omitempty"`
	Email             string        `json:"email,omitempty"`
	PanNumber         string        `json:"panNumber,omitempty"`
	AadhaarNumber     string        `json:"aadhaarNumber,omitempty"`
	DateOfAppointment *time.Time    `json:"dateOfAppointment,omitempty"`
	DateOfRetirement  *time.Time    `json:"dateOfRetirement,omitempty"`
	DepartmentID      string        `json:"departmentId,omitempty"`
	DesignationID     string        `json:"designationId,omitempty"`
	EmployeeType      string        `json:"employeeType,omitempty"`
	Status            string        `json:"status,omitempty"`
	IsActive          *bool         `json:"isActive,omitempty"`
	Jurisdictions     []*Jurisdiction `json:"jurisdictions,omitempty"`
}

// UpdateEmployeeRequest represents the request payload for updating an employee
type UpdateEmployeeRequest struct {
	Code              *string    `json:"code,omitempty"`
	Name              *string    `json:"name,omitempty"`
	DateOfBirth       *time.Time `json:"dateOfBirth,omitempty"`
	Gender            *string    `json:"gender,omitempty"`
	MobileNumber      *string    `json:"mobileNumber,omitempty"`
	Email             *string    `json:"email,omitempty"`
	PanNumber         *string    `json:"panNumber,omitempty"`
	AadhaarNumber     *string    `json:"aadhaarNumber,omitempty"`
	DateOfAppointment *time.Time `json:"dateOfAppointment,omitempty"`
	DateOfRetirement  *time.Time `json:"dateOfRetirement,omitempty"`
	DepartmentID      *string    `json:"departmentId,omitempty"`
	DesignationID     *string    `json:"designationId,omitempty"`
	EmployeeType      *string    `json:"employeeType,omitempty"`
	Status            *string    `json:"status,omitempty"`
	IsActive          *bool      `json:"isActive,omitempty"`
}

// EmployeeResponse represents the response payload for employee operations
type EmployeeResponse struct {
	ID                string                `json:"id"`
	Code              string                `json:"code"`
	Name              string                `json:"name"`
	DateOfBirth       *time.Time            `json:"dateOfBirth,omitempty"`
	Gender            string                `json:"gender,omitempty"`
	MobileNumber      string                `json:"mobileNumber,omitempty"`
	Email             string                `json:"email,omitempty"`
	PanNumber         string                `json:"panNumber,omitempty"`
	AadhaarNumber     string                `json:"aadhaarNumber,omitempty"`
	DateOfAppointment *time.Time            `json:"dateOfAppointment,omitempty"`
	DateOfRetirement  *time.Time            `json:"dateOfRetirement,omitempty"`
	DepartmentID      string                `json:"departmentId,omitempty"`
	DesignationID     string                `json:"designationId,omitempty"`
	EmployeeType      string                `json:"employeeType,omitempty"`
	Status            string                `json:"status"`
	IsActive          bool                  `json:"isActive"`
	TenantID          string                `json:"tenantId"`
	Jurisdictions     []*JurisdictionResponse `json:"jurisdictions,omitempty"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
}

// EmployeeSearchCriteria represents the search criteria for employees
type EmployeeSearchCriteria struct {
	IDs            []string `form:"ids"`
	UUIDs          []string `form:"uuids"`
	Codes          []string `form:"codes"`
	Status         []string `form:"status"`
	EmployeeTypes  []string `form:"employeeTypes"`
	Departments    []string `form:"departments"`
	Designations   []string `form:"designations"`
	Limit          int      `form:"limit,default=10"`
	Offset         int      `form:"offset,default=0"`
	SortBy         string   `form:"sortBy,default=createdAt"`
	SortOrder      string   `form:"sortOrder,default=desc"`
	TenantID       string
}

// TableName specifies the table name for the Employee model
func (Employee) TableName() string {
	return "eg_hrms_employee_v3"
}
