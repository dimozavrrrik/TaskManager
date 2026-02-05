package service

import (
	"context"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/google/uuid"
)

type EmployeeService struct {
	repo   repository.EmployeeRepository
	logger *logger.Logger
}

func NewEmployeeService(repo repository.EmployeeRepository, logger *logger.Logger) *EmployeeService {
	return &EmployeeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, name, department, position, email string) (*domain.Employee, error) {
	employee := domain.NewEmployee(name, department, position, email)

	if err := s.repo.Create(ctx, employee); err != nil {
		return nil, err
	}

	s.logger.Info("Сотрудник создан", "employee_id", employee.ID, "email", email)

	return employee, nil
}

func (s *EmployeeService) GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EmployeeService) GetAllEmployees(ctx context.Context, filter repository.EmployeeFilter) ([]*domain.Employee, int, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s *EmployeeService) UpdateEmployee(ctx context.Context, employee *domain.Employee) error {
	return s.repo.Update(ctx, employee)
}

func (s *EmployeeService) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.Info("Сотрудник удалён", "employee_id", id)

	return nil
}
