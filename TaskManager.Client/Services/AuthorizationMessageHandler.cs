using Blazored.LocalStorage;
using System.Net.Http.Headers;

namespace TaskManager.Client.Services;

public class AuthorizationMessageHandler : DelegatingHandler
{
    private readonly ILocalStorageService _localStorage;

    public AuthorizationMessageHandler(ILocalStorageService localStorage)
    {
        _localStorage = localStorage;
    }

    protected override async Task<HttpResponseMessage> SendAsync(
        HttpRequestMessage request,
        CancellationToken cancellationToken)
    {
        // Получаем токен доступа из локального хранилища
        var accessToken = await _localStorage.GetItemAsync<string>("access_token");

        // ОТЛАДКА: Логируем для отладки
        Console.WriteLine($"[AuthHandler] URL: {request.RequestUri}");
        Console.WriteLine($"[AuthHandler] Токен: {(string.IsNullOrEmpty(accessToken) ? "ОТСУТСТВУЕТ" : $"ЕСТЬ (первые 20 символов: {accessToken.Substring(0, Math.Min(20, accessToken.Length))}...)")}");

        if (!string.IsNullOrEmpty(accessToken))
        {
            // Убираем возможные кавычки, которые могли быть добавлены при JSON сериализации
            accessToken = accessToken.Trim('"').Trim();
            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", accessToken);
            Console.WriteLine($"[AuthHandler] ✓ Authorization header добавлен");
        }
        else
        {
            Console.WriteLine($"[AuthHandler] ✗ ВНИМАНИЕ: Токен отсутствует!");
        }

        var response = await base.SendAsync(request, cancellationToken);

        // При получении 401 пытаемся обновить токен
        if (response.StatusCode == System.Net.HttpStatusCode.Unauthorized)
        {
            // TODO: Реализовать автоматическое обновление токена при необходимости
            // Пока просто возвращаем ошибку, и пользователю нужно войти заново
        }

        return response;
    }
}
