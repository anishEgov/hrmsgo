package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"hrms/internal/config"
	"hrms/internal/handler"
	"hrms/internal/middleware"
)

// SetupRouter initializes and configures the HTTP router
func SetupRouter(
	cfg *config.Config,
	employeeHandler *handler.EmployeeHandler,
	jurisdictionHandler *handler.JurisdictionHandler,
	logger *logrus.Logger,
) *gin.Engine {
	// Create router with default middleware
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Headers(logger))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// API v3 routes
	v3 := r.Group(cfg.Server.ContextPath + "/employees/v3")
	{
		// Employee endpoints
		v3.POST("", employeeHandler.CreateEmployees)
		v3.GET("", employeeHandler.SearchEmployees)

		// Employee by ID endpoints
		employeeID := v3.Group("/:id")
		{
			employeeID.GET("", employeeHandler.GetEmployeeByUUID)
			employeeID.PUT("", employeeHandler.UpdateEmployee)
			employeeID.DELETE("", employeeHandler.HardDeleteEmployee)
			employeeID.PATCH("", employeeHandler.PatchEmployee)

			// Employee status management
			employeeID.POST("deactivate", employeeHandler.DeactivateEmployee)
			employeeID.POST("reactivate", employeeHandler.ReactivateEmployee)
		}

		// Jurisdiction endpoints
		jurisdiction := v3.Group("/jurisdictions")
		{
			jurisdiction.POST("", jurisdictionHandler.CreateJurisdiction)
			jurisdiction.GET("", jurisdictionHandler.SearchJurisdictions)

			// Jurisdiction by UUID endpoints
			jurisdictionUUID := jurisdiction.Group("/:uuid")
			{
				jurisdictionUUID.GET("", jurisdictionHandler.GetJurisdictionByUUID)
				jurisdictionUUID.PUT("", jurisdictionHandler.ReplaceJurisdiction)
			}
		}
	}

	return r
}
