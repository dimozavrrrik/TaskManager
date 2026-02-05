using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace TaskManager.Client.Models;

public class TaskResponse
{
    [JsonPropertyName("id")]
    public Guid Id { get; set; }

    [JsonPropertyName("title")]
    public string Title { get; set; } = string.Empty;

    [JsonPropertyName("description")]
    public string Description { get; set; } = string.Empty;

    [JsonPropertyName("status")]
    public string Status { get; set; } = string.Empty;

    [JsonPropertyName("priority")]
    public int Priority { get; set; } // 0=низкий, 1=средний, 2=высокий

    [JsonPropertyName("created_by")]
    public Guid CreatedBy { get; set; }

    [JsonPropertyName("archived")]
    public bool Archived { get; set; }

    [JsonPropertyName("due_date")]
    public string? DueDate { get; set; }

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }

    [JsonPropertyName("updated_at")]
    public DateTime UpdatedAt { get; set; }

    // Вспомогательные свойства для отображения
    public string PriorityDisplay => Priority switch
    {
        0 => "Низкий",
        1 => "Средний",
        2 => "Высокий",
        _ => "Неизвестно"
    };
}

public class CreateTaskRequest
{
    [JsonPropertyName("title")]
    [Required(ErrorMessage = "Название обязательно")]
    [StringLength(500, MinimumLength = 3, ErrorMessage = "Название должно быть от 3 до 500 символов")]
    public string Title { get; set; } = string.Empty;

    [JsonPropertyName("description")]
    [Required(ErrorMessage = "Описание обязательно")]
    public string Description { get; set; } = string.Empty;

    [JsonPropertyName("priority")]
    [Required(ErrorMessage = "Приоритет обязателен")]
    [Range(0, 2, ErrorMessage = "Приоритет должен быть от 0 до 2")]
    public int Priority { get; set; } = 1; // 0=низкий, 1=средний, 2=высокий

    [JsonPropertyName("due_date")]
    public string? DueDate { get; set; }

    [JsonPropertyName("participants")]
    public List<ParticipantInput> Participants { get; set; } = new();
}

public class ParticipantInput
{
    [JsonPropertyName("employee_id")]
    [Required(ErrorMessage = "ID сотрудника обязателен")]
    public string EmployeeId { get; set; } = string.Empty;

    [JsonPropertyName("role")]
    [Required(ErrorMessage = "Роль обязательна")]
    public string Role { get; set; } = "executor"; // executor=исполнитель, responsible=ответственный, customer=заказчик
}

public class UpdateTaskStatusRequest
{
    [Required(ErrorMessage = "Статус обязателен")]
    public string Status { get; set; } = string.Empty;
}

public class TaskParticipantResponse
{
    [JsonPropertyName("task_id")]
    public Guid TaskId { get; set; }

    [JsonPropertyName("employee_id")]
    public Guid EmployeeId { get; set; }

    [JsonPropertyName("role")]
    public string Role { get; set; } = string.Empty;

    [JsonPropertyName("assigned_at")]
    public DateTime AssignedAt { get; set; }

    [JsonPropertyName("employee")]
    public EmployeeResponse? Employee { get; set; }
}

public class AddParticipantRequest
{
    [Required(ErrorMessage = "ID сотрудника обязателен")]
    public Guid EmployeeId { get; set; }

    [Required(ErrorMessage = "Роль обязательна")]
    [StringLength(50, ErrorMessage = "Роль не должна превышать 50 символов")]
    public string Role { get; set; } = "participant";
}

public enum TaskStatus
{
    Pending,
    InProgress,
    Completed,
    Cancelled
}

public enum TaskPriority
{
    Low,
    Medium,
    High,
    Critical
}

public static class TaskStatusExtensions
{
    public static string ToApiString(this TaskStatus status) => status switch
    {
        TaskStatus.Pending => "pending",
        TaskStatus.InProgress => "in_progress",
        TaskStatus.Completed => "completed",
        TaskStatus.Cancelled => "cancelled",
        _ => "pending"
    };

    public static string ToDisplayString(this TaskStatus status) => status switch
    {
        TaskStatus.Pending => "Ожидает",
        TaskStatus.InProgress => "В работе",
        TaskStatus.Completed => "Завершена",
        TaskStatus.Cancelled => "Отменена",
        _ => "Неизвестно"
    };
}

public static class TaskPriorityExtensions
{
    public static string ToApiString(this TaskPriority priority) => priority switch
    {
        TaskPriority.Low => "low",
        TaskPriority.Medium => "medium",
        TaskPriority.High => "high",
        TaskPriority.Critical => "critical",
        _ => "medium"
    };

    public static string ToDisplayString(this TaskPriority priority) => priority switch
    {
        TaskPriority.Low => "Низкий",
        TaskPriority.Medium => "Средний",
        TaskPriority.High => "Высокий",
        TaskPriority.Critical => "Критический",
        _ => "Средний"
    };
}
