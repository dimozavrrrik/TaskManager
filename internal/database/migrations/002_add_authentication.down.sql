-- Remove password_hash column from employees table
ALTER TABLE employees DROP COLUMN IF EXISTS password_hash;
