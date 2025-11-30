package repository

import (
	"context"

	"gorm.io/gorm"

	"hrms/internal/models"
	"hrms/pkg/errors"
)

// EmployeeRepository defines the interface for employee data access operations
type EmployeeRepository interface {
	// Create creates a new employee
	Create(ctx context.Context, employee *models.Employee) error

	// FindByUUID finds an employee by UUID
	FindByUUID(ctx context.Context, uuid, tenantID string) (*models.Employee, error)

	FindByCode(ctx context.Context, code, tenantID string) (*models.Employee, error)

	// Update updates an existing employee
	Update(ctx context.Context, employee *models.Employee) error

	// Delete deletes an employee by ID
	Delete(ctx context.Context, id, tenantID string) error

	// Search searches for employees based on criteria
	Search(ctx context.Context, criteria *models.EmployeeSearchCriteria) ([]*models.Employee, error)

	// UpdateStatus updates the status of an employee
	UpdateStatus(ctx context.Context, id, status, tenantID string) error

	// EmployeeCodeExists checks if an employee with the given code already exists
	EmployeeCodeExists(ctx context.Context, code, tenantID string) (bool, error)

	// UpdateIsActive updates the is_active status of an employee
	UpdateIsActive(ctx context.Context, id string, isActive bool, tenantID string) error
}

// employeeRepository implements the EmployeeRepository interface
type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (r *employeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	tx := r.db.WithContext(ctx).Table(models.Employee{}.TableName()).Create(employee)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to create employee")
	}
	return nil
}

func (r *employeeRepository) FindByCode(ctx context.Context, code, tenantID string) (*models.Employee, error) {
	var employee models.Employee
	tx := r.db.WithContext(ctx).Where("code = ? AND tenant_id = ?", code, tenantID).First(&employee)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithDescription("employee not found")
		}
		return nil, errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to find employee by ID")
	}
	return &employee, nil
}

func (r *employeeRepository) FindByUUID(ctx context.Context, uuid, tenantID string) (*models.Employee, error) {
	var employee models.Employee
	tx := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", uuid, tenantID).First(&employee)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithDescription("employee not found")
		}
		return nil, errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to find employee by ID")
	}
	return &employee, nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *models.Employee) error {
	tx := r.db.WithContext(ctx).Model(&models.Employee{}).Where("id = ?", employee.ID).Updates(employee)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to update employee")
	}
	if tx.RowsAffected == 0 {
		return errors.ErrNotFound.WithDescription("employee not found")
	}
	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id, tenantID string) error {
	tx := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Employee{})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to delete employee")
	}
	if tx.RowsAffected == 0 {
		return errors.ErrNotFound.WithDescription("employee not found")
	}
	return nil
}

func (r *employeeRepository) Search(ctx context.Context, criteria *models.EmployeeSearchCriteria) ([]*models.Employee, error) {
	var employees []*models.Employee

	tx := r.db.WithContext(ctx).Model(&models.Employee{}).Where("tenant_id = ?", criteria.TenantID)

	// Apply filters

	if len(criteria.UUIDs) > 0 {
		tx = tx.Where("id IN ?", criteria.UUIDs)
	}

	if len(criteria.Codes) > 0 {
		tx = tx.Where("code IN ?", criteria.Codes)
	}

	if len(criteria.Departments) > 0 {
		tx = tx.Where("department_id IN ?", criteria.Departments)
	}

	if len(criteria.Designations) > 0 {
		tx = tx.Where("designation_id IN ?", criteria.Designations)
	}

	if criteria.Phone != "" {
		tx = tx.Where("mobile_number = ?", criteria.Phone)
	}

	if criteria.IsActive != nil {
		tx = tx.Where("is_active = ?", *criteria.IsActive)
	}

	// Apply sorting
	orderBy := criteria.SortBy
	if orderBy == "" {
		orderBy = "created_time"
	}

	orderClause := orderBy
	if criteria.SortOrder == "desc" {
		orderClause += " DESC"
	} else {
		orderClause += " ASC"
	}
	tx = tx.Order(orderClause)

	// Apply pagination
	if criteria.Limit > 0 {
		tx = tx.Limit(criteria.Limit)
	}

	if criteria.Offset > 0 {
		tx = tx.Offset(criteria.Offset)
	}

	// Execute query
	tx = tx.Find(&employees)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to search employees")
	}

	return employees, nil
}

func (r *employeeRepository) UpdateStatus(ctx context.Context, id, status, tenantID string) error {
	err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("status", status).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound.WithDescription("employee not found")
		}
		return errors.Wrap(err, "DATABASE_ERROR", "failed to update employee status")
	}

	return nil
}

// EmployeeCodeExists checks if an employee with the given code already exists in the database
func (r *employeeRepository) EmployeeCodeExists(ctx context.Context, code, tenantID string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("code = ? AND tenant_id = ?", code, tenantID).
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, "DATABASE_ERROR", "failed to check if employee code exists")
	}

	return count > 0, nil
}

// UpdateIsActive updates the is_active status of an employee
func (r *employeeRepository) UpdateIsActive(ctx context.Context, id string, isActive bool, tenantID string) error {
	err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("is_active", isActive).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound.WithDescription("employee not found")
		}
		return errors.Wrap(err, "DATABASE_ERROR", "failed to update employee status")
	}
	return nil
}
