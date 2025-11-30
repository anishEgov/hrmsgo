package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"hrms/db"
	"hrms/internal/clients/boundary"
	"hrms/internal/clients/idgen"
	hrmsConfig "hrms/internal/config"
	"hrms/internal/handler"
	"hrms/internal/repository"
	"hrms/internal/router"
	hrmsService "hrms/internal/service"
)

func main() {
	// Initialize logger
	logger := initLogger()
	logger.Info("Starting HRMS Service...")

	// Load configuration
	cfg, err := hrmsConfig.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbConn, err := initDatabase(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	logger.Info("Database connection established")

	// Initialize repositories
	employeeRepo := repository.NewEmployeeRepository(dbConn)
	jurisdictionRepo := repository.NewJurisdictionRepository(dbConn)

	// Initialize ID generation client
	idGenClient := idgen.NewClient(idgen.Config{
		Host:      cfg.IDGen.Host,
		Path:      cfg.IDGen.Path,
		Enabled:   cfg.IDGen.Enabled,
		IDGenName: cfg.IDGen.IDGenName,
	})

	boundaryClient := boundary.NewClient(cfg.Boundary.BaseURL)

	// First, create employee service with a nil jurisdiction service
	employeeSvc := hrmsService.NewEmployeeService(employeeRepo, nil, idGenClient)

	// Then create jurisdiction service with the employee service
	jurisdictionSvc := hrmsService.NewJurisdictionService(
		jurisdictionRepo,
		employeeSvc,
		boundaryClient,
	)
	// Now update the employee service with the jurisdiction service
	employeeSvc = hrmsService.NewEmployeeService(employeeRepo, jurisdictionSvc, idGenClient)

	// Initialize handlers
	employeeHandler := handler.NewEmployeeHandler(employeeSvc, logger)
	jurisdictionHandler := handler.NewJurisdictionHandler(jurisdictionSvc, logger)

	// Setup router
	r := router.SetupRouter(cfg, employeeHandler, jurisdictionHandler, logger)

	// Start server in a goroutine
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Infof("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	sqlDB, err := dbConn.DB()
	if err == nil {
		sqlDB.Close()
	}

	logger.Info("Server exiting")
}

// initLogger initializes and configures the logger
func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Set log level based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warnf("Invalid log level '%s', defaulting to 'info'", logLevel)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	return logger
}

// initDatabase initializes the database connection
func initDatabase(cfg *hrmsConfig.Config, logger *logrus.Logger) (*gorm.DB, error) {
	dbConn, err := db.InitDB(cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable debug mode in development
	if os.Getenv("GIN_MODE") != "release" {
		dbConn = dbConn.Debug()
	}

	return dbConn, nil
}
