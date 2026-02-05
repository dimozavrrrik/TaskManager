using TaskManager.Client.Models;

namespace TaskManager.Client.Services;

public interface IAuthService
{
    Task<AuthResponse> RegisterAsync(RegisterRequest request);
    Task<AuthResponse> LoginAsync(LoginRequest request);
    Task<TokenResponse> RefreshTokenAsync();
    Task LogoutAsync();
    Task<EmployeeResponse?> GetCurrentUserAsync();
}
