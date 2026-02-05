using System.Net.Http.Json;
using TaskManager.Client.Models;

namespace TaskManager.Client.Services;

public class TaskService : ITaskService
{
    private readonly IHttpClientFactory _httpClientFactory;

    public TaskService(IHttpClientFactory httpClientFactory)
    {
        _httpClientFactory = httpClientFactory;
    }

    private HttpClient CreateClient() => _httpClientFactory.CreateClient("TaskManager.API");

    public async Task<List<TaskResponse>> GetAllTasksAsync()
    {
        var client = CreateClient();
        var response = await client.GetAsync("tasks");
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<PaginatedResponse<TaskResponse>>>();
        return apiResponse?.Data?.Data ?? new List<TaskResponse>();
    }

    public async Task<TaskResponse> GetTaskByIdAsync(Guid id)
    {
        var client = CreateClient();
        var response = await client.GetAsync($"tasks/{id}");
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<TaskResponse>>();
        return apiResponse?.Data ?? throw new Exception("Задача не найдена");
    }

    public async Task<TaskResponse> CreateTaskAsync(CreateTaskRequest request)
    {
        var client = CreateClient();
        var response = await client.PostAsJsonAsync("tasks", request);
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<TaskResponse>>();
        return apiResponse?.Data ?? throw new Exception("Не удалось создать задачу");
    }

    public async Task<TaskResponse> UpdateTaskStatusAsync(Guid id, UpdateTaskStatusRequest request)
    {
        var client = CreateClient();
        var response = await client.PatchAsJsonAsync($"tasks/{id}/status", request);
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<TaskResponse>>();
        return apiResponse?.Data ?? throw new Exception("Не удалось обновить статус задачи");
    }

    public async Task ArchiveTaskAsync(Guid id)
    {
        var client = CreateClient();
        var response = await client.PatchAsync($"tasks/{id}/archive", null);
        response.EnsureSuccessStatusCode();
    }

    public async Task<List<TaskParticipantResponse>> GetTaskParticipantsAsync(Guid taskId)
    {
        var client = CreateClient();
        var response = await client.GetAsync($"tasks/{taskId}/participants");
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<List<TaskParticipantResponse>>>();
        return apiResponse?.Data ?? new List<TaskParticipantResponse>();
    }

    public async Task<TaskParticipantResponse> AddParticipantAsync(Guid taskId, AddParticipantRequest request)
    {
        var client = CreateClient();
        var response = await client.PostAsJsonAsync($"tasks/{taskId}/participants", request);
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<TaskParticipantResponse>>();
        return apiResponse?.Data ?? throw new Exception("Не удалось добавить участника");
    }

    public async Task<List<TaskResponse>> GetEmployeeTasksAsync(Guid employeeId)
    {
        var client = CreateClient();
        var response = await client.GetAsync($"employees/{employeeId}/tasks");
        response.EnsureSuccessStatusCode();

        var apiResponse = await response.Content.ReadFromJsonAsync<ApiResponse<PaginatedResponse<TaskResponse>>>();
        return apiResponse?.Data?.Data ?? new List<TaskResponse>();
    }
}
