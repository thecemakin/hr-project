-- Create asset_assignments table
CREATE TABLE asset_assignments (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    asset_id INTEGER NOT NULL,
    employee_id INTEGER NOT NULL,
    assigned_by INTEGER NOT NULL,
    assigned_date TIMESTAMP WITH TIME ZONE NOT NULL,
    returned_date TIMESTAMP WITH TIME ZONE,
    returned_by INTEGER,
    notes TEXT,
    status VARCHAR(20) DEFAULT 'assigned' NOT NULL
);

-- Indexes
CREATE INDEX idx_asset_assignments_asset_id ON asset_assignments(asset_id);
CREATE INDEX idx_asset_assignments_employee_id ON asset_assignments(employee_id);
CREATE INDEX idx_asset_assignments_assigned_by ON asset_assignments(assigned_by);
CREATE INDEX idx_asset_assignments_returned_by ON asset_assignments(returned_by);
CREATE INDEX idx_asset_assignments_status ON asset_assignments(status);
CREATE INDEX idx_asset_assignments_assigned_date ON asset_assignments(assigned_date);
CREATE INDEX idx_asset_assignments_deleted_at ON asset_assignments(deleted_at);

-- Foreign key constraints
ALTER TABLE asset_assignments ADD CONSTRAINT fk_asset_assignments_asset FOREIGN KEY (asset_id) REFERENCES assets(id);
ALTER TABLE asset_assignments ADD CONSTRAINT fk_asset_assignments_employee FOREIGN KEY (employee_id) REFERENCES employees(id);
ALTER TABLE asset_assignments ADD CONSTRAINT fk_asset_assignments_assigned_by FOREIGN KEY (assigned_by) REFERENCES employees(id);
ALTER TABLE asset_assignments ADD CONSTRAINT fk_asset_assignments_returned_by FOREIGN KEY (returned_by) REFERENCES employees(id);