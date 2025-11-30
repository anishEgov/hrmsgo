-- Updated schema to align with OpenAPI specification

-- Drop existing tables if they exist
DROP TABLE IF EXISTS eg_hrms_jurisdiction_v3 CASCADE;
DROP TABLE IF EXISTS eg_hrms_employee_v3 CASCADE;

-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create employee table
CREATE TABLE IF NOT EXISTS eg_hrms_employee_v3 (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(64),
    user_id VARCHAR(64),
    individual_id VARCHAR(64),
    status VARCHAR(20),
    employee_type VARCHAR(50) NOT NULL,
    date_of_appointment TIMESTAMP WITH TIME ZONE,
    department VARCHAR(64) NOT NULL,
    designation VARCHAR(64) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    tenant_id VARCHAR(64) NOT NULL,
    created_by VARCHAR(64) NOT NULL,
    last_modified_by VARCHAR(64),
    created_time BIGINT NOT NULL,
    last_modified_time BIGINT,
    CONSTRAINT uk_employee_code_tenant UNIQUE (code, tenant_id)
);

-- Create jurisdiction table
CREATE TABLE IF NOT EXISTS eg_hrms_jurisdiction_v3 (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES eg_hrms_employee_v3(id) ON DELETE CASCADE,
    boundary_relation JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active BOOLEAN DEFAULT TRUE,
    tenant_id VARCHAR(64) NOT NULL,
    created_by VARCHAR(64) NOT NULL,
    last_modified_by VARCHAR(64),
    created_time BIGINT NOT NULL,
    last_modified_time BIGINT,
    CONSTRAINT fk_jurisdiction_employee FOREIGN KEY (employee_id) REFERENCES eg_hrms_employee_v3(id)
);

