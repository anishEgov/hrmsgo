package service

import (
	"context"
	"hrms/internal/models"
)

// JurisdictionService defines the interface for jurisdiction operations

type JurisdictionService interface {
	CreateJurisdiction(ctx context.Context, req *models.CreateJurisdictionRequest, tenantID string) (*models.JurisdictionResponse, error)
	UpdateJurisdiction(ctx context.Context, id string, req *models.UpdateJurisdictionRequest, tenantID string) (*models.JurisdictionResponse, error)
	DeleteJurisdiction(ctx context.Context, id, tenantID string) error
	SearchJurisdictions(ctx context.Context, criteria *models.JurisdictionSearchCriteria) ([]*models.JurisdictionResponse, error)
	GetJurisdictionByUUID(ctx context.Context, uuid, tenantID string) (*models.JurisdictionResponse, error)
	ReplaceJurisdiction(ctx context.Context, uuid string, req *models.UpdateJurisdictionRequest, tenantID string) (*models.JurisdictionResponse, error)
}
