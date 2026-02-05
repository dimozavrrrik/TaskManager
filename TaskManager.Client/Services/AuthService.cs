using System.Net.Http.Json;
using Blazored.LocalStorage;
using TaskManager.Client.Models;
using Microsoft.AspNetCore.Components.Authorization;

namespace TaskManager.Client.Services;

public class AuthService : IAuthService
{
    private readonly HttpClient _httpClient;
    private readonly ILocalStorageService _localStorage;
    private readonly AuthenticationStateProvider _authStateProvider;

    public AuthService(
        HttpClient httpClient,
        ILocalStorageService localStorage,
        AuthenticationStateProvider authStateProvider)
    {
        _httpClient = httpClient;
        _localStorage = localStorage;
        _authStateProvider = authStateProvider;
    }

    public async Task<AuthResponse> RegisterAsync(RegisterRequest request)
    {
        var response = await _httpClient.PostAsJsonAsync("auth/register", request);

        if (!response.IsSuccessStatusCode)
        {
            var errorContent = await response.Content.ReadAsStringAsync();
            throw new HttpRequestException($"Ошибка регистрации: {errorContent}");
        }

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<AuthResponse>>();

        if (apiResponse?.Success != true || apiResponse.Data == null)
        {
            throw new Exception(apiResponse?.Error?.Message ?? "Ошибка регистрации");
        }

        var authResponse = apiResponse.Data;

        // ОТЛАДКА: Логируем токен
        Console.WriteLine($"[AuthService.Register] Получен токен (первые 30 символов): {authResponse.AccessToken.Substring(0, Math.Min(30, authResponse.AccessToken.Length))}...");
        Console.WriteLine($"[AuthService.Register] Длина токена: {authResponse.AccessToken.Length}");

        // Сохраняем только access token (refresh token находится в HttpOnly cookie)
        await _localStorage.SetItemAsync("access_token", authResponse.AccessToken);
        await _localStorage.SetItemAsync("expires_at", authResponse.ExpiresAt);

        Console.WriteLine($"[AuthService.Register] ✓ Токен сохранён в localStorage");

        // Уведомляем об изменении состояния аутентификации
        ((CustomAuthenticationStateProvider)_authStateProvider).NotifyUserAuthentication(authResponse.Employee);

        return authResponse;
    }

    public async Task<AuthResponse> LoginAsync(LoginRequest request)
    {
        var response = await _httpClient.PostAsJsonAsync("auth/login", request);

        if (!response.IsSuccessStatusCode)
        {
            var errorContent = await response.Content.ReadAsStringAsync();
            throw new HttpRequestException($"Ошибка входа: {errorContent}");
        }

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<AuthResponse>>();

        if (apiResponse?.Success != true || apiResponse.Data == null)
        {
            throw new Exception(apiResponse?.Error?.Message ?? "Ошибка входа");
        }

        var authResponse = apiResponse.Data;

        // ОТЛАДКА: Логируем токен
        Console.WriteLine($"[AuthService.Login] Получен токен (первые 30 символов): {authResponse.AccessToken.Substring(0, Math.Min(30, authResponse.AccessToken.Length))}...");
        Console.WriteLine($"[AuthService.Login] Длина токена: {authResponse.AccessToken.Length}");

        // Сохраняем только access token (refresh token находится в HttpOnly cookie)
        await _localStorage.SetItemAsync("access_token", authResponse.AccessToken);
        await _localStorage.SetItemAsync("expires_at", authResponse.ExpiresAt);

        Console.WriteLine($"[AuthService.Login] ✓ Токен сохранён в localStorage");

        // Уведомляем об изменении состояния аутентификации
        ((CustomAuthenticationStateProvider)_authStateProvider).NotifyUserAuthentication(authResponse.Employee);

        return authResponse;
    }

    public async Task<TokenResponse> RefreshTokenAsync()
    {
        // Refresh token автоматически отправляется в HttpOnly cookie
        var response = await _httpClient.PostAsJsonAsync("auth/refresh", new { });

        if (!response.IsSuccessStatusCode)
        {
            await LogoutAsync();
            throw new HttpRequestException("Не удалось обновить токен");
        }

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<TokenResponse>>();

        if (apiResponse?.Success != true || apiResponse.Data == null)
        {
            await LogoutAsync();
            throw new Exception(apiResponse?.Error?.Message ?? "Не удалось обновить токен");
        }

        var tokenResponse = apiResponse.Data;

        // Обновляем только access token (refresh token находится в HttpOnly cookie)
        await _localStorage.SetItemAsync("access_token", tokenResponse.AccessToken);
        await _localStorage.SetItemAsync("expires_at", tokenResponse.ExpiresAt);

        return tokenResponse;
    }

    public async Task LogoutAsync()
    {
        try
        {
            // Refresh token автоматически отправляется через cookie
            await _httpClient.PostAsJsonAsync("auth/logout", new { });
        }
        catch
        {
            // Игнорируем ошибки при выходе
        }
        finally
        {
            // Очищаем локальное хранилище (только access token)
            await _localStorage.RemoveItemAsync("access_token");
            await _localStorage.RemoveItemAsync("expires_at");

            // Уведомляем об изменении состояния аутентификации
            ((CustomAuthenticationStateProvider)_authStateProvider).NotifyUserLogout();
        }
    }

    public async Task<EmployeeResponse?> GetCurrentUserAsync()
    {
        var accessToken = await _localStorage.GetItemAsync<string>("access_token");

        if (string.IsNullOrEmpty(accessToken))
        {
            return null;
        }

        // Парсим JWT для получения данных сотрудника
        // В продакшене рекомендуется вызывать отдельный API эндпоинт
        try
        {
            var payload = ParseJwtPayload(accessToken);
            return payload;
        }
        catch
        {
            return null;
        }
    }

    private EmployeeResponse? ParseJwtPayload(string token)
    {
        var parts = token.Split('.');
        if (parts.Length != 3)
        {
            return null;
        }

        var payload = parts[1];
        var base64 = payload.Replace('-', '+').Replace('_', '/');

        switch (base64.Length % 4)
        {
            case 2: base64 += "=="; break;
            case 3: base64 += "="; break;
        }

        var json = System.Text.Encoding.UTF8.GetString(Convert.FromBase64String(base64));
        var claims = System.Text.Json.JsonSerializer.Deserialize<System.Text.Json.JsonElement>(json);

        if (claims.TryGetProperty("employee_id", out var employeeId) &&
            claims.TryGetProperty("email", out var email) &&
            claims.TryGetProperty("name", out var name))
        {
            return new EmployeeResponse
            {
                Id = Guid.Parse(employeeId.GetString() ?? Guid.Empty.ToString()),
                Email = email.GetString() ?? string.Empty,
                Name = name.GetString() ?? string.Empty,
                Department = string.Empty,
                Position = string.Empty
            };
        }

        return null;
    }
}
