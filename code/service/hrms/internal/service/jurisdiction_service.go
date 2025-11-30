// In jurisdiction_service.go
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"hrms/internal/clients/boundary"
	"hrms/internal/models"
	"hrms/internal/repository"
	"hrms/pkg/errors"
)

type jurisdictionService struct {
	repo           repository.JurisdictionRepository
	employeeSvc    EmployeeService
	boundaryClient *boundary.Client
}

// NewJurisdictionService creates a new jurisdiction service
func NewJurisdictionService(repo repository.JurisdictionRepository, employeeSvc EmployeeService, boundaryClient *boundary.Client) JurisdictionService {
	return &jurisdictionService{
		repo:           repo,
		employeeSvc:    employeeSvc,
		boundaryClient: boundaryClient,
	}
}

func (s *jurisdictionService) CreateJurisdiction(ctx context.Context, req *models.CreateJurisdictionRequest, tenantID string) (*models.JurisdictionResponse, error) {
	// Validate boundary codes
	// if err := s.validateBoundaryCodes(ctx, tenantID, req.BoundaryRelation); err != nil {
	// 	return nil, err
	// }

	// Create jurisdiction model
	now := time.Now()
	lastModTime := now.Unix()
	jurisdiction := &models.Jurisdiction{
		ID:               uuid.New().String(),
		EmployeeID:       req.EmployeeID,
		BoundaryRelation: req.BoundaryRelation,
		IsActive:         true, // Default to true if not provided
		TenantID:         tenantID,
		CreatedTime:      now.Unix(),
		LastModifiedTime: &lastModTime,
	}

	// Override with provided IsActive if it's not nil
	if req.IsActive != nil {
		jurisdiction.IsActive = *req.IsActive
	}

	// Save to database
	if err := s.repo.Create(ctx, jurisdiction); err != nil {
		return nil, fmt.Errorf("failed to create jurisdiction: %w", err)
	}

	return toJurisdictionResponse(jurisdiction), nil
}

func (s *jurisdictionService) validateBoundaryCodes(ctx context.Context, tenantID string, codes []string) error {
	if len(codes) == 0 {
		return nil
	}

	boundaries, err := s.boundaryClient.SearchByCodes(ctx, tenantID, codes)
	if err != nil {
		return fmt.Errorf("failed to validate boundary codes: %w", err)
	}

	// Create a map of valid codes for quick lookup
	validCodes := make(map[string]bool)
	for _, b := range boundaries {
		validCodes[b.Code] = true
	}

	// Check for invalid codes
	var invalidCodes []string
	for _, code := range codes {
		if !validCodes[code] {
			invalidCodes = append(invalidCodes, code)
		}
	}

	if len(invalidCodes) > 0 {
		return fmt.Errorf("invalid boundary codes for tenant %s: %s",
			tenantID, strings.Join(invalidCodes, ", "))
	}

	return nil
}

func (s *jurisdictionService) GetJurisdictionByUUID(ctx context.Context, uuid, tenantID string) (*models.JurisdictionResponse, error) {
	jurisdiction, err := s.repo.FindByUUID(ctx, uuid, tenantID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.WithDescription("jurisdiction not found").WithOperation("GetJurisdictionByUUID")
		}
		logrus.WithError(err).Error("Failed to get jurisdiction")
		return nil, errors.Wrap(err, "DATABASE_ERROR", "failed to get jurisdiction").WithOperation("GetJurisdictionByUUID")
	}

	return toJurisdictionResponse(jurisdiction), nil
}

// GetJurisdictionsByEmployeeID retrieves all jurisdictions for a specific employee
func (s *jurisdictionService) GetJurisdictionsByEmployeeID(ctx context.Context, employeeID, tenantID string) ([]*models.JurisdictionResponse, error) {
	// First, verify the employee exists
	_, err := s.employeeSvc.GetEmployeeByUUID(ctx, employeeID, tenantID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.New("NOT_FOUND", "employee not found").WithOperation("GetJurisdictionsByEmployeeID")
		}
		logrus.WithError(err).Error("Failed to validate employee")
		return nil, errors.Wrap(err, "VALIDATION_ERROR", "failed to validate employee").WithOperation("GetJurisdictionsByEmployeeID")
	}

	// Create search criteria
	criteria := &models.JurisdictionSearchCriteria{
		EmployeeIDs: []string{employeeID},
		TenantID:    tenantID,
	}

	// Search for jurisdictions
	jurisdictions, err := s.repo.Search(ctx, criteria)
	if err != nil {
		logrus.WithError(err).Error("Failed to get jurisdictions by employee ID")
		return nil, errors.Wrap(err, "DATABASE_ERROR", "failed to get jurisdictions by employee ID").WithOperation("GetJurisdictionsByEmployeeID")
	}

	// Convert to response objects
	responses := make([]*models.JurisdictionResponse, 0, len(jurisdictions))
	for _, j := range jurisdictions {
		responses = append(responses, toJurisdictionResponse(j))
	}

	return responses, nil
}

