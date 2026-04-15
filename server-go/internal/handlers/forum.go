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
	category := c.Query("category")
	offset := (page - 1) * pageSize

	countQuery := "SELECT COUNT(*) FROM public_posts"
	query := `SELECT p.id, p.title, p.content, p.author_id, p.category, p.is_pinned, p.is_featured, 
		p.view_count, p.like_count, p.comment_count, p.create_time,
		u.username, u.real_name, u.avatar
		FROM public_posts p 
		LEFT JOIN users u ON p.author_id = u.id`
	args := []interface{}{}

	if category != "" && category != "全部" {
		query += " WHERE p.category = ?"
		countQuery += " WHERE category = ?"
		args = append(args, category)
	}

	query += " ORDER BY p.is_pinned DESC, p.create_time DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	var total int
	if len(args) > 2 {
		countArgs := args[:len(args)-2]
		_ = h.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	} else {
		_ = h.DB.QueryRow(countQuery).Scan(&total)
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{
				"list":     []map[string]interface{}{},
				"total":    0,
				"page":     page,
				"pageSize": pageSize,
			},
		})
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content, cat string
		var authorID sql.NullInt64
		var isPinned, isFeatured int
		var viewCount, likeCount, commentCount int
		var createTime string
		var username, realName, avatar sql.NullString
		err = rows.Scan(&id, &title, &content, &authorID, &cat, &isPinned, &isFeatured,
			&viewCount, &likeCount, &commentCount, &createTime,
			&username, &realName, &avatar)
		if err != nil {
			continue
		}

		posts = append(posts, map[string]interface{}{
			"id":          id,
			"title":       title,
			"content":     content,
			"author_id":   authorID.Int64,
			"username":    username.String,
			"real_name":   realName.String,
			"avatar":      avatar.String,
			"category":    cat,
			"is_pinned":   isPinned > 0,
			"is_featured": isFeatured > 0,
			"views":       viewCount,
			"likes":       likeCount,
			"comments":    commentCount,
			"created_at":  createTime,
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

func (h *ForumHandler) GetPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var title, content, category string
	var authorID sql.NullInt64
	var isPinned, isFeatured int
	var viewCount, likeCount, commentCount int
	var createTime string
	var username, realName, avatar sql.NullString

	err := h.DB.QueryRow(`SELECT p.id, p.title, p.content, p.author_id, p.category, p.is_pinned, p.is_featured,
		p.view_count, p.like_count, p.comment_count, p.create_time,
		u.username, u.real_name, u.avatar
		FROM public_posts p 
		LEFT JOIN users u ON p.author_id = u.id 
		WHERE p.id = ?`, id).Scan(
		&id, &title, &content, &authorID, &category, &isPinned, &isFeatured,
		&viewCount, &likeCount, &commentCount, &createTime,
		&username, &realName, &avatar)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	h.DB.Exec("UPDATE public_posts SET view_count = view_count + 1 WHERE id = ?", id)

	commentRows, err := h.DB.Query(`SELECT c.id, c.content, c.author_id, c.create_time,
		u.username, u.real_name, u.avatar
		FROM public_post_comments c 
		LEFT JOIN users u ON c.author_id = u.id 
		WHERE c.post_id = ?
		ORDER BY c.create_time ASC`, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": map[string]interface{}{
				"id":          id,
				"title":       title,
				"content":     content,
				"author_id":   authorID.Int64,
				"username":    username.String,
				"real_name":   realName.String,
				"avatar":      avatar.String,
				"category":    category,
				"is_pinned":   isPinned > 0,
				"is_featured": isFeatured > 0,
				"views":       viewCount + 1,
				"likes":       likeCount,
				"comments":    []map[string]interface{}{},
				"created_at":  createTime,
			},
		})
		return
	}
	defer commentRows.Close()

	var comments []map[string]interface{}
	for commentRows.Next() {
		var cID int
		var cContent string
		var cAuthorID sql.NullInt64
		var cCreatedAt string
		var cUsername, cRealName, cAvatar sql.NullString
		commentRows.Scan(&cID, &cContent, &cAuthorID, &cCreatedAt, &cUsername, &cRealName, &cAvatar)

		comments = append(comments, map[string]interface{}{
			"id":         cID,
			"content":    cContent,
			"author_id":  cAuthorID.Int64,
			"username":   cUsername.String,
			"real_name":  cRealName.String,
			"avatar":     cAvatar.String,
			"created_at": cCreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":          id,
			"title":       title,
			"content":     content,
			"author_id":   authorID.Int64,
			"username":    username.String,
			"real_name":   realName.String,
			"avatar":      avatar.String,
			"category":    category,
			"is_pinned":   isPinned > 0,
			"is_featured": isFeatured > 0,
			"views":       viewCount + 1,
			"likes":       likeCount,
			"comments":    comments,
			"created_at":  createTime,
		},
	})
}

