# eGov HRMS - Deployment and Testing Guide

## 🎯 Production-Ready Go Implementation

This is a complete production-ready migration of the eGov HRMS Java service to Go 1.21+, following DIGIT ecosystem principles.

---

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start](#quick-start)
3. [Configuration](#configuration)
4. [Database Setup](#database-setup)
5. [Running the Application](#running-the-application)
6. [API Testing](#api-testing)
7. [Docker Deployment](#docker-deployment)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software

- **Go**: 1.21 or higher
- **PostgreSQL**: 15 or higher
- **Apache Kafka**: 3.0+ (for event-driven features)
- **Docker & Docker Compose**: Latest versions (for containerized deployment)

### Installation

```bash
# Install Go 1.21+
wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Verify installation
go version
```

---

## Quick Start

### Option 1: Using Make (Recommended)

```bash
# Clone and navigate to project
cd /home/admin/Desktop/egov-hrms

# Download dependencies
make deps

# Run database migrations
make migrate

# Build and run the application
make run
```

### Option 2: Using Docker Compose (Complete Stack)

```bash
# Start all services (PostgreSQL, Kafka, Zookeeper, HRMS)
make docker-run

# Or manually
docker-compose up -d
```

### Option 3: Manual Build

```bash
# Download dependencies
go mod download

# Build the binary
go build -o egov-hrms ./cmd/server

# Run
./egov-hrms
```

---

## Configuration

### Environment Variables

The application can be configured using environment variables or a `config.yaml` file.

#### Database Configuration

```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=postgres
export DATABASE_PASSWORD=postgres
export DATABASE_DBNAME=egov_hrms
export DATABASE_SSL_MODE=disable
export DATABASE_MAX_OPEN_CONNS=25
export DATABASE_MAX_IDLE_CONNS=5
```

#### Kafka Configuration

```bash
export KAFKA_BOOTSTRAP_SERVERS=localhost:9092
export KAFKA_CONSUMER_GROUP=egov-hrms-consumer-group
export KAFKA_TOPIC_SAVE_EMPLOYEE=save-hrms-employee
export KAFKA_TOPIC_UPDATE_EMPLOYEE=update-hrms-employee
export KAFKA_TOPIC_NOTIFICATION_SMS=egov.core.notification.sms
```

#### Server Configuration

```bash
export SERVER_PORT=8080
export SERVER_READ_TIMEOUT=30s
export SERVER_WRITE_TIMEOUT=30s
export LOG_LEVEL=info
```

#### HRMS Specific Configuration

```bash
export HRMS_NOTIFICATION_ENABLED=true
export HRMS_EMPLOYEE_APP_LINK=https://mseva.lgpunjab.gov.in/employee/user/login
export HRMS_DEFAULT_PAGINATION_LIMIT=200
export HRMS_DEFAULT_PWD_LENGTH=10
export HRMS_IDGEN_NAME=hrms.employeecode
export HRMS_IDGEN_FORMAT="EMP-[city]-[SEQ_EG_HRMS_EMP_CODE]"
```

#### External Services (DIGIT Ecosystem)

```bash
export USER_SERVICE_HOST=http://localhost:8081
export IDGEN_SERVICE_HOST=http://localhost:8082
export MDMS_SERVICE_HOST=http://localhost:8083
export FILESTORE_SERVICE_HOST=http://localhost:8084
export LOCALIZATION_SERVICE_HOST=http://localhost:8085
export NOTIFICATION_SERVICE_HOST=http://localhost:8086
```

### Configuration File (config.yaml)

Alternatively, create a `config/config.yaml` file:

```yaml
server:
  port: 8080
  readTimeout: 30s
  writeTimeout: 30s

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: egov_hrms
  sslMode: disable
  maxOpenConns: 25
  maxIdleConns: 5

kafka:
  bootstrapServers: localhost:9092
  consumer:
    groupId: egov-hrms-consumer-group
  topics:
    saveEmployee: save-hrms-employee
    hrmsUpdate: update-hrms-employee
    notificationSms: egov.core.notification.sms

hrms:
  notificationEnabled: true
  employeeAppLink: https://mseva.lgpunjab.gov.in/employee/user/login
  defaultPaginationLimit: 200
  defaultPwdLength: 10
  idgenName: hrms.employeecode
  idgenFormat: "EMP-[city]-[SEQ_EG_HRMS_EMP_CODE]"
```

---

## Database Setup

### Using Migration Script

```bash
# Make migration script executable
chmod +x db/migrate.sh

# Run migrations
./db/migrate.sh
```

### Manual Migration

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE egov_hrms;

# Connect to database
\c egov_hrms

# Run migration SQL
\i db/migrations/001_initial_schema.up.sql
```

### Migration Details

The migration creates the following tables:

1. **eg_hrms_employee** - Main employee table
2. **eg_hrms_jurisdiction** - Employee jurisdictions
3. **eg_hrms_assignment** - Employee assignments (designations, departments)
4. **eg_hrms_service_history** - Service history records
5. **eg_hrms_education** - Educational qualifications
6. **eg_hrms_departmental_test** - Departmental test records
7. **eg_hrms_document** - Employee documents
8. **eg_hrms_deactivation_details** - Deactivation records
9. **eg_hrms_reactivation_details** - Reactivation records

### Verify Tables

```bash
psql -U postgres -d egov_hrms -c "\dt"
```

---

## Running the Application

### Development Mode

```bash
# Run with live reload (using air or similar)
make run

# Or direct execution
go run cmd/server/main.go
```

### Production Mode

```bash
# Build optimized binary
make build

# Run binary
./server
```

### Health Check

```bash
# Check application health
curl http://localhost:8080/health

# Expected response:
# {"status":"UP"}
```

### Swagger Documentation

Access API documentation at: `http://localhost:8080/swagger/index.html`

---

## API Testing

### Create Employee

```bash
curl -X POST http://localhost:8080/egov-hrms/employees/_create \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {
      "apiId": "emp-services",
      "ver": "1.0",
      "ts": null,
      "action": "create",
      "did": "",
      "key": "",
      "msgId": "20170310130900",
      "requesterId": "",
      "authToken": "{{authToken}}",
      "userInfo": {
        "id": 1,
        "uuid": "11b0e02b-0145-4de2-bc42-c97b96264807",
        "userName": "admin",
        "name": "Admin",
        "mobileNumber": "9999999999",
        "emailId": "admin@example.com",
        "type": "EMPLOYEE",
        "roles": [
          {
            "name": "Employee Admin",
            "code": "EMPLOYEE_ADMIN",
            "tenantId": "pb.amritsar"
          }
        ],
        "tenantId": "pb.amritsar"
      }
    },
    "Employees": [
      {
        "tenantId": "pb.amritsar",
        "employeeStatus": "EMPLOYED",
        "employeeType": "PERMANENT",
        "dateOfAppointment": 1609459200000,
        "user": {
          "userName": "",
          "name": "John Doe",
          "gender": "MALE",
          "mobileNumber": "9876543210",
          "emailId": "john.doe@example.com",
          "correspondenceAddress": "123 Main Street, Amritsar",
          "type": "EMPLOYEE",
          "roles": [
            {
              "code": "EMPLOYEE",
              "name": "Employee",
              "tenantId": "pb.amritsar"
            }
          ],
          "tenantId": "pb.amritsar",
          "dob": 631152000000,
          "permanentAddress": "123 Main Street, Amritsar"
        },
        "jurisdictions": [
          {
            "hierarchy": "ADMIN",
            "boundaryType": "City",
            "boundary": "pb.amritsar"
          }
        ],
        "assignments": [
          {
            "fromDate": 1609459200000,
            "toDate": null,
            "department": "DEPT_1",
            "designation": "DESIG_1",
            "isHOD": false,
            "isCurrentAssignment": true
          }
        ]
      }
    ]
  }'
```

### Search Employee

```bash
# Search by tenant ID
curl -X POST http://localhost:8080/egov-hrms/employees/_search \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {
      "apiId": "emp-services",
      "ver": "1.0",
      "ts": null,
      "action": "search",
      "did": "",
      "key": "",
      "msgId": "20170310130900",
      "authToken": "{{authToken}}",
      "userInfo": {
        "id": 1,
        "uuid": "11b0e02b-0145-4de2-bc42-c97b96264807",
        "userName": "admin",
        "tenantId": "pb.amritsar"
      }
    },
    "tenantId": "pb.amritsar",
    "codes": [],
    "names": [],
    "roles": [],
    "departments": [],
    "designations": [],
    "isActive": true,
    "limit": 10,
    "offset": 0
  }'

# Search by employee code
curl -X POST "http://localhost:8080/egov-hrms/employees/_search" \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {...},
    "tenantId": "pb.amritsar",
    "codes": ["EMP-PB-AMRITSAR-001"]
  }'

# Search by name
curl -X POST "http://localhost:8080/egov-hrms/employees/_search" \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {...},
    "tenantId": "pb.amritsar",
    "names": ["John Doe"]
  }'

# Search by department
curl -X POST "http://localhost:8080/egov-hrms/employees/_search" \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {...},
    "tenantId": "pb.amritsar",
    "departments": ["DEPT_1"]
  }'
```

### Update Employee

```bash
curl -X POST http://localhost:8080/egov-hrms/employees/_update \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {
      "apiId": "emp-services",
      "ver": "1.0",
      "ts": null,
      "action": "update",
      "authToken": "{{authToken}}",
      "userInfo": {
        "id": 1,
        "uuid": "11b0e02b-0145-4de2-bc42-c97b96264807",
        "userName": "admin",
        "tenantId": "pb.amritsar"
      }
    },
    "Employees": [
      {
        "uuid": "{{employee_uuid}}",
        "tenantId": "pb.amritsar",
        "code": "EMP-PB-AMRITSAR-001",
        "employeeStatus": "EMPLOYED",
        "employeeType": "PERMANENT",
        "user": {
          "uuid": "{{user_uuid}}",
          "name": "John Doe Updated",
          "mobileNumber": "9876543210",
          "emailId": "john.updated@example.com"
        },
        "assignments": [
          {
            "uuid": "{{assignment_uuid}}",
            "fromDate": 1609459200000,
            "toDate": null,
            "department": "DEPT_2",
            "designation": "DESIG_2",
            "isCurrentAssignment": true
          }
        ]
      }
    ]
  }'
```

### Employee Count

```bash
curl -X POST http://localhost:8080/egov-hrms/employees/_count \
  -H "Content-Type: application/json" \
  -d '{
    "RequestInfo": {
      "apiId": "emp-services",
      "ver": "1.0",
      "authToken": "{{authToken}}",
      "userInfo": {
        "tenantId": "pb.amritsar"
      }
    },
    "tenantId": "pb.amritsar"
  }'

# Expected response:
# {
#   "ResponseInfo": {...},
#   "EmployeesCount": {
#     "active": 50,
#     "inactive": 5
#   }
# }
```

---

## Docker Deployment

### Build Docker Image

```bash
# Using Make
make docker-build

# Or manually
docker build -f pkg/Dockerfile -t egov-hrms:latest .
```

### Run with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f egov-hrms

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Docker Services

The `docker-compose.yml` includes:

1. **PostgreSQL** - Database (Port 5432)
2. **Zookeeper** - Kafka coordination (Port 2181)
3. **Kafka** - Message broker (Port 9092)
4. **egov-hrms** - Application service (Port 8080)

### Access Services

- **HRMS API**: http://localhost:8080
- **Swagger**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **PostgreSQL**: localhost:5432
- **Kafka**: localhost:9092

---

## Makefile Commands

```bash
# Install dependencies
make deps

# Run database migrations
make migrate

# Build the application
make build

# Run the application
make run

# Run tests
make test

# Generate swagger docs
make swagger

# Format code
make fmt

# Run linter
make lint

# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run

# Stop Docker containers
make docker-stop

# Clean build artifacts
make clean

# Show help
make help
```

---

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U postgres -h localhost -p 5432 -c "SELECT version();"

# Check database exists
psql -U postgres -c "\l" | grep egov_hrms
```

### Kafka Connection Issues

```bash
# Check Kafka is running
docker ps | grep kafka

# List Kafka topics
docker exec -it egov-hrms-kafka kafka-topics --list --bootstrap-server localhost:9092

# Create missing topics
docker exec -it egov-hrms-kafka kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic save-hrms-employee \
  --partitions 3 \
  --replication-factor 1
```

### Port Already in Use

```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill the process
sudo kill -9 <PID>

# Or change the port
export SERVER_PORT=8081
```

### Build Errors

```bash
# Clean and rebuild
make clean
go clean -cache
go mod tidy
make build
```

### View Application Logs

```bash
# If running with Docker
docker-compose logs -f egov-hrms

# If running locally
tail -f /var/log/egov-hrms.log
```

---

## Performance Tuning

### Database Connection Pooling

```bash
export DATABASE_MAX_OPEN_CONNS=50
export DATABASE_MAX_IDLE_CONNS=10
export DATABASE_CONN_MAX_LIFETIME=300s
```

### Kafka Consumer Tuning

```bash
export KAFKA_CONSUMER_SESSION_TIMEOUT=30000
export KAFKA_CONSUMER_MAX_POLL_RECORDS=500
export KAFKA_CONSUMER_FETCH_MIN_BYTES=1024
```

### Server Tuning

```bash
export SERVER_READ_TIMEOUT=60s
export SERVER_WRITE_TIMEOUT=60s
export SERVER_IDLE_TIMEOUT=120s
export SERVER_MAX_HEADER_BYTES=1048576
```

---

## Monitoring and Health Checks

### Health Endpoint

```bash
curl http://localhost:8080/health
```

### Metrics (if Prometheus is enabled)

```bash
curl http://localhost:8080/metrics
```

### Kubernetes Liveness/Readiness Probes

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

---

## Production Checklist

- [ ] All environment variables configured
- [ ] Database migrations executed successfully
- [ ] Kafka topics created
- [ ] External service endpoints configured (User, MDMS, IDGen, etc.)
- [ ] SSL/TLS certificates configured
- [ ] Logging configured to centralized system
- [ ] Monitoring and alerting set up
- [ ] Backup strategy in place
- [ ] Load balancer configured
- [ ] API rate limiting configured
- [ ] Security headers configured
- [ ] CORS policy configured
- [ ] Multi-tenancy tested

---

## Support and Documentation

- **API Documentation**: http://localhost:8080/swagger/index.html
- **GitHub Repository**: https://github.com/egovernments/egov-hrms
- **eGov DIGIT Documentation**: https://core.digit.org/
- **Swagger Contract**: https://editor.swagger.io/?url=https://raw.githubusercontent.com/egovernments/business-services/master/Docs/hrms-v1.0.0.yaml

---

## License

Apache 2.0 - See LICENSE file for details.

---

## Contributors

Migrated from Java to Go following eGov DIGIT principles and best practices.
