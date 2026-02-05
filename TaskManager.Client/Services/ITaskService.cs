using TaskManager.Client.Models;

namespace TaskManager.Client.Services;

public interface ITaskService
{
    Task<List<TaskResponse>> GetAllTasksAsync();
    Task<TaskResponse> GetTaskByIdAsync(Guid id);
    Task<TaskResponse> CreateTaskAsync(CreateTaskRequest request);
    Task<TaskResponse> UpdateTaskStatusAsync(Guid id, UpdateTaskStatusRequest request);
    Task ArchiveTaskAsync(Guid id);
    Task<List<TaskParticipantResponse>> GetTaskParticipantsAsync(Guid taskId);
    Task<TaskParticipantResponse> AddParticipantAsync(Guid taskId, AddParticipantRequest request);
    Task<List<TaskResponse>> GetEmployeeTasksAsync(Guid employeeId);
}