func (h *ForumHandler) CreatePost(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" || req.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "标题和内容不能为空"})
		return
	}

	if req.Category == "" {
		req.Category = "讨论"
	}

	now := time.Now().Format(time.RFC3339)
	result, _ := h.DB.Exec(`INSERT INTO public_posts (title, content, author_id, category, create_time) 
		VALUES (?, ?, ?, ?, ?)`,
		req.Title, req.Content, userID, req.Category, now)
	id, _ := result.LastInsertId()

	var username, avatar sql.NullString
	h.DB.QueryRow("SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)

	database.AddLog("发帖", realName, realName, "公共论坛发帖【"+req.Title+"】")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":          id,
			"title":       req.Title,
			"content":     req.Content,
			"author_id":   userID,
			"username":    username.String,
			"real_name":   realName,
			"avatar":      avatar.String,
			"category":    req.Category,
			"is_pinned":   false,
			"is_featured": false,
			"views":       0,
			"likes":       0,
			"comments":    0,
			"created_at":  now,
		},
	})
}

func (h *ForumHandler) CreateComment(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "评论内容不能为空"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	result, _ := h.DB.Exec(`INSERT INTO public_post_comments (post_id, content, author_id, create_time) 
		VALUES (?, ?, ?, ?)`,
		postID, req.Content, userID, now)
	id, _ := result.LastInsertId()

	h.DB.Exec("UPDATE public_posts SET comment_count = comment_count + 1 WHERE id = ?", postID)

	var username, avatar sql.NullString
	h.DB.QueryRow("SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)

	var postAuthorID sql.NullInt64
	var postTitle string
	h.DB.QueryRow("SELECT author_id, title FROM public_posts WHERE id = ?", postID).Scan(&postAuthorID, &postTitle)

	if postAuthorID.Valid && postAuthorID.Int64 != int64(userID) {
		h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, post_id, from_user, create_time)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			postAuthorID.Int64, "post_comment", "💬 帖子收到新评论", "您的帖子【"+postTitle+"】收到了【"+realName+"】的评论",
			"normal", postID, realName, now)
	}

	database.AddLog("评论", realName, realName, "评论帖子【"+postTitle+"】")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"id":         id,
			"post_id":    postID,
			"content":    req.Content,
			"author_id":  userID,
			"username":   username.String,
			"real_name":  realName,
			"avatar":     avatar.String,
			"created_at": now,
		},
	})
}
func (h *ForumHandler) DeletePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var postAuthorID sql.NullInt64
	var title string
	h.DB.QueryRow("SELECT author_id, title FROM public_posts WHERE id = ?", postID).Scan(&postAuthorID, &title)

	if !postAuthorID.Valid {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	if role != "admin" && postAuthorID.Int64 != int64(userID) {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此帖子"})
		return
	}

	h.DB.Exec("DELETE FROM public_post_comments WHERE post_id = ?", postID)
	h.DB.Exec("DELETE FROM public_post_likes WHERE post_id = ?", postID)
	h.DB.Exec("DELETE FROM public_post_collects WHERE post_id = ?", postID)
	h.DB.Exec("DELETE FROM public_posts WHERE id = ?", postID)

	database.AddLog("删帖", "公共论坛帖子【"+title+"】", realName, "删除公共论坛帖子")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *ForumHandler) DeleteComment(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	commentID, _ := strconv.Atoi(c.Param("commentId"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var commentAuthorID sql.NullInt64
	h.DB.QueryRow("SELECT author_id FROM public_post_comments WHERE id = ?", commentID).Scan(&commentAuthorID)

	if !commentAuthorID.Valid {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "评论不存在"})
		return
	}

	if role != "admin" && commentAuthorID.Int64 != int64(userID) {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此评论"})
		return
	}

	h.DB.Exec("DELETE FROM public_post_comments WHERE id = ?", commentID)
	h.DB.Exec("UPDATE public_posts SET comment_count = CASE WHEN comment_count > 0 THEN comment_count - 1 ELSE 0 END WHERE id = ?", postID)

	database.AddLog("删除评论", "公共论坛评论", realName, "删除公共论坛评论")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *ForumHandler) LikePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var liked int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&liked)

	if liked > 0 {
		h.DB.Exec("DELETE FROM public_post_likes WHERE post_id = ? AND user_id = ?", postID, userID)
		h.DB.Exec("UPDATE public_posts SET like_count = like_count - 1 WHERE id = ?", postID)
	} else {
		h.DB.Exec("INSERT INTO public_post_likes (post_id, user_id) VALUES (?, ?)", postID, userID)
		h.DB.Exec("UPDATE public_posts SET like_count = like_count + 1 WHERE id = ?", postID)
	}

	var likeCount int
	h.DB.QueryRow("SELECT like_count FROM public_posts WHERE id = ?", postID).Scan(&likeCount)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"liked": liked == 0, "likes": likeCount}})
}

