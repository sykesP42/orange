package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"server-go/internal/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ForumHandler struct {
	DB *sql.DB
}

func NewForumHandler(db *sql.DB) *ForumHandler {
	return &ForumHandler{DB: db}
}

func (h *ForumHandler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	var total int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE is_deleted = 0").Scan(&total)

	rows, _ := h.DB.Query(`SELECT p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, p.views, p.likes,
		u.username, u.real_name, u.avatar
		FROM public_posts p 
		JOIN users u ON p.user_id = u.id 
		WHERE p.is_deleted = 0
		ORDER BY p.created_at DESC LIMIT ? OFFSET ?`, pageSize, offset)

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content string
		var userID int
		var createdAt, updatedAt string
		var views, likes int
		var username, realName, avatar sql.NullString
		rows.Scan(&id, &title, &content, &userID, &createdAt, &updatedAt, &views, &likes,
			&username, &realName, &avatar)

		var commentCount int
		h.DB.QueryRow("SELECT COUNT(*) FROM public_comments WHERE post_id = ? AND is_deleted = 0", id).Scan(&commentCount)

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"title":      title,
			"content":    content,
			"user_id":    userID,
			"username":   username.String,
			"real_name":  realName.String,
			"avatar":     avatar.String,
			"created_at": createdAt,
			"updated_at": updatedAt,
			"views":      views,
			"likes":      likes,
			"comments":   commentCount,
		})
	}
	rows.Close()

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

func (h *ForumHandler) GetPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var title, content string
	var userID int
	var createdAt, updatedAt string
	var views, likes int
	var username, realName, avatar sql.NullString

	err := h.DB.QueryRow(`SELECT p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, p.views, p.likes,
		u.username, u.real_name, u.avatar
		FROM public_posts p 
		JOIN users u ON p.user_id = u.id 
		WHERE p.id = ? AND p.is_deleted = 0`, id).Scan(
		&id, &title, &content, &userID, &createdAt, &updatedAt, &views, &likes,
		&username, &realName, &avatar)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	h.DB.Exec("UPDATE public_posts SET views = views + 1 WHERE id = ?", id)

	var commentCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_comments WHERE post_id = ? AND is_deleted = 0", id).Scan(&commentCount)

	commentRows, _ := h.DB.Query(`SELECT c.id, c.content, c.user_id, c.created_at,
		u.username, u.real_name, u.avatar
		FROM public_comments c 
		JOIN users u ON c.user_id = u.id 
		WHERE c.post_id = ? AND c.is_deleted = 0
		ORDER BY c.created_at ASC`, id)

	var comments []map[string]interface{}
	for commentRows.Next() {
		var cID int
		var cContent string
		var cUserID int
		var cCreatedAt string
		var cUsername, cRealName, cAvatar sql.NullString
		commentRows.Scan(&cID, &cContent, &cUserID, &cCreatedAt, &cUsername, &cRealName, &cAvatar)

		comments = append(comments, map[string]interface{}{
			"id":         cID,
			"content":    cContent,
			"user_id":    cUserID,
			"username":   cUsername.String,
			"real_name":  cRealName.String,
			"avatar":     cAvatar.String,
			"created_at": cCreatedAt,
		})
	}
	commentRows.Close()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":         id,
			"title":      title,
			"content":    content,
			"user_id":    userID,
			"username":   username.String,
			"real_name":  realName.String,
			"avatar":     avatar.String,
			"created_at": createdAt,
			"updated_at": updatedAt,
			"views":      views + 1,
			"likes":      likes,
			"comments":   comments,
		},
	})
}

