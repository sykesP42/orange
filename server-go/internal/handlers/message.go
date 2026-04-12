package handlers

import (
	"database/sql"
	"net/http"
	"server-go/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	DB *sql.DB
}

func NewMessageHandler(db *sql.DB) *MessageHandler {
	return &MessageHandler{DB: db}
}

type Conversation struct {
	UserID         int                    `json:"user_id"`
	Username       string                 `json:"username"`
	RealName       string                 `json:"real_name"`
	Avatar         string                 `json:"avatar"`
	LastMessage    string                 `json:"last_message"`
	LastMessageTime string                `json:"last_message_time"`
	UnreadCount    int                    `json:"unread_count"`
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	query := `
		SELECT 
			u.id, u.username, u.real_name, u.avatar,
			pm.content as last_message,
			pm.create_time as last_message_time,
			(SELECT COUNT(*) FROM private_messages 
			 WHERE (from_user_id = u.id AND to_user_id = ? AND is_read = 0)) as unread_count
		FROM (
			SELECT DISTINCT 
				CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END as user_id
			FROM private_messages
			WHERE from_user_id = ? OR to_user_id = ?
		) as conv
		JOIN users u ON conv.user_id = u.id
		LEFT JOIN private_messages pm ON 
			(pm.from_user_id = u.id AND pm.to_user_id = ?) OR 
			(pm.from_user_id = ? AND pm.to_user_id = u.id)
		WHERE pm.id IN (
			SELECT MAX(id) FROM private_messages 
			WHERE (from_user_id = u.id AND to_user_id = ?) OR 
			      (from_user_id = ? AND to_user_id = u.id)
		)
		ORDER BY pm.create_time DESC`

	rows, err := h.DB.Query(query, userID, userID, userID, userID, userID, userID, userID, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []Conversation{}})
		return
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		rows.Scan(&conv.UserID, &conv.Username, &conv.RealName, &conv.Avatar, 
			&conv.LastMessage, &conv.LastMessageTime, &conv.UnreadCount)
		conversations = append(conversations, conv)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": conversations})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	otherUserID, _ := strconv.Atoi(c.Query("user_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))
	offset := (page - 1) * pageSize

	if otherUserID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "用户ID不能为空"})
		return
	}

	query := `SELECT id, from_user_id, from_user_name, to_user_id, content, is_read, create_time
		FROM private_messages 
		WHERE (from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)
		ORDER BY create_time DESC 
		LIMIT ? OFFSET ?`

	rows, err := h.DB.Query(query, userID, otherUserID, otherUserID, userID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.PrivateMessage{}})
		return
	}
	defer rows.Close()

	var messages []models.PrivateMessage
	for rows.Next() {
		var msg models.PrivateMessage
		rows.Scan(&msg.ID, &msg.FromUserID, &msg.FromUserName, &msg.ToUserID, &msg.Content, &msg.IsRead, &msg.CreateTime)
		messages = append(messages, msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	h.DB.Exec(`UPDATE private_messages SET is_read = 1 WHERE from_user_id = ? AND to_user_id = ? AND is_read = 0`, otherUserID, userID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": messages})
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
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
		ToUserID int    `json:"to_user_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "消息内容不能为空"})
		return
	}

	if req.ToUserID == userID {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不能给自己发消息"})
		return
	}

	var targetExists int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND is_deleted = 0", req.ToUserID).Scan(&targetExists)
	if targetExists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	result, err := h.DB.Exec(`INSERT INTO private_messages 
		(from_user_id, from_user_name, to_user_id, content, is_read, create_time)
		VALUES (?, ?, ?, ?, 0, ?)`,
		userID, realName, req.ToUserID, req.Content, now)
	
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "发送失败"})
		return
	}

	msgID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "发送成功",
		"data": map[string]interface{}{
			"id":              msgID,
			"from_user_id":    userID,
			"from_user_name":  realName,
			"to_user_id":      req.ToUserID,
			"content":         req.Content,
			"is_read":         0,
			"create_time":     now,
		},
	})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	otherUserID, _ := strconv.Atoi(c.Query("user_id"))

	if otherUserID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "用户ID不能为空"})
		return
	}

	h.DB.Exec(`UPDATE private_messages SET is_read = 1 WHERE from_user_id = ? AND to_user_id = ? AND is_read = 0`, otherUserID, userID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已标记为已读"})
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var totalUnread int
	h.DB.QueryRow("SELECT COUNT(*) FROM private_messages WHERE to_user_id = ? AND is_read = 0", userID).Scan(&totalUnread)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"total": totalUnread}})
}
