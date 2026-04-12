package models

type Blogger struct {
	ID               int    `json:"id"`
	Nickname         string `json:"nickname"`
	Category         string `json:"category"`
	Product          string `json:"product"`
	Contact          string `json:"contact"`
	Wechat           string `json:"wechat"`
	CustomContact    string `json:"custom_contact"`
	Platform         string `json:"platform"`
	PlatformAccount  string `json:"platform_account"`
	Description      string `json:"description"`
	Tags             string `json:"tags"`
	Avatar           string `json:"avatar"`
	UserName         string `json:"user_name"`
	RealName         string `json:"real_name"`
	Status           string `json:"status"`
	StatusUpdateTime string `json:"status_update_time"`
	StatusRemark     string `json:"status_remark"`
	IsDeleted        int    `json:"is_deleted"`
	IsInvalid        int    `json:"is_invalid"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`

	OwnerUsername string `json:"owner_username,omitempty"`
	OwnerRealName string `json:"owner_real_name,omitempty"`
	OwnerAvatar   string `json:"owner_avatar,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Bio           string `json:"bio,omitempty"`
	Role          string `json:"role,omitempty"`
	TeamID        *int64 `json:"team_id,omitempty"`
	TeamName      string `json:"team_name,omitempty"`
	TeamColor     string `json:"team_color,omitempty"`
}

type Category struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Icon       string `json:"icon"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"create_time"`
}

type Platform struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	CreateTime string `json:"create_time"`
}

type BloggerStatusHistory struct {
	ID         int    `json:"id"`
	BloggerID  int    `json:"blogger_id"`
	OldStatus  string `json:"old_status"`
	NewStatus  string `json:"new_status"`
	Operator   string `json:"operator"`
	Remark     string `json:"remark"`
	CreateTime string `json:"create_time"`
}

type OperationLog struct {
	ID         int    `json:"id"`
	Action     string `json:"action"`
	Target     string `json:"target"`
	Operator   string `json:"operator"`
	Detail     string `json:"detail"`
	CreateTime string `json:"create_time"`
}

type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"-"`
	RealName      string `json:"real_name"`
	Role          string `json:"role"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Bio           string `json:"bio"`
	Status        string `json:"status"`
	TeamID        *int64 `json:"team_id"`
	TeamName      string `json:"team_name"`
	TeamColor     string `json:"team_color"`
	ApprovedBy    string `json:"approved_by"`
	ApprovedTime  string `json:"approved_time"`
	IsDeleted     int    `json:"is_deleted"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`
}

type Team struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Logo          string `json:"logo"`
	BgImage       string `json:"bg_image"`
	Color         string `json:"color"`
	LeaderID      int    `json:"leader_id"`
	LeaderName    string `json:"leader_name"`
	AdminIDs      string `json:"admin_ids"`
	Announcement  string `json:"announcement"`
	IsPublic      int    `json:"is_public"`
	NeedApprove   int    `json:"need_approve"`
	MemberCount   int    `json:"member_count"`
	InviteCode    string `json:"invite_code"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`
}

type TeamPost struct {
	ID              int    `json:"id"`
	TeamID          int    `json:"team_id"`
	AuthorID        int    `json:"author_id"`
	AuthorName      string `json:"author_name"`
	AuthorUsername  string `json:"author_username"`
	AuthorAvatar    string `json:"author_avatar"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Category        string `json:"category"`
	CommentCount    int    `json:"comment_count"`
	LikeCount       int    `json:"like_count"`
	ViewCount       int    `json:"view_count"`
	IsPinned        int    `json:"is_pinned"`
	IsFeatured      int    `json:"is_featured"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`
}

type TeamPostComment struct {
	ID             int    `json:"id"`
	PostID         int    `json:"post_id"`
	AuthorID       int    `json:"author_id"`
	AuthorName     string `json:"author_name"`
	AuthorUsername string `json:"author_username"`
	AuthorAvatar   string `json:"author_avatar"`
	Content        string `json:"content"`
	CreateTime     string `json:"create_time"`
}

