using System.Net.Http;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using OrangeDesktop.Services;

namespace OrangeDesktop.ViewModels;

public partial class LoginViewModel : ObservableObject
{
    private readonly AuthService _authService;
    private readonly ApiService _apiService;
    private readonly SettingsService _settingsService;

    [ObservableProperty]
    private string _username = "";

    [ObservableProperty]
    private string _password = "";

    [ObservableProperty]
    private string _serverUrl = "";

    [ObservableProperty]
    private string _errorMessage = "";

    [ObservableProperty]
    private bool _isLoading;

    [ObservableProperty]
    private bool _isRegisterMode;

    [ObservableProperty]
    private string _realName = "";

    public event Action? LoginSucceeded;

    public LoginViewModel(ApiService apiService, AuthService authService, SettingsService settingsService)
    {
        _apiService = apiService;
        _authService = authService;
        _settingsService = settingsService;
        _serverUrl = _settingsService.Settings.ServerUrl;
    }

    [RelayCommand]
    private async Task LoginAsync()
    {
        if (string.IsNullOrWhiteSpace(Username) || string.IsNullOrWhiteSpace(Password))
        {
            ErrorMessage = "请输入用户名和密码";
            return;
        }

        await ExecuteWithLoadingAsync(async () =>
        {
            UpdateServerUrl();
            var success = await _authService.LoginAsync(Username, Password);
            if (success)
            {
                LoginSucceeded?.Invoke();
            }
            else
            {
                ErrorMessage = "登录失败，请检查用户名和密码";
            }
        });
    }

    [RelayCommand]
    private async Task RegisterAsync()
    {
        if (string.IsNullOrWhiteSpace(Username) || string.IsNullOrWhiteSpace(Password))
        {
            ErrorMessage = "请输入用户名和密码";
            return;
        }

        if (string.IsNullOrWhiteSpace(RealName))
        {
            ErrorMessage = "请输入真实姓名";
            return;
        }

        await ExecuteWithLoadingAsync(async () =>
        {
            UpdateServerUrl();
            var success = await _authService.RegisterAsync(Username, Password, RealName);
            if (success)
            {
                LoginSucceeded?.Invoke();
            }
            else
            {
                ErrorMessage = "注册失败，用户名可能已存在";
            }
        });
    }

    [RelayCommand]
    private void ToggleMode()
    {
        IsRegisterMode = !IsRegisterMode;
        ErrorMessage = "";
    }

    [RelayCommand]
    private async Task TestConnectionAsync()
    {
        await ExecuteWithLoadingAsync(async () =>
        {
            UpdateServerUrl();
            _apiService.SetBaseUrl(ServerUrl);
            var healthy = await _apiService.CheckHealthAsync();
            ErrorMessage = healthy ? "服务器连接成功 ✓" : "无法连接到服务器 ✗";
        });
    }

    private void UpdateServerUrl()
    {
        if (_settingsService.Settings.ServerUrl != ServerUrl)
        {
            _settingsService.UpdateServerUrl(ServerUrl);
        }
        _apiService.SetBaseUrl(ServerUrl);
    }

    private async Task ExecuteWithLoadingAsync(Func<Task> action)
    {
        IsLoading = true;
        ErrorMessage = "";
        try
        {
            await action();
        }
        catch (HttpRequestException ex)
        {
            ErrorMessage = $"网络错误: {ex.Message}";
        }
        catch (Exception ex)
        {
            ErrorMessage = $"发生错误: {ex.Message}";
        }
        finally
        {
            IsLoading = false;
        }
    }
}
