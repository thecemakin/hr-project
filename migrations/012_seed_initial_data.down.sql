-- Reverse initial seed data
DELETE FROM leave_balances;
DELETE FROM leave_types;
DELETE FROM asset_assignments;
DELETE FROM assets;
UPDATE users SET employee_id = NULL WHERE email = 'admin@hr-project.com';
DELETE FROM employees;
DELETE FROM positions;
DELETE FROM departments;
