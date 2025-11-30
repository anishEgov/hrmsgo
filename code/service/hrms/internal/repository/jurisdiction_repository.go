// internal/repository/jurisdiction_repository.go
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"hrms/internal/models"
	"hrms/pkg/errors"
)

type JurisdictionRepository interface {
	Create(ctx context.Context, jurisdiction *models.Jurisdiction) error
	FindByID(ctx context.Context, id, tenantID string) (*models.Jurisdiction, error)
	FindByUUID(ctx context.Context, uuid, tenantID string) (*models.Jurisdiction, error)
	Update(ctx context.Context, jurisdiction *models.Jurisdiction) error
	Delete(ctx context.Context, id, tenantID string) error
	Search(ctx context.Context, criteria *models.JurisdictionSearchCriteria) ([]*models.Jurisdiction, error)
}

type jurisdictionRepository struct {
	db *gorm.DB
}

func NewJurisdictionRepository(db *gorm.DB) JurisdictionRepository {
	return &jurisdictionRepository{
		db: db,
	}
}

func (r *jurisdictionRepository) Create(ctx context.Context, jurisdiction *models.Jurisdiction) error {
	// Set timestamps
	now := time.Now()
	jurisdiction.CreatedAt = now
	jurisdiction.UpdatedAt = now

	// Create the record using GORM
	tx := r.db.WithContext(ctx).Create(jurisdiction)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to create jurisdiction")
	}
	return nil
}

func (r *jurisdictionRepository) FindByID(ctx context.Context, id, tenantID string) (*models.Jurisdiction, error) {
	var jurisdiction models.Jurisdiction
	tx := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&jurisdiction)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithDescription("jurisdiction not found")
		}
		return nil, errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to find jurisdiction")
	}

	return &jurisdiction, nil
}

func (r *jurisdictionRepository) FindByUUID(ctx context.Context, uuid, tenantID string) (*models.Jurisdiction, error) {
	var jurisdiction models.Jurisdiction
	tx := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", uuid, tenantID).
		First(&jurisdiction)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound.WithDescription("jurisdiction not found")
		}
		return nil, errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to find jurisdiction")
	}

	return &jurisdiction, nil
}

func (r *jurisdictionRepository) Update(ctx context.Context, jurisdiction *models.Jurisdiction) error {
	jurisdiction.UpdatedAt = time.Now()
	tx := r.db.WithContext(ctx).Save(jurisdiction)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to update jurisdiction")
	}
	return nil
}

func (r *jurisdictionRepository) Delete(ctx context.Context, id, tenantID string) error {
	tx := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.Jurisdiction{})
	
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to delete jurisdiction")
	}
	
	if tx.RowsAffected == 0 {
		return errors.ErrNotFound.WithDescription("jurisdiction not found")
	}
	
	return nil
}

func (r *jurisdictionRepository) Search(ctx context.Context, criteria *models.JurisdictionSearchCriteria) ([]*models.Jurisdiction, error) {
	var jurisdictions []*models.Jurisdiction

	tx := r.db.WithContext(ctx).Model(&models.Jurisdiction{})

	if criteria.TenantID != "" {
		tx = tx.Where("tenant_id = ?", criteria.TenantID)
	}

	if len(criteria.IDs) > 0 {
		tx = tx.Where("id IN ?", criteria.IDs)
	}

	if len(criteria.EmployeeIDs) > 0 {
		tx = tx.Where("employee_id IN ?", criteria.EmployeeIDs)
	}

	if len(criteria.BoundaryRelations) > 0 {
		// Using GORM's JSONB contains operator
		for _, code := range criteria.BoundaryRelations {
			tx = tx.Where("boundary_relation @> ?", json.RawMessage(fmt.Sprintf(`["%s"]`, code)))
		}
	}

	if criteria.IsActive != nil {
		tx = tx.Where("is_active = ?", *criteria.IsActive)
	}

	// Apply pagination
	if criteria.Limit > 0 {
		tx = tx.Limit(criteria.Limit)
	}
	if criteria.Offset > 0 {
		tx = tx.Offset(criteria.Offset)
	}

	// Apply sorting
	if criteria.SortBy != "" {
		order := criteria.SortBy
		if criteria.SortOrder != "" {
			order = order + " " + criteria.SortOrder
		}
		tx = tx.Order(order)
	}

	tx = tx.Find(&jurisdictions)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to search jurisdictions")
	}

	return jurisdictions, nil
}

func (r *jurisdictionRepository) FindByEmployeeID(ctx context.Context, employeeID, tenantID string) ([]*models.Jurisdiction, error) {
	var jurisdictions []*models.Jurisdiction
	tx := r.db.WithContext(ctx).
		Where("employee_id = ? AND tenant_id = ?", employeeID, tenantID).
		Find(&jurisdictions)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "DATABASE_ERROR", "failed to find jurisdictions by employee ID")
	}

	return jurisdictions, nil
}
