package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"server-go/internal/database"
	"server-go/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	DB *sql.DB
}

func NewTeamHandler(db *sql.DB) *TeamHandler {
	return &TeamHandler{DB: db}
}

func (h *TeamHandler) GetTeams(c *gin.Context) {
	rows, _ := h.DB.Query(`SELECT id, name, color, logo, bg_image, description, announcement, leader_id, create_time FROM teams ORDER BY id`)
	var teams []map[string]interface{}
	for rows.Next() {
		var id int
		var name, color, createTime string
		var logo, bgImage, description, announcement sql.NullString
		var leaderID sql.NullInt64
		rows.Scan(&id, &name, &color, &logo, &bgImage, &description, &announcement, &leaderID, &createTime)

		var memberCount int
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE team_id = ? AND is_deleted = 0", id).Scan(&memberCount)

		teams = append(teams, map[string]interface{}{
			"id":           id,
			"name":         name,
			"color":        color,
			"logo":         logo.String,
			"bg_image":     bgImage.String,
			"description":  description.String,
			"announcement": announcement.String,
			"leader_id":    leaderID.Int64,
			"create_time":  createTime,
			"member_count": memberCount,
		})
	}
	rows.Close()

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": teams})
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		Name         string `json:"name"`
		Color        string `json:"color"`
		Description  string `json:"description"`
		Logo         string `json:"logo"`
		BgImage      string `json:"bg_image"`
		Announcement string `json:"announcement"`
		LeaderID     *int64 `json:"leader_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "团队名称不能为空"})
		return
	}

	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM teams WHERE name = ?", req.Name).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "团队名称已存在"})
		return
	}

	color := req.Color
	if color == "" {
		color = "#3b82f6"
	}

	result, _ := h.DB.Exec(`INSERT INTO teams (name, color, description, logo, bg_image, announcement, leader_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.Name, color, req.Description, req.Logo, req.BgImage, req.Announcement, req.LeaderID)
	id, _ := result.LastInsertId()

	adminRealName, _ := c.Get("realName")
	database.AddLog("新增", "团队【"+req.Name+"】", adminRealName.(string), "新增团队【"+req.Name+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"id": id, "name": req.Name, "color": color}})
}

