-- Clean up existing data to ensure predictable IDs for seeding
DELETE FROM leave_balances;
DELETE FROM leave_requests;
DELETE FROM leave_types;
DELETE FROM asset_assignments;
DELETE FROM assets;
DELETE FROM employees;
DELETE FROM positions;
DELETE FROM departments;

-- Reset sequences
ALTER SEQUENCE departments_id_seq RESTART WITH 1;
ALTER SEQUENCE positions_id_seq RESTART WITH 1;
ALTER SEQUENCE employees_id_seq RESTART WITH 1;
ALTER SEQUENCE assets_id_seq RESTART WITH 1;
ALTER SEQUENCE asset_assignments_id_seq RESTART WITH 1;
ALTER SEQUENCE leave_types_id_seq RESTART WITH 1;
ALTER SEQUENCE leave_balances_id_seq RESTART WITH 1;

-- 1. Departments
INSERT INTO departments (name, description) VALUES
('Management', 'Executive management of the company'), -- ID 1
('Engineering', 'Software development and technical operations'), -- ID 2
('Human Resources', 'People management, recruitment, and culture'), -- ID 3
('Marketing', 'Brand awareness and digital marketing campaigns'), -- ID 4
('Sales', 'Direct sales and customer acquisitions'); -- ID 5

-- 2. Positions
INSERT INTO positions (title, description, level) VALUES
('CEO', 'Chief Executive Officer', 0),                             -- ID 1
('Engineering Manager', 'Oversees engineering teams and projects', 1), -- ID 2
('HR Manager', 'Manages HR operations and policy', 1),             -- ID 3
('Senior Software Engineer', 'Advanced software development and mentoring', 2), -- ID 4
('Software Engineer', 'Full-stack software development', 3),        -- ID 5
('HR Specialist', 'Handles recruitment and employee relations', 2), -- ID 6
('Marketing Specialist', 'Digital marketing and SEO expert', 2),    -- ID 7
('Sales Executive', 'Direct B2B and B2C sales', 2);                 -- ID 8

-- 3. Employees
INSERT INTO employees (first_name, last_name, email, employee_number, status, department_id, position_id, manager_id) VALUES
('Robert', 'Biggs', 'ceo@hr-project.com', 'EMP001', 'active', 1, 1, NULL),               -- ID 1
('Ahmet', 'Yilmaz', 'ahmet.yilmaz@hr-project.com', 'EMP002', 'active', 2, 2, 1),      -- ID 2
('Sarah', 'Conner', 'sarah.conner@hr-project.com', 'EMP003', 'active', 3, 3, 1),       -- ID 3
('Zeynep', 'Demir', 'zeynep.demir@hr-project.com', 'EMP004', 'active', 2, 4, 2),       -- ID 4
('Arda', 'Bulut', 'arda.bulut@hr-project.com', 'EMP005', 'active', 2, 5, 2),          -- ID 5
('Michael', 'Scott', 'michael.scott@hr-project.com', 'EMP006', 'active', 5, 8, 1),     -- ID 6
('Pam', 'Beesly', 'pam.beesly@hr-project.com', 'EMP007', 'active', 3, 6, 3),          -- ID 7
('Jim', 'Halpert', 'jim.halpert@hr-project.com', 'EMP008', 'active', 5, 8, 6),         -- ID 8
('Dwight', 'Schrute', 'dwight.schrute@hr-project.com', 'EMP009', 'active', 5, 8, 6),    -- ID 9
('Angela', 'Martin', 'angela.martin@hr-project.com', 'EMP010', 'active', 3, 6, 3);      -- ID 10

-- 4. Assets
INSERT INTO assets (name, asset_tag, serial_number, type, model, brand, status) VALUES
('MacBook Pro 14', 'AST-LAP-001', 'SN-MBP14-001', 'Laptop', 'M3 Pro', 'Apple', 'available'), -- ID 1
('MacBook Pro 16', 'AST-LAP-002', 'SN-MBP16-002', 'Laptop', 'M3 Max', 'Apple', 'assigned'),  -- ID 2
('Dell Latitude 5520', 'AST-LAP-003', 'SN-DEL5520-003', 'Laptop', '5520', 'Dell', 'available'), -- ID 3
('Sony WH-1000XM4', 'AST-HED-001', 'SN-SONY-001', 'Headphones', 'XM4', 'Sony', 'assigned'), -- ID 4
('Dell 27 Monitor', 'AST-MON-001', 'SN-DELMON-001', 'Monitor', 'P2723D', 'Dell', 'available'); -- ID 5

-- 5. Asset Assignments
INSERT INTO asset_assignments (asset_id, employee_id, assigned_by, assigned_date, status) VALUES
(2, 2, 1, CURRENT_TIMESTAMP - INTERVAL '30 days', 'assigned'),
(4, 4, 1, CURRENT_TIMESTAMP - INTERVAL '10 days', 'assigned');

-- Update asset status to assigned
UPDATE assets SET status = 'assigned' WHERE id IN (2, 4);

-- 6. Leave Types
INSERT INTO leave_types (name, description, default_days, is_active) VALUES
('Annual Leave', 'Paid time off for holidays', 20, true), -- ID 1
('Sick Leave', 'Medical absence leave', 15, true),      -- ID 2
('Casual Leave', 'Short period leave for personal matters', 5, true),
('Parental Leave', 'Leave for new parents', 120, true),
('Bereavement Leave', 'Leave for family emergencies', 3, true);

-- 7. Leave Balances
INSERT INTO leave_balances (employee_id, leave_type_id, total_days, used_days) VALUES
(1, 1, 20, 0),
(1, 2, 15, 0),
(2, 1, 20, 2),
(2, 2, 15, 0),
(3, 1, 20, 0),
(4, 1, 20, 0);

-- 8. Users (matching some employees)
UPDATE users SET employee_id = 1 WHERE email = 'admin@hr-project.com';
