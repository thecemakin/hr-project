-- Create audit_logs table
CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    actor_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL, -- CREATE, UPDATE, DELETE, APPROVE, REJECT, LOGIN
    resource_type VARCHAR(50) NOT NULL, -- employees, assets, leave_requests, etc.
    resource_id INTEGER,
    
    old_values JSONB, -- The state before the change
    new_values JSONB, -- The state after the change
    
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB -- Any additional context (e.g., "Reason for rejection")
);

-- Indexes for performance
CREATE INDEX idx_audit_logs_actor_id ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