func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Name         string `json:"name"`
		Color        string `json:"color"`
		Description  string `json:"description"`
		Logo         string `json:"logo"`
		BgImage      string `json:"bg_image"`
		Announcement string `json:"announcement"`
		LeaderID     *int64 `json:"leader_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var oldName string
	h.DB.QueryRow("SELECT name FROM teams WHERE id = ?", id).Scan(&oldName)
	if oldName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "团队不存在"})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, req.Name)
	}
	if req.Color != "" {
		updates = append(updates, "color = ?")
		args = append(args, req.Color)
	}
	if req.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, req.Description)
	}
	if req.Logo != "" {
		updates = append(updates, "logo = ?")
		args = append(args, req.Logo)
	}
	if req.BgImage != "" {
		updates = append(updates, "bg_image = ?")
		args = append(args, req.BgImage)
	}
	if req.Announcement != "" {
		updates = append(updates, "announcement = ?")
		args = append(args, req.Announcement)
	}
	if req.LeaderID != nil {
		updates = append(updates, "leader_id = ?")
		args = append(args, req.LeaderID)
	}

	if len(updates) > 0 {
		updates = append(updates, "update_time = CURRENT_TIMESTAMP")
		args = append(args, id)
		query := "UPDATE teams SET " + string(join(updates, ", ")) + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	adminRealName, _ := c.Get("realName")
	newName := req.Name
	if newName == "" {
		newName = oldName
	}
	database.AddLog("编辑", "团队【"+oldName+"】", adminRealName.(string), "编辑团队【"+oldName+"】→【"+newName+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var name string
	h.DB.QueryRow("SELECT name FROM teams WHERE id = ?", id).Scan(&name)
	if name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "团队不存在"})
		return
	}

	var memberCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE team_id = ? AND is_deleted = 0", id).Scan(&memberCount)
	if memberCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "团队下还有成员，无法删除"})
		return
	}

	h.DB.Exec("DELETE FROM teams WHERE id = ?", id)

	adminRealName, _ := c.Get("realName")
	database.AddLog("删除", "团队【"+name+"】", adminRealName.(string), "删除团队【"+name+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *TeamHandler) GetTeamPosts(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	var teamID interface{}
	if teamIDStr := c.Param("id"); teamIDStr != "" {
		if tid, err := strconv.ParseInt(teamIDStr, 10, 64); err == nil {
			teamID = tid
		}
	} else if role != "admin" {
		var tid sql.NullInt64
		h.DB.QueryRow("SELECT team_id FROM users WHERE id = ?", userID).Scan(&tid)
		if !tid.Valid {
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": []interface{}{}, "total": 0}})
			return
		}
		teamID = tid.Int64
	}

	countQuery := "SELECT COUNT(*) FROM team_posts WHERE 1=1"
	query := `SELECT p.id, p.title, p.content, p.category, p.author_id, p.create_time, 
		u.username, u.real_name, u.avatar
		FROM team_posts p 
		LEFT JOIN users u ON p.author_id = u.id 
		WHERE 1=1`
	args := []interface{}{}

	if teamID != nil {
		query += " AND p.team_id = ?"
		countQuery += " AND team_id = ?"
		args = append(args, teamID)
	}

	query += " ORDER BY p.create_time DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	var total int
	if teamID != nil {
		h.DB.QueryRow(countQuery, teamID).Scan(&total)
	} else {
		h.DB.QueryRow(countQuery).Scan(&total)
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": []interface{}{}, "total": 0}})
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content, category string
		var authorID sql.NullInt64
		var createTime string
		var username, realName, avatar sql.NullString
		rows.Scan(&id, &title, &content, &category, &authorID, &createTime, &username, &realName, &avatar)

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"title":      title,
			"content":    content,
			"category":   category,
			"user_id":    authorID.Int64,
			"username":   username.String,
			"real_name":  realName.String,
			"avatar":     avatar.String,
			"created_at": createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     posts,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h *TeamHandler) CreateTeamPost(c *gin.Context) {
	userID, _ := c.Get("userID")
	realName, _ := c.Get("realName")
	role, _ := c.Get("role")

	var req struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" || req.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "标题和内容不能为空"})
		return
	}

	var targetTeamID int64 = 0
	if teamIDStr := c.Param("id"); teamIDStr != "" {
		if tid, err := strconv.ParseInt(teamIDStr, 10, 64); err == nil {
			targetTeamID = tid
		}
	}

	var teamID sql.NullInt64
	if targetTeamID > 0 {
		teamID.Valid = true
		teamID.Int64 = targetTeamID
	} else {
		h.DB.QueryRow("SELECT team_id FROM users WHERE id = ?", userID).Scan(&teamID)
		if !teamID.Valid {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "您不在任何团队中"})
			return
		}
	}

	if role != "admin" && targetTeamID > 0 {
		var userTeamID sql.NullInt64
		h.DB.QueryRow("SELECT team_id FROM users WHERE id = ?", userID).Scan(&userTeamID)
		if !userTeamID.Valid || userTeamID.Int64 != teamID.Int64 {
			c.JSON(http.StatusOK, gin.H{"code": 403, "message": "您不在该团队中"})
			return
		}
	}

	tagsStr := ""
	if len(req.Tags) > 0 {
		tagsStr = stringsJoin(req.Tags, ",")
	}

	result, _ := h.DB.Exec(`INSERT INTO team_posts (title, content, tags, user_id, team_id, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		req.Title, req.Content, tagsStr, userID, teamID.Int64, time.Now().Format(time.RFC3339))
	id, _ := result.LastInsertId()

	var username, avatar sql.NullString
	h.DB.QueryRow("SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)

	rows, _ := h.DB.Query("SELECT id, username, real_name FROM users WHERE team_id = ? AND id != ? AND is_deleted = 0", teamID.Int64, userID)
	for rows.Next() {
		var uID int
		var uName, uRealName string
		rows.Scan(&uID, &uName, &uRealName)
		h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, post_id, from_user, create_time)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			uID, "team_post", "📢 团队新消息", "【"+realName.(string)+"】发布了新话题："+req.Title,
			"normal", id, realName, time.Now().Format(time.RFC3339))
	}
	rows.Close()

	database.AddLog("发帖", realName.(string), realName.(string), "团队发帖【"+req.Title+"】")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":         id,
			"title":      req.Title,
			"content":    req.Content,
			"tags":       tagsStr,
			"user_id":    userID,
			"username":   username.String,
			"real_name":  realName,
			"avatar":     avatar.String,
			"created_at": time.Now().Format(time.RFC3339),
		},
	})
}