func (s *jurisdictionService) SearchJurisdictions(ctx context.Context, criteria *models.JurisdictionSearchCriteria) ([]*models.JurisdictionResponse, error) {
	jurisdictions, err := s.repo.Search(ctx, criteria)
	if err != nil {
		logrus.WithError(err).Error("Failed to search jurisdictions")
		return nil, errors.Wrap(err, "DATABASE_ERROR", "failed to search jurisdictions").WithOperation("SearchJurisdictions")
	}

	responses := make([]*models.JurisdictionResponse, 0, len(jurisdictions))
	for _, j := range jurisdictions {
		responses = append(responses, toJurisdictionResponse(j))
	}

	return responses, nil
}

func (s *jurisdictionService) UpdateJurisdiction(ctx context.Context, uuid string, req *models.UpdateJurisdictionRequest, tenantID string) (*models.JurisdictionResponse, error) {
	// Get existing jurisdiction
	existing, err := s.repo.FindByUUID(ctx, uuid, tenantID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.WithDescription("jurisdiction not found").WithOperation("UpdateJurisdiction")
		}
		logrus.WithError(err).Error("Failed to find jurisdiction")
		return nil, errors.Wrap(err, "DATABASE_ERROR", "failed to find jurisdiction").WithOperation("UpdateJurisdiction")
	}

	if req.EmployeeID != "" {
		existing.EmployeeID = req.EmployeeID
	}

	if req.BoundaryRelation != nil && len(*req.BoundaryRelation) > 0 {
		existing.BoundaryRelation = *req.BoundaryRelation
	}

	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	// existing.LastModifiedTime = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, existing); err != nil {
		logrus.WithError(err).Error("Failed to update jurisdiction")
		return nil, errors.Wrap(err, "DATABASE_ERROR", "failed to update jurisdiction").WithOperation("UpdateJurisdiction")
	}

	return toJurisdictionResponse(existing), nil
}

// toJurisdictionResponse converts a Jurisdiction model to JurisdictionResponse
func toJurisdictionResponse(j *models.Jurisdiction) *models.JurisdictionResponse {
	if j == nil {
		return nil
	}

	return &models.JurisdictionResponse{
		ID:               j.ID,
		EmployeeID:       j.EmployeeID,
		BoundaryRelation: j.BoundaryRelation,
		IsActive:         j.IsActive,
		TenantID:         j.TenantID,
		CreatedTime:      j.CreatedTime,
		LastModifiedTime: j.LastModifiedTime,
	}
}

func (s *jurisdictionService) DeleteJurisdiction(ctx context.Context, uuid, tenantID string) error {
	// Check if jurisdiction exists
	_, err := s.repo.FindByUUID(ctx, uuid, tenantID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return errors.ErrNotFound.WithDescription("jurisdiction not found").WithOperation("DeleteJurisdiction")
		}
		logrus.WithError(err).Error("Failed to find jurisdiction")
		return errors.Wrap(err, "DATABASE_ERROR", "failed to find jurisdiction").WithOperation("DeleteJurisdiction")
	}

	// Delete jurisdiction
	if err := s.repo.Delete(ctx, uuid, tenantID); err != nil {
		logrus.WithError(err).Error("Failed to delete jurisdiction")
		return errors.Wrap(err, "DATABASE_ERROR", "failed to delete jurisdiction").WithOperation("DeleteJurisdiction")
	}

	return nil
}

// ReplaceJurisdiction completely replaces an existing jurisdiction with new data
func (s *jurisdictionService) ReplaceJurisdiction(ctx context.Context, uuid string, req *models.UpdateJurisdictionRequest, tenantID string) (*models.JurisdictionResponse, error) {
	// Get existing jurisdiction
	existing, err := s.repo.FindByUUID(ctx, uuid, tenantID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.ErrNotFound.WithDescription("jurisdiction not found").WithOperation("ReplaceJurisdiction")
		}
		logrus.WithError(err).Error("Failed to find jurisdiction")
		return nil, errors.Wrap(err, "DATABASE_ERROR", "failed to find jurisdiction").WithOperation("ReplaceJurisdiction")
	}

	// Validate employee exists if EmployeeID is being updated
	if req.EmployeeID != "" && req.EmployeeID != existing.EmployeeID {
		_, err := s.employeeSvc.GetEmployeeByUUID(ctx, req.EmployeeID, tenantID)
		if err != nil {
			if errors.Is(err, errors.ErrNotFound) {
				return nil, errors.New("NOT_FOUND", "employee not found").WithOperation("ReplaceJurisdiction")
			}
			logrus.WithError(err).Error("Failed to validate employee")
			return nil, errors.Wrap(err, "VALIDATION_ERROR", "failed to validate employee").WithOperation("ReplaceJurisdiction")
		}
	}

	// Update all fields
	if req.EmployeeID != "" {
		existing.EmployeeID = req.EmployeeID
	}

	if req.BoundaryRelation != nil && len(*req.BoundaryRelation) > 0 {
		existing.BoundaryRelation = *req.BoundaryRelation
	}

	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	// existing.LastModifiedTime = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, existing); err != nil {
		logrus.WithError(err).Error("Failed to replace jurisdiction")
		return nil, errors.Wrap(err, "DATABASE_ERROR", "failed to replace jurisdiction").WithOperation("ReplaceJurisdiction")
	}

	return toJurisdictionResponse(existing), nil
}
