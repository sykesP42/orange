namespace OrangeDesktop.Models;

public class Team
{
    public int Id { get; set; }
    public string Name { get; set; } = "";
    public string Description { get; set; } = "";
    public string Logo { get; set; } = "";
    public string BgImage { get; set; } = "";
    public string Color { get; set; } = "";
    public int LeaderId { get; set; }
    public string LeaderName { get; set; } = "";
    public string AdminIds { get; set; } = "";
    public string Announcement { get; set; } = "";
    public int IsPublic { get; set; }
    public int NeedApprove { get; set; }
    public int MemberCount { get; set; }
    public string InviteCode { get; set; } = "";
    public string CreateTime { get; set; } = "";
    public string UpdateTime { get; set; } = "";
}

public class TeamPost
{
    public int Id { get; set; }
    public int TeamId { get; set; }
    public int AuthorId { get; set; }
    public string AuthorName { get; set; } = "";
    public string AuthorUsername { get; set; } = "";
    public string AuthorAvatar { get; set; } = "";
    public string Title { get; set; } = "";
    public string Content { get; set; } = "";
    public string Category { get; set; } = "";
    public int CommentCount { get; set; }
    public int LikeCount { get; set; }
    public int ViewCount { get; set; }
    public int IsPinned { get; set; }
    public int IsFeatured { get; set; }
    public string CreateTime { get; set; } = "";
    public string UpdateTime { get; set; } = "";
}

public class TeamMessage
{
    public int Id { get; set; }
    public int TeamId { get; set; }
    public int SenderId { get; set; }
    public string SenderName { get; set; } = "";
    public string SenderAvatar { get; set; } = "";
    public string Title { get; set; } = "";
    public string Content { get; set; } = "";
    public string Type { get; set; } = "";
    public int IsRead { get; set; }
    public string CreateTime { get; set; } = "";
}

public class Notification
{
    public int Id { get; set; }
    public int UserId { get; set; }
    public string Type { get; set; } = "";
    public string Title { get; set; } = "";
    public string Content { get; set; } = "";
    public string Priority { get; set; } = "";
    public long? TeamId { get; set; }
    public string FromUser { get; set; } = "";
    public long? BloggerId { get; set; }
    public long? PostId { get; set; }
    public int IsRead { get; set; }
    public string CreateTime { get; set; } = "";
}