func stringsJoin(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += sep + slice[i]
	}
	return result
}

func (h *TeamHandler) GetTeamMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var teamID interface{}
	if teamIDStr := c.Param("id"); teamIDStr != "" {
		if tid, err := strconv.ParseInt(teamIDStr, 10, 64); err == nil {
			teamID = tid
		}
	} else if role != "admin" {
		var tid sql.NullInt64
		h.DB.QueryRow("SELECT team_id FROM users WHERE id = ?", userID).Scan(&tid)
		if !tid.Valid {
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.TeamMessage{}})
			return
		}
		teamID = tid.Int64
	}

	var query string
	var args []interface{}

	if teamID != nil {
		query = `SELECT id, team_id, sender_id, sender_name, sender_avatar, title, content, type, is_read, create_time 
			FROM team_messages WHERE team_id = ? ORDER BY create_time DESC`
		args = append(args, teamID)
	} else {
		query = `SELECT id, team_id, sender_id, sender_name, sender_avatar, title, content, type, is_read, create_time 
			FROM team_messages ORDER BY create_time DESC`
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.TeamMessage{}})
		return
	}
	defer rows.Close()

	var messages []models.TeamMessage
	for rows.Next() {
		var m models.TeamMessage
		rows.Scan(&m.ID, &m.TeamID, &m.SenderID, &m.SenderName, &m.SenderAvatar, &m.Title, &m.Content, &m.Type, &m.IsRead, &m.CreateTime)
		messages = append(messages, m)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": messages})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil || teamID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "团队不存在"})
		return
	}

	var id int
	var name, color, createTime string
	var logo, bgImage, description, announcement sql.NullString
	var leaderID sql.NullInt64

	err = h.DB.QueryRow(`SELECT id, name, color, logo, bg_image, description, announcement, leader_id, create_time 
		FROM teams WHERE id = ?`, teamID).Scan(
		&id, &name, &color, &logo, &bgImage, &description, &announcement, &leaderID, &createTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "团队不存在"})
		return
	}

	var memberCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE team_id = ? AND is_deleted = 0", id).Scan(&memberCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":           id,
			"name":         name,
			"color":        color,
			"logo":         logo.String,
			"bg_image":     bgImage.String,
			"description":  description.String,
			"announcement": announcement.String,
			"leader_id":    leaderID.Int64,
			"create_time":  createTime,
			"member_count": memberCount,
		},
	})
}

