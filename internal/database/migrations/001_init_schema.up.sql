-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Employees table
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    department VARCHAR(100) NOT NULL,
    position VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_department ON employees(department);
CREATE INDEX idx_employees_deleted_at ON employees(deleted_at);

-- Tasks table
CREATE TYPE task_status AS ENUM (
    'new',
    'in_progress',
    'code_review',
    'testing',
    'returned_with_errors',
    'closed'
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'new',
    priority INTEGER DEFAULT 0,
    created_by UUID NOT NULL REFERENCES employees(id),
    archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    due_date TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_tasks_archived ON tasks(archived);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);

-- Task participants table
CREATE TYPE participant_role AS ENUM (
    'executor',
    'responsible',
    'customer'
);

CREATE TABLE task_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    role participant_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(task_id, employee_id, role)
);

CREATE INDEX idx_task_participants_task ON task_participants(task_id);
CREATE INDEX idx_task_participants_employee ON task_participants(employee_id);
CREATE INDEX idx_task_participants_role ON task_participants(role);
CREATE INDEX idx_task_participants_composite ON task_participants(employee_id, task_id);

-- Task messages table
CREATE TABLE task_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    author_id UUID REFERENCES employees(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    is_system_message BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_task_messages_task ON task_messages(task_id);
CREATE INDEX idx_task_messages_author ON task_messages(author_id);
CREATE INDEX idx_task_messages_created ON task_messages(created_at DESC);
CREATE INDEX idx_task_messages_system ON task_messages(is_system_message);

-- Time entries table
CREATE TABLE time_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    hours DECIMAL(10, 2) NOT NULL CHECK (hours > 0),
    description TEXT,
    entry_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_time_entries_task ON time_entries(task_id);
CREATE INDEX idx_time_entries_employee ON time_entries(employee_id);
CREATE INDEX idx_time_entries_date ON time_entries(entry_date);
CREATE INDEX idx_time_entries_deleted ON time_entries(deleted_at);

-- View for aggregated time per task
CREATE OR REPLACE VIEW task_time_summary AS
SELECT
    task_id,
    SUM(hours) AS total_hours,
    COUNT(*) AS entry_count,
    COUNT(DISTINCT employee_id) AS unique_employees
FROM time_entries
WHERE deleted_at IS NULL
GROUP BY task_id;

-- View for main screen (user's tasks)
CREATE OR REPLACE VIEW employee_tasks_view AS
SELECT
    t.id AS task_id,
    t.title,
    t.status,
    tp.role AS user_role,
    tp.employee_id,
    t.priority,
    t.due_date,
    t.created_at,
    t.updated_at,
    COALESCE(tts.total_hours, 0) AS total_hours
FROM tasks t
INNER JOIN task_participants tp ON t.id = tp.task_id
LEFT JOIN task_time_summary tts ON t.id = tts.task_id
WHERE t.deleted_at IS NULL AND t.archived = FALSE;

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_employees_updated_at BEFORE UPDATE ON employees
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_task_messages_updated_at BEFORE UPDATE ON task_messages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_time_entries_updated_at BEFORE UPDATE ON time_entries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
