using System.IO;
using System.Text.Json;
using OrangeDesktop.Models;

namespace OrangeDesktop.Services;

public class AuthService
{
    private const string AuthFile = "auth.json";
    private readonly ApiService _apiService;

    public User? CurrentUser { get; private set; }
    public string? Token { get; private set; }
    public bool IsLoggedIn => Token != null && CurrentUser != null;

    public event Action? AuthStateChanged;

    public AuthService(ApiService apiService)
    {
        _apiService = apiService;
        LoadSavedAuth();
    }

    public async Task<bool> LoginAsync(string username, string password)
    {
        try
        {
            var response = await _apiService.LoginAsync(new LoginRequest
            {
                Username = username,
                Password = password
            });

            if (response != null && !string.IsNullOrEmpty(response.Token))
            {
                Token = response.Token;
                CurrentUser = response.User;
                _apiService.SetToken(Token);
                SaveAuth();
                AuthStateChanged?.Invoke();
                return true;
            }
            return false;
        }
        catch
        {
            return false;
        }
    }

    public async Task<bool> RegisterAsync(string username, string password, string realName)
    {
        try
        {
            var response = await _apiService.RegisterAsync(new RegisterRequest
            {
                Username = username,
                Password = password,
                RealName = realName
            });

            if (response != null && !string.IsNullOrEmpty(response.Token))
            {
                Token = response.Token;
                CurrentUser = response.User;
                _apiService.SetToken(Token);
                SaveAuth();
                AuthStateChanged?.Invoke();
                return true;
            }
            return false;
        }
        catch
        {
            return false;
        }
    }

    public void Logout()
    {
        Token = null;
        CurrentUser = null;
        _apiService.ClearToken();
        ClearSavedAuth();
        AuthStateChanged?.Invoke();
    }

    public async Task<bool> RefreshProfileAsync()
    {
        if (!IsLoggedIn) return false;
        try
        {
            var profile = await _apiService.GetProfileAsync();
            if (profile != null)
            {
                CurrentUser = profile;
                SaveAuth();
                AuthStateChanged?.Invoke();
                return true;
            }
            return false;
        }
        catch
        {
            return false;
        }
    }

    private void SaveAuth()
    {
        var data = new AuthData { Token = Token!, User = CurrentUser! };
        var json = JsonSerializer.Serialize(data, new JsonSerializerOptions { PropertyNamingPolicy = JsonNamingPolicy.CamelCase });
        File.WriteAllText(AuthFile, json);
    }

    private void LoadSavedAuth()
    {
        if (!File.Exists(AuthFile)) return;
        try
        {
            var json = File.ReadAllText(AuthFile);
            var data = JsonSerializer.Deserialize<AuthData>(json, new JsonSerializerOptions { PropertyNameCaseInsensitive = true });
            if (data != null && !string.IsNullOrEmpty(data.Token))
            {
                Token = data.Token;
                CurrentUser = data.User;
                _apiService.SetToken(Token);
            }
        }
        catch
        {
            Token = null;
            CurrentUser = null;
        }
    }

    private void ClearSavedAuth()
    {
        if (File.Exists(AuthFile))
            File.Delete(AuthFile);
    }

    private class AuthData
    {
        public string Token { get; set; } = "";
        public User User { get; set; } = new();
    }
}