func (h *TeamHandler) CreateTeamMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	realName, _ := c.Get("realName")
	role, _ := c.Get("role")

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "标题和内容不能为空"})
		return
	}

	var teamID sql.NullInt64
	if role != "admin" {
		h.DB.QueryRow("SELECT team_id FROM users WHERE id = ?", userID).Scan(&teamID)
		if !teamID.Valid {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "您不在任何团队中"})
			return
		}
	}

	msgType := req.Type
	if msgType == "" {
		msgType = "normal"
	}

	var senderUsername, senderAvatar sql.NullString
	h.DB.QueryRow("SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&senderUsername, &senderAvatar)

	now := time.Now().Format(time.RFC3339)
	result, _ := h.DB.Exec(`INSERT INTO team_messages 
		(team_id, sender_id, sender_name, sender_avatar, title, content, type, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		teamID.Int64, userID, realName, senderAvatar.String, req.Title, req.Content, msgType, now)

	msgID, _ := result.LastInsertId()

	if teamID.Valid {
		rows, _ := h.DB.Query("SELECT id, username FROM users WHERE team_id = ? AND id != ? AND is_deleted = 0", teamID.Int64, userID)
		for rows.Next() {
			var uID int
			var uName string
			rows.Scan(&uID, &uName)
			h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, team_id, from_user, create_time)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				uID, "team_message", "📢 团队消息", "【"+realName.(string)+"】发布了团队消息："+req.Title,
				"normal", teamID.Int64, realName, now)
		}
		rows.Close()
	}

	database.AddLog("团队消息", "消息【"+req.Title+"】", realName.(string), "发布团队消息")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送成功",
		"data":    gin.H{"id": msgID},
	})
}

func (h *TeamHandler) GetTeamPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var p models.TeamPost
	err := h.DB.QueryRow(`SELECT p.id, p.team_id, p.author_id, p.author_name, p.author_username, p.author_avatar, 
		p.title, p.content, p.category, p.comment_count, p.like_count, p.view_count, p.is_pinned, p.is_featured, p.create_time, p.update_time
		FROM team_posts p WHERE p.id = ?`, postID).Scan(
		&p.ID, &p.TeamID, &p.AuthorID, &p.AuthorName, &p.AuthorUsername, &p.AuthorAvatar,
		&p.Title, &p.Content, &p.Category, &p.CommentCount, &p.LikeCount, &p.ViewCount, &p.IsPinned, &p.IsFeatured, &p.CreateTime, &p.UpdateTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var teamID sql.NullInt64
	if role != "admin" {
		h.DB.QueryRow("SELECT team_id FROM users WHERE id = ?", userID).Scan(&teamID)
		if !teamID.Valid || teamID.Int64 != int64(p.TeamID) {
			c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限查看此帖子"})
			return
		}
	}

	h.DB.Exec("UPDATE team_posts SET view_count = view_count + 1 WHERE id = ?", postID)

	var liked, collected int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&liked)
	h.DB.QueryRow("SELECT COUNT(*) FROM team_post_collects WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&collected)

	commentRows, _ := h.DB.Query(`SELECT id, post_id, author_id, author_name, author_username, author_avatar, content, create_time 
		FROM team_post_comments WHERE post_id = ? ORDER BY create_time ASC`, postID)

	var comments []models.TeamPostComment
	for commentRows.Next() {
		var c models.TeamPostComment
		commentRows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.AuthorName, &c.AuthorUsername, &c.AuthorAvatar, &c.Content, &c.CreateTime)
		comments = append(comments, c)
	}
	commentRows.Close()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"post":       p,
			"liked":      liked > 0,
			"collected":  collected > 0,
			"comments":   comments,
			"view_count": p.ViewCount + 1,
		},
	})
}

func (h *TeamHandler) LikeTeamPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var liked int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&liked)

	if liked > 0 {
		h.DB.Exec("DELETE FROM team_post_likes WHERE post_id = ? AND user_id = ?", postID, userID)
		h.DB.Exec("UPDATE team_posts SET like_count = like_count - 1 WHERE id = ?", postID)
	} else {
		h.DB.Exec("INSERT INTO team_post_likes (post_id, user_id) VALUES (?, ?)", postID, userID)
		h.DB.Exec("UPDATE team_posts SET like_count = like_count + 1 WHERE id = ?", postID)
	}

	var likeCount int
	h.DB.QueryRow("SELECT like_count FROM team_posts WHERE id = ?", postID).Scan(&likeCount)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"liked": liked == 0, "like_count": likeCount}})
}

func (h *TeamHandler) CollectTeamPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var collected int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_post_collects WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&collected)

	if collected > 0 {
		h.DB.Exec("DELETE FROM team_post_collects WHERE post_id = ? AND user_id = ?", postID, userID)
	} else {
		h.DB.Exec("INSERT INTO team_post_collects (post_id, user_id) VALUES (?, ?)", postID, userID)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"collected": collected == 0}})
}

func (h *TeamHandler) CreateTeamPostComment(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")
	realName, _ := c.Get("realName")

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	var username, avatar sql.NullString
	h.DB.QueryRow("SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)

	now := time.Now().Format(time.RFC3339)
	result, _ := h.DB.Exec(`INSERT INTO team_post_comments 
		(post_id, author_id, author_name, author_username, author_avatar, content, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		postID, userID, realName, username.String, avatar.String, req.Content, now)

	commentID, _ := result.LastInsertId()

	h.DB.Exec("UPDATE team_posts SET comment_count = comment_count + 1 WHERE id = ?", postID)

	var postAuthorID int
	var postTitle string
	h.DB.QueryRow("SELECT author_id, title FROM team_posts WHERE id = ?", postID).Scan(&postAuthorID, &postTitle)

	if postAuthorID != userID.(int) {
		h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, post_id, from_user, create_time)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			postAuthorID, "post_comment", "💬 团队帖子收到新评论",
			"您的团队帖子【"+postTitle+"】收到了【"+realName.(string)+"】的评论",
			"normal", postID, realName, now)
	}

	database.AddLog("评论", "团队帖子评论", realName.(string), "评论团队帖子")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":          commentID,
			"content":     req.Content,
			"author_id":   userID,
			"author_name": realName,
			"avatar":      avatar.String,
			"created_at":  now,
		},
	})
}

