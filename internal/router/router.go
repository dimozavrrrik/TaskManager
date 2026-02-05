package router

import (
	"net/http"

	"github.com/dmitry/taskmanager/internal/handler"
	"github.com/dmitry/taskmanager/internal/middleware"
	"github.com/dmitry/taskmanager/internal/service"
	"github.com/dmitry/taskmanager/pkg/logger"
	"github.com/gorilla/mux"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	employeeHandler *handler.EmployeeHandler,
	taskHandler *handler.TaskHandler,
	jwtService *service.JWTService,
	frontendURL string,
	logger *logger.Logger,
) http.Handler {
	r := mux.NewRouter()

	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.CORSMiddleware(frontendURL))

	api := r.PathPrefix("/api/v1").Subrouter()

	// Публичные маршруты (аутентификация не требуется)
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"taskmanager"}`))
	}).Methods("GET")

	// Маршруты аутентификации (аутентификация не требуется)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	auth.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Защищенные маршруты (требуется JWT аутентификация)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtService))

	// Эндпоинты для работы с сотрудниками
	protected.HandleFunc("/employees", employeeHandler.CreateEmployee).Methods("POST")
	protected.HandleFunc("/employees", employeeHandler.GetAllEmployees).Methods("GET")
	protected.HandleFunc("/employees/{id}", employeeHandler.GetEmployee).Methods("GET")
	protected.HandleFunc("/employees/{id}", employeeHandler.UpdateEmployee).Methods("PUT")
	protected.HandleFunc("/employees/{id}", employeeHandler.DeleteEmployee).Methods("DELETE")
	protected.HandleFunc("/employees/{id}/tasks", taskHandler.GetEmployeeTasks).Methods("GET")

	// Эндпоинты для работы с задачами
	protected.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	protected.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	protected.HandleFunc("/tasks/{id}/status", taskHandler.UpdateTaskStatus).Methods("PATCH")
	protected.HandleFunc("/tasks/{id}/archive", taskHandler.ArchiveTask).Methods("PATCH")
	protected.HandleFunc("/tasks/{id}/participants", taskHandler.GetTaskParticipants).Methods("GET")
	protected.HandleFunc("/tasks/{id}/participants", taskHandler.AddParticipant).Methods("POST")

	return r
}