type PublicPost struct {
	ID            int    `json:"id"`
	AuthorID      int    `json:"author_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Category      string `json:"category"`
	IsPinned      int    `json:"is_pinned"`
	IsFeatured    int    `json:"is_featured"`
	ViewCount     int    `json:"view_count"`
	LikeCount     int    `json:"like_count"`
	CommentCount  int    `json:"comment_count"`
	CreateTime    string `json:"create_time"`
	AuthorName    string `json:"author_name,omitempty"`
	AuthorAvatar  string `json:"author_avatar,omitempty"`
	AuthorColor   string `json:"author_color,omitempty"`
	TeamName      string `json:"team_name,omitempty"`
}

type PublicPostComment struct {
	ID           int    `json:"id"`
	PostID       int    `json:"post_id"`
	AuthorID     int    `json:"author_id"`
	AuthorName   string `json:"author_name"`
	AuthorAvatar string `json:"author_avatar"`
	Content      string `json:"content"`
	CreateTime   string `json:"create_time"`
}

type Notification struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Priority   string `json:"priority"`
	TeamID     *int64 `json:"team_id"`
	FromUser   string `json:"from_user"`
	BloggerID  *int64 `json:"blogger_id"`
	PostID     *int64 `json:"post_id"`
	IsRead     int    `json:"is_read"`
	CreateTime string `json:"create_time"`
	IsTeamMessage bool `json:"is_team_message,omitempty"`
}

type PrivateMessage struct {
	ID            int    `json:"id"`
	FromUserID    int    `json:"from_user_id"`
	FromUserName  string `json:"from_user_name"`
	ToUserID      int    `json:"to_user_id"`
	Content       string `json:"content"`
	IsRead        int    `json:"is_read"`
	CreateTime    string `json:"create_time"`
}

type TeamMessage struct {
	ID           int    `json:"id"`
	TeamID       int    `json:"team_id"`
	SenderID     int    `json:"sender_id"`
	SenderName   string `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Type         string `json:"type"`
	IsRead       int    `json:"is_read"`
	CreateTime   string `json:"create_time"`
}

type Followup struct {
	ID         int    `json:"id"`
	BloggerID  int    `json:"blogger_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Operator   string `json:"operator"`
	CreateTime string `json:"create_time"`
}

type Cooperation struct {
	ID         int     `json:"id"`
	BloggerID  int     `json:"blogger_id"`
	Title      string  `json:"title"`
	Date       string  `json:"date"`
	Status     string  `json:"status"`
	Product    string  `json:"product"`
	Amount     float64 `json:"amount"`
	Note       string  `json:"note"`
	CreateTime string  `json:"create_time"`
	UpdateTime string  `json:"update_time"`
}

type BloggerTransferRequest struct {
	ID              int    `json:"id"`
	BloggerID       int    `json:"blogger_id"`
	BloggerName     string `json:"blogger_name"`
	FromUserID      int    `json:"from_user_id"`
	FromUserName    string `json:"from_user_name"`
	ToUserID        int    `json:"to_user_id"`
	ToUserName      string `json:"to_user_name"`
	Reason          string `json:"reason"`
	Status          string `json:"status"`
	FromConfirmed   int    `json:"from_confirmed"`
	ToConfirmed     int    `json:"to_confirmed"`
	AdminConfirmed  int    `json:"admin_confirmed"`
	CreateTime      string `json:"create_time"`
}

type SystemSetting struct {
	ID         int    `json:"id"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type Snapshot struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Filename   string `json:"filename"`
	Size       int    `json:"size"`
	Created    string `json:"created"`
	Mtime      string `json:"mtime"`
	CreateTime string `json:"create_time"`
}

type Announcement struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

type SavedFilter struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Filters   string `json:"filters"`
	IsDefault int    `json:"is_default"`
	CreateTime string `json:"create_time"`
}

type WorkflowRule struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	TriggerType  string `json:"trigger_type"`
	TriggerValue string `json:"trigger_value"`
	ActionType   string `json:"action_type"`
	ActionConfig string `json:"action_config"`
	Enabled      int    `json:"enabled"`
	Priority     int    `json:"priority"`
	CreateTime   string `json:"create_time"`
}
