using System.Windows;
using OrangeDesktop.Services;
using OrangeDesktop.ViewModels;
using OrangeDesktop.Views;

namespace OrangeDesktop;

public partial class App : Application
{
    private ApiService _apiService = null!;
    private AuthService _authService = null!;
    private SettingsService _settingsService = null!;

    private void OnStartup(object sender, StartupEventArgs e)
    {
        _apiService = new ApiService();
        _settingsService = new SettingsService();
        _settingsService.Load();
        _apiService.SetBaseUrl(_settingsService.Settings.ServerUrl);

        _authService = new AuthService(_apiService);

        if (_authService.IsLoggedIn)
        {
            ShowMainWindow();
        }
        else
        {
            ShowLoginWindow();
        }
    }

    private void ShowLoginWindow()
    {
        var viewModel = new LoginViewModel(_apiService, _authService, _settingsService);
        var loginWindow = new LoginWindow(viewModel);
        viewModel.LoginSucceeded += () =>
        {
            loginWindow.DialogResult = true;
        };

        if (loginWindow.ShowDialog() != true)
        {
            Shutdown();
            return;
        }

        ShowMainWindow();
    }

    private void ShowMainWindow()
    {
        var viewModel = new MainViewModel(_apiService, _authService);
        var mainWindow = new MainWindowView(viewModel);
        viewModel.LogoutRequested += () =>
        {
            mainWindow.Close();
            ShowLoginWindow();
        };
        mainWindow.Show();
    }
}
