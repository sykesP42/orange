namespace OrangeDesktop.Models;

public class PaginatedResult<T>
{
    public List<T> Items { get; set; } = [];
    public int Total { get; set; }
    public int Page { get; set; }
    public int PageSize { get; set; }
}

public class Category
{
    public int Id { get; set; }
    public string Name { get; set; } = "";
    public string Color { get; set; } = "";
    public string Icon { get; set; } = "";
}

public class Platform
{
    public int Id { get; set; }
    public string Name { get; set; } = "";
    public string Icon { get; set; } = "";
}

public class Product
{
    public int Id { get; set; }
    public string Name { get; set; } = "";
}

public class PrivateMessage
{
    public int Id { get; set; }
    public int FromUserId { get; set; }
    public string FromUserName { get; set; } = "";
    public int ToUserId { get; set; }
    public string Content { get; set; } = "";
    public int IsRead { get; set; }
    public string CreateTime { get; set; } = "";
}

public class Announcement
{
    public int Id { get; set; }
    public string Title { get; set; } = "";
    public string Content { get; set; } = "";
    public string CreateTime { get; set; } = "";
}
