-- Create positions table
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    title VARCHAR(100) NOT NULL,
    description TEXT,
    code VARCHAR(20) UNIQUE,
    level INTEGER DEFAULT 1,
    salary_min BIGINT DEFAULT 0,
    salary_max BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    department_id INTEGER
);

-- Indexes
CREATE INDEX idx_positions_code ON positions(code);
CREATE INDEX idx_positions_department_id ON positions(department_id);
CREATE INDEX idx_positions_is_active ON positions(is_active);
CREATE INDEX idx_positions_deleted_at ON positions(deleted_at);