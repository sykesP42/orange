using System.IO;
using System.Text.Json;

namespace OrangeDesktop.Services;

public class AppSettings
{
    public string ServerUrl { get; set; } = "http://localhost:8080";
}

public class SettingsService
{
    private const string SettingsFile = "settings.json";
    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
        WriteIndented = true
    };

    public AppSettings Settings { get; private set; } = new();

    public event Action? SettingsChanged;

    public void Load()
    {
        if (!File.Exists(SettingsFile)) return;
        try
        {
            var json = File.ReadAllText(SettingsFile);
            var settings = JsonSerializer.Deserialize<AppSettings>(json, new JsonSerializerOptions { PropertyNameCaseInsensitive = true });
            if (settings != null)
                Settings = settings;
        }
        catch
        {
        }
    }

    public void Save()
    {
        var json = JsonSerializer.Serialize(Settings, JsonOptions);
        File.WriteAllText(SettingsFile, json);
        SettingsChanged?.Invoke();
    }

    public void UpdateServerUrl(string url)
    {
        Settings.ServerUrl = url.TrimEnd('/');
        Save();
    }
}
