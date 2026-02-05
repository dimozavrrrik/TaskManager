using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using TaskManager.Client;
using TaskManager.Client.Services;
using MudBlazor.Services;
using Blazored.LocalStorage;
using Microsoft.AspNetCore.Components.Authorization;

var builder = WebAssemblyHostBuilder.CreateDefault(args);
builder.RootComponents.Add<App>("#app");
builder.RootComponents.Add<HeadOutlet>("head::after");

// Читаем базовый URL API из конфигурации
var apiBaseUrlPath = builder.Configuration["ApiBaseUrl"] ?? "/api/v1";

// Формируем полный URL API (относительно базового адреса приложения)
// Завершающий слэш ОБЯЗАТЕЛЕН для корректного добавления относительных путей в HttpClient
var baseAddress = new Uri(builder.HostEnvironment.BaseAddress);
var apiBaseUrl = new Uri(baseAddress, apiBaseUrlPath).ToString().TrimEnd('/') + "/";

// Добавляем сервисы MudBlazor
builder.Services.AddMudServices();

// Добавляем Blazored.LocalStorage
builder.Services.AddBlazoredLocalStorage();

// Добавляем обработчики сообщений
builder.Services.AddTransient<AuthorizationMessageHandler>();
builder.Services.AddTransient<CookieHandler>();

// Настраиваем HttpClient по умолчанию для AuthService (только cookies, без auth header)
builder.Services.AddHttpClient("AuthClient", client =>
{
    client.BaseAddress = new Uri(apiBaseUrl);
})
.AddHttpMessageHandler<CookieHandler>();

builder.Services.AddScoped(sp =>
{
    var factory = sp.GetRequiredService<IHttpClientFactory>();
    return factory.CreateClient("AuthClient");
});

// Настраиваем аутентифицированный HttpClient для защищенных эндпоинтов
// ВАЖНО: порядок handlers - сначала Authorization, потом Cookie!
builder.Services.AddHttpClient("TaskManager.API", client =>
{
    client.BaseAddress = new Uri(apiBaseUrl);
})
.AddHttpMessageHandler<AuthorizationMessageHandler>()
.AddHttpMessageHandler<CookieHandler>();

// Регистрируем сервисы
builder.Services.AddScoped<IAuthService, AuthService>();
builder.Services.AddScoped<ITaskService, TaskService>();
builder.Services.AddScoped<IEmployeeService, EmployeeService>();

// Добавляем сервисы аутентификации
builder.Services.AddAuthorizationCore();
builder.Services.AddScoped<AuthenticationStateProvider, CustomAuthenticationStateProvider>();

await builder.Build().RunAsync();
