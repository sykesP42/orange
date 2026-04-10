using System.Net.Http;
using System.Net.Http.Headers;
using System.Text.Json;
using OrangeDesktop.Models;

namespace OrangeDesktop.Services;

public class ApiService
{
    private readonly HttpClient _httpClient;
    private readonly JsonSerializerOptions _jsonOptions;

    public ApiService()
    {
        _httpClient = new HttpClient
        {
            Timeout = TimeSpan.FromSeconds(30)
        };

        _jsonOptions = new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true,
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase
        };
    }

    public void SetBaseUrl(string baseUrl)
    {
        _httpClient.BaseAddress = new Uri(baseUrl.TrimEnd('/'));
    }

    public void SetToken(string token)
    {
        _httpClient.DefaultRequestHeaders.Authorization =
            new AuthenticationHeaderValue("Bearer", token);
    }

    public void ClearToken()
    {
        _httpClient.DefaultRequestHeaders.Authorization = null;
    }

    private async Task<T?> GetAsync<T>(string endpoint)
    {
        var response = await _httpClient.GetAsync(endpoint);
        response.EnsureSuccessStatusCode();
        var json = await response.Content.ReadAsStringAsync();
        return JsonSerializer.Deserialize<T>(json, _jsonOptions);
    }

    private async Task<T?> PostAsync<T>(string endpoint, object? data = null)
    {
        var json = data != null
            ? JsonSerializer.Serialize(data, _jsonOptions)
            : null;
        var content = json != null
            ? new StringContent(json, System.Text.Encoding.UTF8, "application/json")
            : null;

        var response = content != null
            ? await _httpClient.PostAsync(endpoint, content)
            : await _httpClient.PostAsync(endpoint, null);
        response.EnsureSuccessStatusCode();
        var responseJson = await response.Content.ReadAsStringAsync();
        return JsonSerializer.Deserialize<T>(responseJson, _jsonOptions);
    }

    private async Task<T?> PutAsync<T>(string endpoint, object data)
    {
        var json = JsonSerializer.Serialize(data, _jsonOptions);
        var content = new StringContent(json, System.Text.Encoding.UTF8, "application/json");
        var response = await _httpClient.PutAsync(endpoint, content);
        response.EnsureSuccessStatusCode();
        var responseJson = await response.Content.ReadAsStringAsync();
        return JsonSerializer.Deserialize<T>(responseJson, _jsonOptions);
    }

    private async Task DeleteAsync(string endpoint)
    {
        var response = await _httpClient.DeleteAsync(endpoint);
        response.EnsureSuccessStatusCode();
    }

    public async Task<LoginResponse?> LoginAsync(LoginRequest request)
    {
        return await PostAsync<LoginResponse>("/api/login", request);
    }

    public async Task<LoginResponse?> RegisterAsync(RegisterRequest request)
    {
        return await PostAsync<LoginResponse>("/api/register", request);
    }

    public async Task<User?> GetProfileAsync()
    {
        return await GetAsync<User>("/api/user/profile");
    }

    public async Task<User?> UpdateProfileAsync(object data)
    {
        return await PutAsync<User>("/api/user/profile", data);
    }

    public async Task<List<Blogger>?> GetMyBloggersAsync(int page = 1, int pageSize = 20)
    {
        return await GetAsync<List<Blogger>>($"/api/blogger/my?page={page}&pageSize={pageSize}");
    }

    public async Task<List<Blogger>?> GetBloggerListAsync(int page = 1, int pageSize = 20, string? keyword = null)
    {
        var url = $"/api/blogger/list?page={page}&pageSize={pageSize}";
        if (!string.IsNullOrEmpty(keyword))
            url += $"&keyword={Uri.EscapeDataString(keyword)}";
        return await GetAsync<List<Blogger>>(url);
    }

    public async Task<Blogger?> GetBloggerAsync(int id)
    {
        return await GetAsync<Blogger>($"/api/blogger/{id}");
    }

    public async Task<Blogger?> CreateBloggerAsync(object data)
    {
        return await PostAsync<Blogger>("/api/blogger/add", data);
    }

    public async Task<Blogger?> UpdateBloggerAsync(int id, object data)
    {
        return await PutAsync<Blogger>($"/api/blogger/{id}", data);
    }

    public async Task DeleteBloggerAsync(int id)
    {
        await DeleteAsync($"/api/blogger/{id}");
    }

    public async Task<BloggerStat?> GetBloggerStatAsync()
    {
        return await GetAsync<BloggerStat>("/api/blogger/stat");
    }

    public async Task<List<Followup>?> GetFollowupsAsync(int bloggerId)
    {
        return await GetAsync<List<Followup>>($"/api/blogger/{bloggerId}/followups");
    }

    public async Task<Followup?> CreateFollowupAsync(int bloggerId, object data)
    {
        return await PostAsync<Followup>($"/api/blogger/{bloggerId}/followups", data);
    }

    public async Task<List<Cooperation>?> GetCooperationsAsync(int bloggerId)
    {
        return await GetAsync<List<Cooperation>>($"/api/blogger/{bloggerId}/cooperations");
    }

    public async Task<Cooperation?> CreateCooperationAsync(int bloggerId, object data)
    {
        return await PostAsync<Cooperation>($"/api/blogger/{bloggerId}/cooperations", data);
    }

    public async Task<List<Team>?> GetTeamsAsync()
    {
        return await GetAsync<List<Team>>("/api/teams");
    }

    public async Task<Team?> GetTeamAsync(int id)
    {
        return await GetAsync<Team>($"/api/team/{id}");
    }

    public async Task<Team?> GetMyTeamAsync()
    {
        return await GetAsync<Team>("/api/my/team");
    }

    public async Task<List<TeamPost>?> GetTeamPostsAsync(int teamId, int page = 1, int pageSize = 20)
    {
        return await GetAsync<List<TeamPost>>($"/api/team/{teamId}/posts?page={page}&pageSize={pageSize}");
    }

    public async Task<List<TeamMessage>?> GetTeamMessagesAsync(int teamId)
    {
        return await GetAsync<List<TeamMessage>>($"/api/team/{teamId}/messages");
    }

    public async Task<List<Notification>?> GetNotificationsAsync()
    {
        return await GetAsync<List<Notification>>("/api/notifications");
    }

    public async Task<List<Category>?> GetCategoriesAsync()
    {
        return await GetAsync<List<Category>>("/api/category/list");
    }

    public async Task<List<Platform>?> GetPlatformsAsync()
    {
        return await GetAsync<List<Platform>>("/api/platform/list");
    }

    public async Task<List<Product>?> GetProductsAsync()
    {
        return await GetAsync<List<Product>>("/api/product/list");
    }

    public async Task<List<User>?> GetPublicUsersAsync()
    {
        return await GetAsync<List<User>>("/api/users/public");
    }

    public async Task<List<Announcement>?> GetAnnouncementsAsync()
    {
        return await GetAsync<List<Announcement>>("/api/announcements");
    }

    public async Task<bool> CheckHealthAsync()
    {
        try
        {
            var response = await _httpClient.GetAsync("/health");
            return response.IsSuccessStatusCode;
        }
        catch
        {
            return false;
        }
    }
}
