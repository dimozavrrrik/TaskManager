using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace TaskManager.Client.Models;

public class RegisterRequest
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

    [Required(ErrorMessage = "Пароль обязателен")]
    [StringLength(72, MinimumLength = 8, ErrorMessage = "Пароль должен быть от 8 до 72 символов")]
    public string Password { get; set; } = string.Empty;

    [Compare("Password", ErrorMessage = "Пароли не совпадают")]
    public string ConfirmPassword { get; set; } = string.Empty;
}

public class LoginRequest
{
    [Required(ErrorMessage = "Email обязателен")]
    [EmailAddress(ErrorMessage = "Некорректный email")]
    public string Email { get; set; } = string.Empty;

    [Required(ErrorMessage = "Пароль обязателен")]
    public string Password { get; set; } = string.Empty;
}

// RefreshTokenRequest и LogoutRequest больше не нужны - refresh token находится в HttpOnly cookie

public class AuthResponse
{
    [JsonPropertyName("access_token")]
    public string AccessToken { get; set; } = string.Empty;

    [JsonPropertyName("expires_at")]
    public DateTime ExpiresAt { get; set; }

    [JsonPropertyName("employee")]
    public EmployeeResponse Employee { get; set; } = null!;
}

public class TokenResponse
{
    [JsonPropertyName("access_token")]
    public string AccessToken { get; set; } = string.Empty;

    [JsonPropertyName("expires_at")]
    public DateTime ExpiresAt { get; set; }
}

public class ApiResponse<T>
{
    [JsonPropertyName("success")]
    public bool Success { get; set; }

    [JsonPropertyName("data")]
    public T? Data { get; set; }

    [JsonPropertyName("error")]
    public ErrorResponse? Error { get; set; }
}

public class ErrorResponse
{
    [JsonPropertyName("code")]
    public string Code { get; set; } = string.Empty;

    [JsonPropertyName("message")]
    public string Message { get; set; } = string.Empty;
}
