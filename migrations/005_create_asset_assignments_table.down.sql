-- Drop asset_assignments table and its foreign key constraints
ALTER TABLE asset_assignments DROP CONSTRAINT IF EXISTS fk_asset_assignments_asset;
ALTER TABLE asset_assignments DROP CONSTRAINT IF EXISTS fk_asset_assignments_employee;
ALTER TABLE asset_assignments DROP CONSTRAINT IF EXISTS fk_asset_assignments_assigned_by;
ALTER TABLE asset_assignments DROP CONSTRAINT IF EXISTS fk_asset_assignments_returned_by;

DROP TABLE IF EXISTS asset_assignments;