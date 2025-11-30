package models

import (
	"time"
)

// Jurisdiction represents a jurisdiction in the system
type Jurisdiction struct {
	ID               string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	EmployeeID       string     `json:"employeeId" gorm:"not null;index"`
	BoundaryRelation []string   `json:"boundaryRelation" db:"boundary_relation"`
	IsActive         bool       `json:"isActive" gorm:"default:true"`
	TenantID         string     `json:"tenantId" gorm:"not null;index"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

// CreateJurisdictionRequest represents the request payload for creating a jurisdiction
type CreateJurisdictionRequest struct {
	EmployeeID       string   `json:"employeeId" binding:"required"`
	BoundaryRelation []string `json:"boundaryRelation" validate:"required,min=1"`
	IsActive         *bool     `json:"isActive"`
}

// UpdateJurisdictionRequest represents the request payload for updating a jurisdiction
type UpdateJurisdictionRequest struct {
	EmployeeID       string    `json:"employeeId,omitempty"`
	BoundaryRelation *[]string `json:"boundaryRelation,omitempty" validate:"omitempty,min=1"`
	IsActive         *bool     `json:"isActive,omitempty"`
}

// JurisdictionResponse represents the response payload for jurisdiction operations
type JurisdictionResponse struct {
	ID               string    `json:"id"`
	EmployeeID       string    `json:"employeeId"`
	BoundaryRelation []string  `json:"boundaryRelation"`
	IsActive         bool      `json:"isActive"`
	TenantID         string    `json:"tenantId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// JurisdictionSearchCriteria represents the search criteria for jurisdictions
type JurisdictionSearchCriteria struct {
	IDs               []string `form:"ids"`
	EmployeeIDs       []string `form:"employeeIds"`
	BoundaryRelations []string `form:"boundaryRelation"`
	IsActive          *bool    `form:"isActive"`
	Limit             int      `form:"limit,default=10"`
	Offset            int      `form:"offset,default=0"`
	SortBy            string   `form:"sortBy,default=createdAt"`
	SortOrder         string   `form:"sortOrder,default=desc"`
	TenantID          string
}
