using System.Collections.ObjectModel;
using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using OrangeDesktop.Models;
using OrangeDesktop.Services;

namespace OrangeDesktop.ViewModels;

public partial class MainViewModel : ObservableObject
{
    private readonly ApiService _apiService;
    private readonly AuthService _authService;

    [ObservableProperty]
    private User? _currentUser;

    [ObservableProperty]
    private string _pageTitle = "首页";

    [ObservableProperty]
    private ObservableCollection<Blogger> _bloggers = [];

    [ObservableProperty]
    private ObservableCollection<Team> _teams = [];

    [ObservableProperty]
    private ObservableCollection<Notification> _notifications = [];

    [ObservableProperty]
    private BloggerStat? _bloggerStat;

    [ObservableProperty]
    private bool _isLoading;

    [ObservableProperty]
    private string _searchKeyword = "";

    public event Action? LogoutRequested;

    public MainViewModel(ApiService apiService, AuthService authService)
    {
        _apiService = apiService;
        _authService = authService;
        _currentUser = authService.CurrentUser;
    }

    [RelayCommand]
    private async Task LoadDataAsync()
    {
        IsLoading = true;
        try
        {
            await Task.WhenAll(
                LoadBloggersAsync(),
                LoadStatAsync(),
                LoadTeamsAsync(),
                LoadNotificationsAsync()
            );
        }
        finally
        {
            IsLoading = false;
        }
    }

    [RelayCommand]
    private async Task SearchAsync()
    {
        IsLoading = true;
        try
        {
            var bloggers = await _apiService.GetBloggerListAsync(1, 50, string.IsNullOrEmpty(SearchKeyword) ? null : SearchKeyword);
            Bloggers = new ObservableCollection<Blogger>(bloggers ?? []);
        }
        finally
        {
            IsLoading = false;
        }
    }

    [RelayCommand]
    private void Logout()
    {
        _authService.Logout();
        LogoutRequested?.Invoke();
    }

    [RelayCommand]
    private async Task RefreshProfileAsync()
    {
        await _authService.RefreshProfileAsync();
        CurrentUser = _authService.CurrentUser;
    }

    private async Task LoadBloggersAsync()
    {
        var bloggers = await _apiService.GetMyBloggersAsync(1, 50);
        Bloggers = new ObservableCollection<Blogger>(bloggers ?? []);
    }

    private async Task LoadStatAsync()
    {
        BloggerStat = await _apiService.GetBloggerStatAsync();
    }

    private async Task LoadTeamsAsync()
    {
        var teams = await _apiService.GetTeamsAsync();
        Teams = new ObservableCollection<Team>(teams ?? []);
    }

    private async Task LoadNotificationsAsync()
    {
        var notifications = await _apiService.GetNotificationsAsync();
        Notifications = new ObservableCollection<Notification>(notifications ?? []);
    }
}
