package handler

import (
	"net/http"
	"strconv"

	"github.com/dmitry/taskmanager/internal/domain"
	"github.com/dmitry/taskmanager/internal/dto"
	"github.com/dmitry/taskmanager/internal/middleware"
	"github.com/dmitry/taskmanager/internal/repository"
	"github.com/dmitry/taskmanager/internal/service"
	"github.com/dmitry/taskmanager/pkg/errors"
	"github.com/dmitry/taskmanager/pkg/validator"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service   *service.TaskService
	validator *validator.Validator
}

func NewTaskHandler(service *service.TaskService, validator *validator.Validator) *TaskHandler {
	return &TaskHandler{
		service:   service,
		validator: validator,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	if err := h.validator.Validate(req); err != nil {
		RespondError(w, err)
		return
	}

	createdByID, err := middleware.GetEmployeeIDFromContext(r.Context())
	if err != nil {
		RespondError(w, err)
		return
	}

	participants := make([]service.ParticipantInput, len(req.Participants))
	for i, p := range req.Participants {
		empID, err := uuid.Parse(p.EmployeeID)
		if err != nil {
			RespondError(w, errors.BadRequest("Неверный ID сотрудника"))
			return
		}
		participants[i] = service.ParticipantInput{
			EmployeeID: empID,
			Role:       domain.ParticipantRole(p.Role),
		}
	}

	task, err := h.service.CreateTask(r.Context(), service.CreateTaskRequest{
		Title:        req.Title,
		Description:  req.Description,
		Priority:     req.Priority,
		CreatedBy:    createdByID,
		Participants: participants,
	})

	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusCreated, dto.ToTaskResponse(task))
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, dto.ToTaskResponse(task))
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	filter := repository.TaskFilter{
		Page:     page,
		PageSize: pageSize,
	}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		filter.Status = []domain.TaskStatus{domain.TaskStatus(statusStr)}
	}

	tasks, total, err := h.service.GetAllTasks(r.Context(), filter)
	if err != nil {
		RespondError(w, err)
		return
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = dto.ToTaskResponse(task)
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

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	var req dto.UpdateTaskStatusRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	if err := h.validator.Validate(req); err != nil {
		RespondError(w, err)
		return
	}

	err := h.service.UpdateTaskStatus(r.Context(), id, domain.TaskStatus(req.Status))
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "Статус задачи успешно обновлён"})
}

func (h *TaskHandler) ArchiveTask(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	err := h.service.ArchiveTask(r.Context(), id)
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "Задача успешно архивирована"})
}

func (h *TaskHandler) GetTaskParticipants(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	participants, err := h.service.GetParticipants(r.Context(), id)
	if err != nil {
		RespondError(w, err)
		return
	}

	responses := make([]dto.TaskParticipantResponse, len(participants))
	for i, p := range participants {
		responses[i] = dto.ToTaskParticipantResponse(p)
	}

	RespondJSON(w, http.StatusOK, responses)
}

func (h *TaskHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	taskID, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	var req dto.AddParticipantRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	if err := h.validator.Validate(req); err != nil {
		RespondError(w, err)
		return
	}

	employeeID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		RespondError(w, errors.BadRequest("Неверный ID сотрудника"))
		return
	}

	err = h.service.AddParticipant(r.Context(), taskID, employeeID, domain.ParticipantRole(req.Role))
	if err != nil {
		RespondError(w, err)
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]string{"message": "Участник успешно добавлен"})
}

func (h *TaskHandler) GetEmployeeTasks(w http.ResponseWriter, r *http.Request) {
	employeeID, ok := ParseUUID(w, r, "id")
	if !ok {
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	filter := repository.TaskFilter{
		Page:     page,
		PageSize: pageSize,
	}

	tasks, total, err := h.service.GetTasksForEmployee(r.Context(), employeeID, filter)
	if err != nil {
		RespondError(w, err)
		return
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = dto.ToTaskResponse(task)
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
