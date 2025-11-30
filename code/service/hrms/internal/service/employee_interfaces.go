// internal/service/employee_interfaces.go
package service

import (
	"context"
	"hrms/internal/models"
)

// EmployeeService defines the interface for employee operations
type EmployeeService interface {
	// CreateEmployees creates one or more employees
	CreateEmployees(ctx context.Context, req []*models.CreateEmployeeRequest, tenantID string) ([]*models.EmployeeResponse, error)

	// SearchEmployees searches for employees based on criteria
	SearchEmployees(ctx context.Context, criteria *models.EmployeeSearchCriteria) ([]*models.EmployeeResponse, error)

	// GetEmployeeByUUID retrieves an employee by UUID
	GetEmployeeByUUID(ctx context.Context, uuid, tenantID string) (*models.EmployeeResponse, error)

	// UpdateEmployee updates an employee by UUID
	UpdateEmployee(ctx context.Context, uuid string, req *models.CreateEmployeeRequest, tenantID string) (*models.EmployeeResponse, error)

	// HardDeleteEmployee permanently deletes an employee and all related records
	HardDeleteEmployee(ctx context.Context, uuid, tenantID string) error

	// PatchEmployee partially updates an employee
	PatchEmployee(ctx context.Context, uuid string, req *models.UpdateEmployeeRequest, tenantID string) (*models.EmployeeResponse, error)

	// DeactivateEmployee deactivates an employee
	DeactivateEmployee(ctx context.Context, uuid string, req *models.DeactivationDetails, tenantID string) (*models.EmployeeResponse, error)

	// ReactivateEmployee reactivates an inactive employee
	ReactivateEmployee(ctx context.Context, uuid string, req *models.ReactivationDetails, tenantID string) (*models.EmployeeResponse, error)
}
