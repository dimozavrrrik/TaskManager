using Microsoft.AspNetCore.Components.Authorization;
using System.Security.Claims;
using Blazored.LocalStorage;
using TaskManager.Client.Models;

namespace TaskManager.Client.Services;

public class CustomAuthenticationStateProvider : AuthenticationStateProvider
{
    private readonly ILocalStorageService _localStorage;
    private ClaimsPrincipal _anonymous = new ClaimsPrincipal(new ClaimsIdentity());

    public CustomAuthenticationStateProvider(ILocalStorageService localStorage)
    {
        _localStorage = localStorage;
    }

    public override async Task<AuthenticationState> GetAuthenticationStateAsync()
    {
        try
        {
            var accessToken = await _localStorage.GetItemAsync<string>("access_token");

            if (string.IsNullOrEmpty(accessToken))
            {
                return new AuthenticationState(_anonymous);
            }

            // Парсим JWT для получения claims
            var claims = ParseClaimsFromJwt(accessToken);
            var identity = new ClaimsIdentity(claims, "jwt");
            var user = new ClaimsPrincipal(identity);

            return new AuthenticationState(user);
        }
        catch
        {
            return new AuthenticationState(_anonymous);
        }
    }

    public void NotifyUserAuthentication(EmployeeResponse employee)
    {
        var claims = new List<Claim>
        {
            new Claim(ClaimTypes.NameIdentifier, employee.Id.ToString()),
            new Claim(ClaimTypes.Name, employee.Name),
            new Claim(ClaimTypes.Email, employee.Email)
        };

        var identity = new ClaimsIdentity(claims, "jwt");
        var user = new ClaimsPrincipal(identity);

        NotifyAuthenticationStateChanged(Task.FromResult(new AuthenticationState(user)));
    }

    public void NotifyUserLogout()
    {
        var authState = Task.FromResult(new AuthenticationState(_anonymous));
        NotifyAuthenticationStateChanged(authState);
    }

    private IEnumerable<Claim> ParseClaimsFromJwt(string jwt)
    {
        var claims = new List<Claim>();
        var payload = jwt.Split('.')[1];

        var base64 = payload.Replace('-', '+').Replace('_', '/');
        switch (base64.Length % 4)
        {
            case 2: base64 += "=="; break;
            case 3: base64 += "="; break;
        }

        var jsonBytes = Convert.FromBase64String(base64);
        var keyValuePairs = System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(jsonBytes);

        if (keyValuePairs != null)
        {
            keyValuePairs.TryGetValue("employee_id", out var employeeId);
            keyValuePairs.TryGetValue("email", out var email);
            keyValuePairs.TryGetValue("name", out var name);
            keyValuePairs.TryGetValue("exp", out var exp);

            if (employeeId != null)
                claims.Add(new Claim(ClaimTypes.NameIdentifier, employeeId.ToString() ?? string.Empty));

            if (email != null)
                claims.Add(new Claim(ClaimTypes.Email, email.ToString() ?? string.Empty));

            if (name != null)
                claims.Add(new Claim(ClaimTypes.Name, name.ToString() ?? string.Empty));

            if (exp != null)
                claims.Add(new Claim("exp", exp.ToString() ?? string.Empty));
        }

        return claims;
    }
}
