package validator

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"hrms/internal/config"
	"hrms/internal/models"
	"hrms/internal/repository"

	"github.com/go-playground/validator/v10"
)

var (
	// phoneRegex validates 10-digit phone numbers starting with 6-9
	phoneRegex = regexp.MustCompile(`^[6-9][0-9]{9}$`)
	// emailRegex validates standard email format
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// EmployeeValidator validates employee data against the OpenAPI specification
type EmployeeValidator struct {
	repo     repository.EmployeeRepository
	cfg      *config.Config
	validate *validator.Validate
}

// NewEmployeeValidator creates a new EmployeeValidator with custom validations
func NewEmployeeValidator(repo repository.EmployeeRepository, cfg *config.Config) *EmployeeValidator {
	v := validator.New()

	// Register custom validations
	_ = v.RegisterValidation("employeeType", validateEmployeeType)
	_ = v.RegisterValidation("phone", validatePhone)
	_ = v.RegisterValidation("employeeStatus", validateEmployeeStatus)
	_ = v.RegisterValidation("dateTime", validateDateTime)

	return &EmployeeValidator{
		repo:     repo,
		cfg:      cfg,
		validate: v,
	}
}

// ValidateCreate validates a new employee creation request
func (v *EmployeeValidator) ValidateCreate(ctx context.Context, emp *models.Employee) error {
	// Required fields validation
	if emp.TenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}

	if emp.EmployeeType == "" {
		return fmt.Errorf("employee type is required")
	}

	if !isValidEmployeeType(emp.EmployeeType) {
		return fmt.Errorf("invalid employee type: %s. Must be one of: PERMANENT, CONTRACT, TEMPORARY", emp.EmployeeType)
	}

	if emp.Department == "" {
		return fmt.Errorf("department is required")
	}

	if emp.Designation == "" {
		return fmt.Errorf("designation is required")
	}

	// Validate field formats and constraints
	if emp.Code != "" {
		if len(emp.Code) < 2 || len(emp.Code) > 64 {
			return fmt.Errorf("code must be between 2 and 64 characters")
		}
	}

	// Check for duplicate employee code
	if emp.Code != "" {
		existing, err := v.repo.FindByCode(ctx, emp.Code, emp.TenantID)
		if err == nil && existing != nil && existing.ID != emp.ID {
			return fmt.Errorf("employee with code %s already exists", emp.Code)
		}
	}

	return nil
}

// ValidateUpdate validates an employee update request
func (v *EmployeeValidator) ValidateUpdate(ctx context.Context, emp *models.Employee, existing *models.Employee) error {
	if existing == nil {
		return fmt.Errorf("employee not found")
	}
	return nil
}

// ValidatePatch validates a partial employee update
func (v *EmployeeValidator) ValidatePatch(patch *models.UpdateEmployeeRequest) error {
	if patch.EmployeeType != nil {
		if !isValidEmployeeType(*patch.EmployeeType) {
			return fmt.Errorf("invalid employee type: %s. Must be one of: PERMANENT, CONTRACT, TEMPORARY", *patch.EmployeeType)
		}
	}

	if patch.EmployeeStatus != nil {
		if !isValidEmployeeStatus(*patch.EmployeeStatus) {
			return fmt.Errorf("invalid employee status: %s. Must be one of: ACTIVE, INACTIVE, SUSPENDED", *patch.EmployeeStatus)
		}
	}

	if patch.Phone != nil && !phoneRegex.MatchString(*patch.Phone) {
		return fmt.Errorf("invalid mobile number format. Must be a 10-digit number starting with 6-9")
	}

	if patch.EmailId != nil && !emailRegex.MatchString(*patch.EmailId) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateSearch validates employee search criteria
func (v *EmployeeValidator) ValidateSearch(ctx context.Context, criteria *models.EmployeeSearchCriteria) error {
	if criteria.TenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}

	// Validate limit and offset
	if criteria.Limit < 1 {
		criteria.Limit = 10
	}

	if criteria.Offset < 0 {
		criteria.Offset = 0
	}

	// Validate sort order
	if criteria.SortBy != "" {
		validSortFields := map[string]bool{
			"code":        true,
			"createdAt":   true,
			"updatedAt":   true,
			"employeeType": true,
			"status":      true,
		}

		if !validSortFields[criteria.SortBy] {
			return fmt.Errorf("invalid sort field: %s", criteria.SortBy)
		}

		if criteria.SortOrder != "" && criteria.SortOrder != "asc" && criteria.SortOrder != "desc" {
			return fmt.Errorf("invalid sort order: %s. Must be 'asc' or 'desc'", criteria.SortOrder)
		}
	}

	return nil
}

// normalizeValidationErrors converts validation errors to a more user-friendly format
func (v *EmployeeValidator) normalizeValidationErrors(err error) error {
	if err == nil {
		return nil
	}

	// If it's already a formatted error, return as is
	if _, ok := err.(*validator.ValidationErrors); !ok {
		return err
	}

	// Convert validation errors to a more user-friendly format
	var errs validator.ValidationErrors
	errs = err.(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}

	// Just return the first error for simplicity
	switch errs[0].Tag() {
	case "required":
		return fmt.Errorf("%s is required", errs[0].Field())
	case "email":
		return fmt.Errorf("invalid email format")
	case "min":
		return fmt.Errorf("%s must be at least %s characters", errs[0].Field(), errs[0].Param())
	case "max":
		return fmt.Errorf("%s must not exceed %s characters", errs[0].Field(), errs[0].Param())
	default:
		return fmt.Errorf("invalid value for %s", errs[0].Field())
	}
}

// Custom validation functions
func validateEmployeeType(fl validator.FieldLevel) bool {
	return isValidEmployeeType(fl.Field().String())
}

func validatePhone(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true // Empty is valid, use required tag if field is mandatory
	}
	return phoneRegex.MatchString(fl.Field().String())
}

func validateEmployeeStatus(fl validator.FieldLevel) bool {
	return isValidEmployeeStatus(fl.Field().String())
}

func validateDateTime(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

// Helper functions
func isValidEmployeeType(empType string) bool {
	validTypes := map[string]bool{
		"PERMANENT": true,
		"CONTRACT":  true,
		"TEMPORARY": true,
	}
	return validTypes[empType]
}

func isValidEmployeeStatus(status string) bool {
	validStatuses := map[string]bool{
		"ACTIVE":   true,
		"INACTIVE": true,
		"SUSPENDED": true,
	}
	return validStatuses[status]
}
