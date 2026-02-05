package handler

import (
	"net/http"
	"strconv"

	"github.com/dmitry/taskmanager/internal/dto"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/internal/service"
	"github.com/dmitry/taskmanager/pkg/validator"
)

type EmployeeHandler struct {
	service   *service.EmployeeService
	validator *validator.Validator
}

func NewEmployeeHandler(service *service.EmployeeService, validator *validator.Validator) *EmployeeHandler {
	return &EmployeeHandler{
		service:   service,
		validator: validator,
	}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	if err := h.validator.Validate(req); err != nil {
		RespondError(w, err)
		return
	}

	employee, err := h.service.CreateEmployee(r.Context(), req.Name, req.Department, req.Position, req.Email)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusCreated, dto.ToEmployeeResponse(employee))
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	employee, err := h.service.GetEmployee(r.Context(), id)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, dto.ToEmployeeResponse(employee))
}

func (h *EmployeeHandler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	department := r.URL.Query().Get("department")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	filter := repository.EmployeeFilter{
		Department: department,
		Page:       page,
		PageSize:   pageSize,
	}

	employees, total, err := h.service.GetAllEmployees(r.Context(), filter)
	if err != nil {
		RespondError(w, err)
		return
	}

	responses := make([]dto.EmployeeResponse, len(employees))
	for i, emp := range employees {
		responses[i] = dto.ToEmployeeResponse(emp)
	}

	totalPages := (total + pageSize - 1) / pageSize

	RespondJSON(w, http.StatusOK, dto.PaginatedResponse{
		Data:       responses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	var req dto.UpdateEmployeeRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	if err := h.validator.Validate(req); err != nil {
		RespondError(w, err)
		return
	}

	employee, err := h.service.GetEmployee(r.Context(), id)
	if err != nil {
		RespondError(w, err)
		return
	}

	employee.Name = req.Name
	employee.Department = req.Department
	employee.Position = req.Position
	employee.Email = req.Email

	if err := h.service.UpdateEmployee(r.Context(), employee); err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, dto.ToEmployeeResponse(employee))
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteEmployee(r.Context(), id); err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "Сотрудник успешно удалён"})
}
