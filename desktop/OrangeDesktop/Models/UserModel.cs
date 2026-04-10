namespace OrangeDesktop.Models;

public class User
{
    public int Id { get; set; }
    public string Username { get; set; } = "";
    public string RealName { get; set; } = "";
    public string Role { get; set; } = "";
    public string Avatar { get; set; } = "";
    public string Email { get; set; } = "";
    public string Phone { get; set; } = "";
    public string Bio { get; set; } = "";
    public string Status { get; set; } = "";
    public long? TeamId { get; set; }
    public string TeamName { get; set; } = "";
    public string TeamColor { get; set; } = "";
    public string CreateTime { get; set; } = "";
}