func (h *TeamHandler) DeleteTeamPostComment(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	commentID, _ := strconv.Atoi(c.Param("commentId"))
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	realName, _ := c.Get("realName")

	var authorID int
	h.DB.QueryRow("SELECT author_id FROM team_post_comments WHERE id = ?", commentID).Scan(&authorID)
	if authorID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "评论不存在"})
		return
	}

	if role != "admin" && authorID != userID.(int) {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此评论"})
		return
	}

	h.DB.Exec("DELETE FROM team_post_comments WHERE id = ?", commentID)
	h.DB.Exec("UPDATE team_posts SET comment_count = comment_count - 1 WHERE id = ?", postID)

	database.AddLog("删除评论", "团队帖子评论", realName.(string), "删除团队帖子评论")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *TeamHandler) PinTeamPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	role, _ := c.Get("role")
	realName, _ := c.Get("realName")

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var isPinned int
	h.DB.QueryRow("SELECT is_pinned FROM team_posts WHERE id = ?", postID).Scan(&isPinned)

	h.DB.Exec("UPDATE team_posts SET is_pinned = ? WHERE id = ?", 1-isPinned, postID)

	database.AddLog("置顶", "团队帖子", realName.(string), fmt.Sprintf("置顶/取消置顶团队帖子"))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"is_pinned": 1 - isPinned}})
}

func (h *TeamHandler) FeatureTeamPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	role, _ := c.Get("role")
	realName, _ := c.Get("realName")

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM team_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var isFeatured int
	h.DB.QueryRow("SELECT is_featured FROM team_posts WHERE id = ?", postID).Scan(&isFeatured)

	h.DB.Exec("UPDATE team_posts SET is_featured = ? WHERE id = ?", 1-isFeatured, postID)

	database.AddLog("精选", "团队帖子", realName.(string), fmt.Sprintf("精选/取消精选团队帖子"))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"is_featured": 1 - isFeatured}})
}

func (h *TeamHandler) DeleteTeamPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	realName, _ := c.Get("realName")

	var authorID int
	var title string
	h.DB.QueryRow("SELECT author_id, title FROM team_posts WHERE id = ?", postID).Scan(&authorID, &title)
	if authorID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	if role != "admin" && authorID != userID.(int) {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此帖子"})
		return
	}

	h.DB.Exec("DELETE FROM team_posts WHERE id = ?", postID)
	h.DB.Exec("DELETE FROM team_post_comments WHERE post_id = ?", postID)
	h.DB.Exec("DELETE FROM team_post_likes WHERE post_id = ?", postID)
	h.DB.Exec("DELETE FROM team_post_collects WHERE post_id = ?", postID)

	database.AddLog("删帖", "团队帖子【"+title+"】", realName.(string), "删除团队帖子")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *TeamHandler) GetMyTeam(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var teamID sql.NullInt64
	if role == "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": nil})
		return
	}

	h.DB.QueryRow("SELECT team_id FROM users WHERE id = ?", userID).Scan(&teamID)
	if !teamID.Valid {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": nil})
		return
	}

	var id int
	var name, color, createTime string
	err := h.DB.QueryRow("SELECT id, name, color, create_time FROM teams WHERE id = ?", teamID.Int64).Scan(&id, &name, &color, &createTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": nil})
		return
	}

	var memberCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE team_id = ? AND is_deleted = 0", id).Scan(&memberCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":           id,
			"name":         name,
			"color":        color,
			"create_time":  createTime,
			"member_count": memberCount,
		},
	})
}

func (h *TeamHandler) UpdateMyTeam(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		TeamID int64 `json:"team_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的请求"})
		return
	}

	_, err := h.DB.Exec("UPDATE users SET team_id = ? WHERE id = ?", req.TeamID, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}
