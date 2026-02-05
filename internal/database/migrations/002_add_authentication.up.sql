-- Add password_hash column to employees table for authentication
ALTER TABLE employees ADD COLUMN password_hash VARCHAR(255);

-- Note: For existing employees, password_hash will be NULL
-- They will need to set a password through registration or admin action