func (h *ForumHandler) CreatePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	realName, _ := c.Get("realName")

	var req struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" || req.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "标题和内容不能为空"})
		return
	}

	tagsStr := ""
	if len(req.Tags) > 0 {
		tagsStr = stringsJoin(req.Tags, ",")
	}

	result, _ := h.DB.Exec(`INSERT INTO public_posts (title, content, tags, user_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		req.Title, req.Content, tagsStr, userID, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339))
	id, _ := result.LastInsertId()

	var username, avatar sql.NullString
	h.DB.QueryRow("SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)

	database.AddLog("发帖", realName.(string), realName.(string), "公共论坛发帖【"+req.Title+"】")

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
			"updated_at": time.Now().Format(time.RFC3339),
			"views":      0,
			"likes":      0,
		},
	})
}

func (h *ForumHandler) CreateComment(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")
	realName, _ := c.Get("realName")

	var req struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "评论内容不能为空"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ? AND is_deleted = 0", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	result, _ := h.DB.Exec(`INSERT INTO public_comments (post_id, content, user_id, created_at) 
		VALUES (?, ?, ?, ?)`,
		postID, req.Content, userID, time.Now().Format(time.RFC3339))
	id, _ := result.LastInsertId()

	var username, avatar sql.NullString
	h.DB.QueryRow("SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)

	var postUserID int
	var postTitle string
	h.DB.QueryRow("SELECT user_id, title FROM public_posts WHERE id = ?", postID).Scan(&postUserID, &postTitle)

	if postUserID != userID {
		h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, post_id, from_user, create_time)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			postUserID, "post_comment", "💬 帖子收到新评论", "您的帖子【"+postTitle+"】收到了【"+realName.(string)+"】的评论",
			"normal", postID, realName, time.Now().Format(time.RFC3339))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":         id,
			"post_id":    postID,
			"content":    req.Content,
			"user_id":    userID,
			"username":   username.String,
			"real_name":  realName,
			"avatar":     avatar.String,
			"created_at": time.Now().Format(time.RFC3339),
		},
	})
}

func (h *ForumHandler) DeletePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var postUserID int
	var title string
	h.DB.QueryRow("SELECT user_id, title FROM public_posts WHERE id = ? AND is_deleted = 0", postID).Scan(&postUserID, &title)

	if postUserID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	if role != "admin" && postUserID != userID {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此帖子"})
		return
	}

	h.DB.Exec("UPDATE public_posts SET is_deleted = 1 WHERE id = ?", postID)

	realName, _ := c.Get("realName")
	database.AddLog("删帖", realName.(string), realName.(string), "删除帖子【"+title+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *ForumHandler) DeleteComment(c *gin.Context) {
	commentID, _ := strconv.Atoi(c.Param("commentId"))
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	realName, _ := c.Get("realName")

	var commentUserID int
	h.DB.QueryRow("SELECT user_id FROM public_comments WHERE id = ? AND is_deleted = 0", commentID).Scan(&commentUserID)

	if commentUserID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "评论不存在"})
		return
	}

	if role != "admin" && commentUserID != userID {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此评论"})
		return
	}

	h.DB.Exec("UPDATE public_comments SET is_deleted = 1 WHERE id = ?", commentID)
	h.DB.Exec("UPDATE public_posts SET comment_count = comment_count - 1 WHERE id = ?", postID)

	database.AddLog("删除评论", "公共论坛评论", realName.(string), "删除公共论坛评论")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *ForumHandler) LikePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ? AND is_deleted = 0", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var liked int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&liked)

	if liked > 0 {
		h.DB.Exec("DELETE FROM public_post_likes WHERE post_id = ? AND user_id = ?", postID, userID)
		h.DB.Exec("UPDATE public_posts SET likes = likes - 1 WHERE id = ?", postID)
	} else {
		h.DB.Exec("INSERT INTO public_post_likes (post_id, user_id) VALUES (?, ?)", postID, userID)
		h.DB.Exec("UPDATE public_posts SET likes = likes + 1 WHERE id = ?", postID)
	}

	var likes int
	h.DB.QueryRow("SELECT likes FROM public_posts WHERE id = ?", postID).Scan(&likes)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"liked": liked == 0, "likes": likes}})
}

func (h *ForumHandler) CollectPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ? AND is_deleted = 0", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var collected int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_post_collects WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&collected)

	if collected > 0 {
		h.DB.Exec("DELETE FROM public_post_collects WHERE post_id = ? AND user_id = ?", postID, userID)
	} else {
		h.DB.Exec("INSERT INTO public_post_collects (post_id, user_id) VALUES (?, ?)", postID, userID)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"collected": collected == 0}})
}

func (h *ForumHandler) PinPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	role, _ := c.Get("role")
	realName, _ := c.Get("realName")

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ? AND is_deleted = 0", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var isPinned int
	h.DB.QueryRow("SELECT is_pinned FROM public_posts WHERE id = ?", postID).Scan(&isPinned)

	h.DB.Exec("UPDATE public_posts SET is_pinned = ? WHERE id = ?", 1-isPinned, postID)

	database.AddLog("置顶", "公共论坛帖子", realName.(string), fmt.Sprintf("置顶/取消置顶公共论坛帖子"))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"is_pinned": 1 - isPinned}})
}

func (h *ForumHandler) FeaturePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	role, _ := c.Get("role")
	realName, _ := c.Get("realName")

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ? AND is_deleted = 0", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var isFeatured int
	h.DB.QueryRow("SELECT is_featured FROM public_posts WHERE id = ?", postID).Scan(&isFeatured)

	h.DB.Exec("UPDATE public_posts SET is_featured = ? WHERE id = ?", 1-isFeatured, postID)

	database.AddLog("精选", "公共论坛帖子", realName.(string), fmt.Sprintf("精选/取消精选公共论坛帖子"))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"is_featured": 1 - isFeatured}})
}
