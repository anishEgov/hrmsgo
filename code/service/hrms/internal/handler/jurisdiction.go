package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"hrms/internal/models"
	"hrms/internal/service"
	"hrms/pkg/errors"
)

type JurisdictionHandler struct {
	service service.JurisdictionService
	logger  *logrus.Logger
}

func NewJurisdictionHandler(service service.JurisdictionService, logger *logrus.Logger) *JurisdictionHandler {
	if logger == nil {
		logger = logrus.New()
	}

	return &JurisdictionHandler{
		service: service,
		logger:  logger,
	}
}

type createJurisdictionRequest struct {
	Jurisdiction *models.Jurisdiction `json:"jurisdiction" binding:"required"`
}

func (h *JurisdictionHandler) SearchJurisdictions(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	criteria := &models.JurisdictionSearchCriteria{
		TenantID: tID,
	}

	// Parse query parameters
	if employeeID := c.Query("employeeId"); employeeID != "" {
		criteria.EmployeeIDs = []string{employeeID}
	}
	// boundaryRelation
	if relations := c.QueryArray("boundaryRelation"); len(relations) > 0 {
		for _, r := range relations {
			criteria.BoundaryRelations = append(criteria.BoundaryRelations, strings.Split(r, ",")...)
		}
	}
	if isActive := c.Query("isActive"); isActive != "" {
		active := isActive == "true"
		criteria.IsActive = &active
	}

	// Set pagination
	if limit, err := parseIntParam(c, "limit", 10); err == nil {
		criteria.Limit = limit
	}
	if offset, err := parseIntParam(c, "offset", 0); err == nil {
		criteria.Offset = offset
	}

	// Search jurisdictions
	jurisdictions, err := h.service.SearchJurisdictions(c.Request.Context(), criteria)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jurisdictions": jurisdictions})
}

func (h *JurisdictionHandler) GetJurisdictionByUUID(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	uuidStr := c.Param("uuid")
	if _, err := uuid.Parse(uuidStr); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid jurisdiction UUID"))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	jurisdiction, err := h.service.GetJurisdictionByUUID(c.Request.Context(), uuidStr, tID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			h.handleError(c, http.StatusNotFound, errors.New("NOT_FOUND", "Jurisdiction not found"))
			return
		}
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jurisdiction": jurisdiction})
}

func (h *JurisdictionHandler) CreateJurisdiction(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	var req models.CreateJurisdictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_REQUEST", err.Error()))
		return
	}

	jurisdiction, err := h.service.CreateJurisdiction(c.Request.Context(), &req, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"jurisdiction": jurisdiction})
}

func (h *JurisdictionHandler) ReplaceJurisdiction(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	uuidStr := c.Param("uuid")
	if _, err := uuid.Parse(uuidStr); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid jurisdiction UUID"))
		return
	}

	var req models.UpdateJurisdictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_REQUEST", err.Error()))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	jurisdiction, err := h.service.ReplaceJurisdiction(c.Request.Context(), uuidStr, &req, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jurisdiction": jurisdiction})
}

func (h *JurisdictionHandler) handleError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    getErrorCode(err),
			"message": getErrorMessage(err),
		},
	})
}

func getErrorCode(err error) string {
	if e, ok := err.(interface{ Code() string }); ok {
		return e.Code()
	}
	return "INTERNAL_ERROR"
}

func getErrorMessage(err error) string {
	if e, ok := err.(interface{ Error() string }); ok {
		return e.Error()
	}
	return "An unexpected error occurred"
}

// parseIntParam is a helper function to parse integer query parameters
func parseIntParam(c *gin.Context, param string, defaultValue int) (int, error) {
	value := c.DefaultQuery(param, "")
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}
