using TaskManager.Client.Models;

namespace TaskManager.Client.Services;

public interface IEmployeeService
{
    Task<List<EmployeeResponse>> GetAllEmployeesAsync();
    Task<EmployeeResponse> GetEmployeeByIdAsync(Guid id);
    Task<EmployeeResponse> CreateEmployeeAsync(CreateEmployeeRequest request);
    Task<EmployeeResponse> UpdateEmployeeAsync(Guid id, UpdateEmployeeRequest request);
    Task DeleteEmployeeAsync(Guid id);
}
