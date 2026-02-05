using Microsoft.AspNetCore.Components.WebAssembly.Http;

namespace TaskManager.Client.Services;

/// <summary>
/// DelegatingHandler для включения cookies (credentials) в каждый HTTP запрос
/// </summary>
public class CookieHandler : DelegatingHandler
{
    protected override Task<HttpResponseMessage> SendAsync(
        HttpRequestMessage request,
        CancellationToken cancellationToken)
    {
        // Включаем credentials для отправки cookies
        // Это необходимо для работы с HttpOnly cookies
        request.SetBrowserRequestCredentials(BrowserRequestCredentials.Include);

        return base.SendAsync(request, cancellationToken);
    }
}
