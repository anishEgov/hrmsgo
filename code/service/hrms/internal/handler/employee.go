package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"hrms/internal/models"
	"hrms/internal/service"
	"hrms/pkg/errors"
)

type EmployeeHandler struct {
	service service.EmployeeService
	logger  *logrus.Logger
}

func NewEmployeeHandler(service service.EmployeeService, logger *logrus.Logger) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
		logger:  logger,
	}
}

func (h *EmployeeHandler) handleError(c *gin.Context, status int, err error) {
	h.logger.WithError(err).Error("Request failed")
	c.JSON(status, gin.H{"error": err.Error()})
}

func (h *EmployeeHandler) CreateEmployees(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	var req []*models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_REQUEST", err.Error()))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	employees, err := h.service.CreateEmployees(c.Request.Context(), req, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, employees)
}

func (h *EmployeeHandler) SearchEmployees(c *gin.Context) {
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

	criteria := &models.EmployeeSearchCriteria{
		TenantID: tID,
		// Add other search criteria from query params as needed
	}

	employees, err := h.service.SearchEmployees(c.Request.Context(), criteria)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) GetEmployeeByUUID(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid employee UUID"))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	employee, err := h.service.GetEmployeeByUUID(c.Request.Context(), id, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid employee UUID"))
		return
	}

	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_REQUEST", err.Error()))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	employee, err := h.service.UpdateEmployee(c.Request.Context(), id, &req, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) HardDeleteEmployee(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid employee UUID"))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	if err := h.service.HardDeleteEmployee(c.Request.Context(), id, tID); err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EmployeeHandler) PatchEmployee(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid employee UUID"))
		return
	}

	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_REQUEST", err.Error()))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}
	employee, err := h.service.PatchEmployee(c.Request.Context(), id, &req, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) DeactivateEmployee(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid employee UUID"))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	employee, err := h.service.DeactivateEmployee(c.Request.Context(), id, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) ReactivateEmployee(c *gin.Context) {
	tenantID, exists := c.Get("tenantID")
	if !exists {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Tenant ID not found in context"))
		return
	}

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		h.handleError(c, http.StatusBadRequest, errors.New("INVALID_UUID", "Invalid employee UUID"))
		return
	}

	tID, ok := tenantID.(string)
	if !ok {
		h.handleError(c, http.StatusInternalServerError, errors.New("INTERNAL_ERROR", "Invalid tenant ID format"))
		return
	}

	employee, err := h.service.ReactivateEmployee(c.Request.Context(), id, tID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employee)
}
