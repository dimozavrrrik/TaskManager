-- Drop triggers
DROP TRIGGER IF EXISTS update_time_entries_updated_at ON time_entries;
DROP TRIGGER IF EXISTS update_task_messages_updated_at ON task_messages;
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop views
DROP VIEW IF EXISTS employee_tasks_view;
DROP VIEW IF EXISTS task_time_summary;

-- Drop tables
DROP TABLE IF EXISTS time_entries;
DROP TABLE IF EXISTS task_messages;
DROP TABLE IF EXISTS task_participants;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS employees;

-- Drop types
DROP TYPE IF EXISTS participant_role;
DROP TYPE IF EXISTS task_status;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