func (h *ForumHandler) CollectPost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ?", postID).Scan(&exists)
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
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var isPinned int
	h.DB.QueryRow("SELECT is_pinned FROM public_posts WHERE id = ?", postID).Scan(&isPinned)

	h.DB.Exec("UPDATE public_posts SET is_pinned = ? WHERE id = ?", 1-isPinned, postID)

	database.AddLog("置顶", "公共论坛帖子", realName, fmt.Sprintf("置顶/取消置顶公共论坛帖子"))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"is_pinned": 1 - isPinned}})
}

func (h *ForumHandler) FeaturePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_posts WHERE id = ?", postID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "帖子不存在"})
		return
	}

	var isFeatured int
	h.DB.QueryRow("SELECT is_featured FROM public_posts WHERE id = ?", postID).Scan(&isFeatured)

	h.DB.Exec("UPDATE public_posts SET is_featured = ? WHERE id = ?", 1-isFeatured, postID)

	database.AddLog("精选", "公共论坛帖子", realName, fmt.Sprintf("精选/取消精选公共论坛帖子"))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"is_featured": 1 - isFeatured}})
}

func (h *ForumHandler) SearchForumPosts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "搜索关键词不能为空"})
		return
	}

	searchPattern := "%" + keyword + "%"
	query := `SELECT p.id, p.title, p.content, p.author_id, p.category, p.is_pinned, p.is_featured, 
		p.view_count, p.like_count, p.comment_count, p.create_time,
		u.username, u.real_name, u.avatar
		FROM public_posts p 
		LEFT JOIN users u ON p.author_id = u.id
		WHERE (p.title LIKE ? OR p.content LIKE ?)
		ORDER BY p.create_time DESC`

	rows, err := h.DB.Query(query, searchPattern, searchPattern)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []interface{}{}})
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content, cat string
		var authorID sql.NullInt64
		var isPinned, isFeatured int
		var viewCount, likeCount, commentCount int
		var createTime string
		var username, realName, avatar sql.NullString
		err = rows.Scan(&id, &title, &content, &authorID, &cat, &isPinned, &isFeatured,
			&viewCount, &likeCount, &commentCount, &createTime,
			&username, &realName, &avatar)
		if err != nil {
			continue
		}

		posts = append(posts, map[string]interface{}{
			"id":          id,
			"title":       title,
			"content":     content,
			"author_id":   authorID.Int64,
			"username":    username.String,
			"real_name":   realName.String,
			"avatar":      avatar.String,
			"category":    cat,
			"is_pinned":   isPinned > 0,
			"is_featured": isFeatured > 0,
			"views":       viewCount,
			"likes":       likeCount,
			"comments":    commentCount,
			"created_at":  createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": posts})
}

func (h *ForumHandler) GetHotPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	query := `SELECT p.id, p.title, p.content, p.author_id, p.category, p.is_pinned, p.is_featured, 
		p.view_count, p.like_count, p.comment_count, p.create_time,
		u.username, u.real_name, u.avatar
		FROM public_posts p 
		LEFT JOIN users u ON p.author_id = u.id
		ORDER BY p.like_count DESC, p.view_count DESC LIMIT ?`

	rows, err := h.DB.Query(query, limit)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []interface{}{}})
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content, cat string
		var authorID sql.NullInt64
		var isPinned, isFeatured int
		var viewCount, likeCount, commentCount int
		var createTime string
		var username, realName, avatar sql.NullString
		err = rows.Scan(&id, &title, &content, &authorID, &cat, &isPinned, &isFeatured,
			&viewCount, &likeCount, &commentCount, &createTime,
			&username, &realName, &avatar)
		if err != nil {
			continue
		}

		posts = append(posts, map[string]interface{}{
			"id":          id,
			"title":       title,
			"content":     content,
			"author_id":   authorID.Int64,
			"username":    username.String,
			"real_name":   realName.String,
			"avatar":      avatar.String,
			"category":    cat,
			"is_pinned":   isPinned > 0,
			"is_featured": isFeatured > 0,
			"views":       viewCount,
			"likes":       likeCount,
			"comments":    commentCount,
			"created_at":  createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": posts})
}

func (h *ForumHandler) GetCollectedPosts(c *gin.Context) {
	userID, _ := c.Get("userID")

	query := `SELECT p.id, p.title, p.content, p.author_id, p.category, p.is_pinned, p.is_featured, 
		p.view_count, p.like_count, p.comment_count, p.create_time,
		u.username, u.real_name, u.avatar
		FROM public_posts p 
		INNER JOIN public_post_collects pc ON p.id = pc.post_id
		LEFT JOIN users u ON p.author_id = u.id
		WHERE pc.user_id = ?
		ORDER BY pc.create_time DESC`

	rows, err := h.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []interface{}{}})
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content, cat string
		var authorID sql.NullInt64
		var isPinned, isFeatured int
		var viewCount, likeCount, commentCount int
		var createTime string
		var username, realName, avatar sql.NullString
		err = rows.Scan(&id, &title, &content, &authorID, &cat, &isPinned, &isFeatured,
			&viewCount, &likeCount, &commentCount, &createTime,
			&username, &realName, &avatar)
		if err != nil {
			continue
		}

		posts = append(posts, map[string]interface{}{
			"id":          id,
			"title":       title,
			"content":     content,
			"author_id":   authorID.Int64,
			"username":    username.String,
			"real_name":   realName.String,
			"avatar":      avatar.String,
			"category":    cat,
			"is_pinned":   isPinned > 0,
			"is_featured": isFeatured > 0,
			"views":       viewCount,
			"likes":       likeCount,
			"comments":    commentCount,
			"created_at":  createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": posts})
}

func (h *ForumHandler) GetCollectStatus(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")

	var collected int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_post_collects WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&collected)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"collected": collected > 0}})
}

func (h *ForumHandler) GetLikeStatus(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("userID")

	var liked int
	h.DB.QueryRow("SELECT COUNT(*) FROM public_post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&liked)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"liked": liked > 0}})
}
