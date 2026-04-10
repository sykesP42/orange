namespace OrangeDesktop.Models;

public class Blogger
{
    public int Id { get; set; }
    public string Nickname { get; set; } = "";
    public string Category { get; set; } = "";
    public string Product { get; set; } = "";
    public string Contact { get; set; } = "";
    public string Wechat { get; set; } = "";
    public string CustomContact { get; set; } = "";
    public string Platform { get; set; } = "";
    public string PlatformAccount { get; set; } = "";
    public string Description { get; set; } = "";
    public string Tags { get; set; } = "";
    public string Avatar { get; set; } = "";
    public string UserName { get; set; } = "";
    public string RealName { get; set; } = "";
    public string Status { get; set; } = "";
    public string StatusUpdateTime { get; set; } = "";
    public string StatusRemark { get; set; } = "";
    public int IsDeleted { get; set; }
    public int IsInvalid { get; set; }
    public string CreateTime { get; set; } = "";
    public string UpdateTime { get; set; } = "";
    public long? TeamId { get; set; }
    public string TeamName { get; set; } = "";
    public string TeamColor { get; set; } = "";
}

public class BloggerStat
{
    public int Total { get; set; }
    public int Active { get; set; }
    public int Inactive { get; set; }
    public int Invalid { get; set; }
}

public class Followup
{
    public int Id { get; set; }
    public int BloggerId { get; set; }
    public string Content { get; set; } = "";
    public string Type { get; set; } = "";
    public string Operator { get; set; } = "";
    public string CreateTime { get; set; } = "";
}

public class Cooperation
{
    public int Id { get; set; }
    public int BloggerId { get; set; }
    public string Title { get; set; } = "";
    public string Date { get; set; } = "";
    public string Status { get; set; } = "";
    public string Product { get; set; } = "";
    public double Amount { get; set; }
    public string Note { get; set; } = "";
    public string CreateTime { get; set; } = "";
    public string UpdateTime { get; set; } = "";
}
