package models

import "time"

// EmployeePatch represents a partial update for an employee
// All fields are pointers to allow for partial updates (PATCH)
type EmployeePatch struct {
	ID                string     `json:"id,omitempty"`
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
}

// DeactivationDetails contains information about employee deactivation
type DeactivationDetails struct {
	ID            string     `json:"id"`
	EmployeeID    string     `json:"employeeId" binding:"required"`
	EffectiveFrom *time.Time `json:"effectiveFrom" binding:"required"`
	Reason        string     `json:"reason" binding:"required"`
	Remarks       string     `json:"remarks,omitempty"`
	DeactivatedBy string     `json:"deactivatedBy" binding:"required"`
	TenantID      string     `json:"tenantId" binding:"required"`
}
