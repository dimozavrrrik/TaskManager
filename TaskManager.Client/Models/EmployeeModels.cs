using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace TaskManager.Client.Models;

public class EmployeeResponse
{
    [JsonPropertyName("id")]
    public Guid Id { get; set; }

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [JsonPropertyName("department")]
    public string Department { get; set; } = string.Empty;

    [JsonPropertyName("position")]
    public string Position { get; set; } = string.Empty;

    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("created_at")]
    public DateTime CreatedAt { get; set; }

    [JsonPropertyName("updated_at")]
    public DateTime UpdatedAt { get; set; }

    [JsonPropertyName("deleted_at")]
    public DateTime? DeletedAt { get; set; }
}

public class CreateEmployeeRequest
{
    [Required(ErrorMessage = "Имя обязательно")]
    [StringLength(255, MinimumLength = 2, ErrorMessage = "Имя должно быть от 2 до 255 символов")]
    public string Name { get; set; } = string.Empty;

    [Required(ErrorMessage = "Отдел обязателен")]
    [StringLength(100, ErrorMessage = "Отдел не должен превышать 100 символов")]
    public string Department { get; set; } = string.Empty;

    [Required(ErrorMessage = "Должность обязательна")]
    [StringLength(100, ErrorMessage = "Должность не должна превышать 100 символов")]
    public string Position { get; set; } = string.Empty;

    [Required(ErrorMessage = "Email обязателен")]
    [EmailAddress(ErrorMessage = "Некорректный email")]
    public string Email { get; set; } = string.Empty;
}

public class UpdateEmployeeRequest
{
    [StringLength(255, MinimumLength = 2, ErrorMessage = "Имя должно быть от 2 до 255 символов")]
    public string? Name { get; set; }

    [StringLength(100, ErrorMessage = "Отдел не должен превышать 100 символов")]
    public string? Department { get; set; }

    [StringLength(100, ErrorMessage = "Должность не должна превышать 100 символов")]
    public string? Position { get; set; }

    [EmailAddress(ErrorMessage = "Некорректный email")]
    public string? Email { get; set; }
}
