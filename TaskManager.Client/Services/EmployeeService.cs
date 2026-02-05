using System.Net.Http.Json;
using TaskManager.Client.Models;

namespace TaskManager.Client.Services;

public class EmployeeService : IEmployeeService
{
    private readonly IHttpClientFactory _httpClientFactory;

    public EmployeeService(IHttpClientFactory httpClientFactory)
    {
        _httpClientFactory = httpClientFactory;
    }

    private HttpClient CreateClient() => _httpClientFactory.CreateClient("TaskManager.API");

    public async Task<List<EmployeeResponse>> GetAllEmployeesAsync()
    {
        var client = CreateClient();
        var response = await client.GetAsync("employees");
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<PaginatedResponse<EmployeeResponse>>>();
        return apiResponse?.Data?.Data ?? new List<EmployeeResponse>();
    }

    public async Task<EmployeeResponse> GetEmployeeByIdAsync(Guid id)
    {
        var client = CreateClient();
        var response = await client.GetAsync($"employees/{id}");
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<EmployeeResponse>>();
        return apiResponse?.Data ?? throw new Exception("Сотрудник не найден");
    }

    public async Task<EmployeeResponse> CreateEmployeeAsync(CreateEmployeeRequest request)
    {
        var client = CreateClient();
        var response = await client.PostAsJsonAsync("employees", request);
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<EmployeeResponse>>();
        return apiResponse?.Data ?? throw new Exception("Не удалось создать сотрудника");
    }

    public async Task<EmployeeResponse> UpdateEmployeeAsync(Guid id, UpdateEmployeeRequest request)
    {
        var client = CreateClient();
        var response = await client.PutAsJsonAsync($"employees/{id}", request);
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<EmployeeResponse>>();
        return apiResponse?.Data ?? throw new Exception("Не удалось обновить сотрудника");
    }

    public async Task DeleteEmployeeAsync(Guid id)
    {
        var client = CreateClient();
        var response = await client.DeleteAsync($"employees/{id}");
        response.EnsureSuccessStatusCode();
    }
}
