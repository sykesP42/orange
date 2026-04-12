package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"server-go/internal/database"
)

type TeamCollabHandler struct{}

func NewTeamCollabHandler() *TeamCollabHandler {
	return &TeamCollabHandler{}
}

type BloggerTransfer struct {
	ID              int       `json:"id"`
	BloggerID      int       `json:"blogger_id"`
	BloggerName    string    `json:"blogger_name"`
	FromUserID     int       `json:"from_user_id"`
	FromUserName   string    `json:"from_user_name"`
	ToUserID       int       `json:"to_user_id"`
	ToUserName     string    `json:"to_user_name"`
	Status         string    `json:"status"`
	Remark         string    `json:"remark"`
	CreatedAt      string    `json:"created_at"`
	HandledAt      string    `json:"handled_at,omitempty"`
}

func (h *TeamCollabHandler) GetTeamMembers(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var teamID sql.NullInt64
	database.DB.QueryRow(`SELECT team_id FROM users WHERE id = ? AND status = 'active'`, userID).Scan(&teamID)
	if !teamID.Valid || teamID.Int64 == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "未加入团队"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT u.id, u.username, u.real_name, u.avatar, t.name as team_name
		FROM users u
		INNER JOIN teams t ON u.team_id = t.id
		WHERE u.team_id = ? AND u.status = 'active'
		ORDER BY u.id ASC
	`, teamID.Int64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取团队成员失败"})
		return
	}
	defer rows.Close()

	var members []map[string]interface{}
	for rows.Next() {
		var id int
		var username, realName, avatar, teamName string
		rows.Scan(&id, &username, &realName, &avatar, &teamName)
		members = append(members, map[string]interface{}{
			"id":         id,
			"username":  username,
			"real_name": realName,
			"avatar":    avatar,
			"team_name": teamName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": members})
}

func (h *TeamCollabHandler) RequestTransfer(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		BloggerID int    `json:"blogger_id" binding:"required"`
		ToUserID  int    `json:"to_user_id" binding:"required"`
		Remark    string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if userID == req.ToUserID {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不能移交给自己"})
		return
	}

	var bloggerName string
	var ownerName sql.NullString
	database.DB.QueryRow(`SELECT nickname, COALESCE(user_name, '') FROM blogger WHERE id = ?`, req.BloggerID).Scan(&bloggerName, &ownerName)

	var fromRealName sql.NullString
	database.DB.QueryRow(`SELECT COALESCE(real_name, '') FROM users WHERE id = ?`, userID).Scan(&fromRealName)

	var toRealName sql.NullString
	database.DB.QueryRow(`SELECT COALESCE(real_name, '') FROM users WHERE id = ?`, req.ToUserID).Scan(&toRealName)

	_, err := database.DB.Exec(`
		INSERT INTO blogger_transfer_requests
		(blogger_id, blogger_name, from_user_id, from_user_name, to_user_id, to_user_name, status, remark, create_time)
		VALUES (?, ?, ?, ?, ?, ?, 'pending', ?, CURRENT_TIMESTAMP)
	`, req.BloggerID, bloggerName, userID, fromRealName.String, req.ToUserID, toRealName.String, req.Remark)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "发起移交请求失败"})
		return
	}

	database.DB.Exec(`
		INSERT INTO notifications (user_id, type, title, content, redirect_url, is_read, create_time)
		VALUES (?, 'transfer_request', '📦 博主移交请求', ?, ?, 0, CURRENT_TIMESTAMP)
	`, req.ToUserID, fromRealName.String+" 申请移交博主 "+bloggerName, "/team/transfers")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "移交请求已发送"})
}

func (h *TeamCollabHandler) GetTransferRequests(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	status := c.DefaultQuery("status", "all")

	var query string
	var args []interface{}

	if status == "pending" {
		query = `
			SELECT id, blogger_id, blogger_name, from_user_id, from_user_name, to_user_id, to_user_name, status, remark, created_at
			FROM blogger_transfer_requests
			WHERE to_user_id = ? AND status = 'pending'
			ORDER BY create_time DESC`
		args = []interface{}{userID}
	} else {
		query = `
			SELECT id, blogger_id, blogger_name, from_user_id, from_user_name, to_user_id, to_user_name, status, remark, created_at
			FROM blogger_transfer_requests
			WHERE (from_user_id = ? OR to_user_id = ?)
			ORDER BY create_time DESC LIMIT 50`
		args = []interface{}{userID, userID}
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取移交请求失败"})
		return
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id, bloggerID, fromUserID, toUserID int
		var bloggerName, fromUserName, toUserName, reqStatus, remark, createdAt string
		rows.Scan(&id, &bloggerID, &bloggerName, &fromUserID, &fromUserName, &toUserID, &toUserName, &reqStatus, &remark, &createdAt)
		requests = append(requests, map[string]interface{}{
			"id":            id,
			"blogger_id":   bloggerID,
			"blogger_name": bloggerName,
			"from_user_id": fromUserID,
			"from_user_name": fromUserName,
			"to_user_id":   toUserID,
			"to_user_name": toUserName,
			"status":       reqStatus,
			"remark":       remark,
			"created_at":   createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": requests})
}

func (h *TeamCollabHandler) HandleTransferRequest(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		RequestID int    `json:"request_id" binding:"required"`
		Action    string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var bloggerID, fromUserID int
	var bloggerName string
	err := database.DB.QueryRow(`
		SELECT blogger_id, from_user_id, blogger_name FROM blogger_transfer_requests
		WHERE id = ? AND to_user_id = ? AND status = 'pending'
	`, req.RequestID, userID).Scan(&bloggerID, &fromUserID, &bloggerName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "请求不存在或无权处理"})
		return
	}

	if req.Action == "approve" {
		var toUsername sql.NullString
		database.DB.QueryRow(`SELECT username FROM users WHERE id = ?`, userID).Scan(&toUsername)

		database.DB.Exec(`UPDATE blogger SET user_name = ? WHERE id = ?`, toUsername.String, bloggerID)

		database.DB.Exec(`UPDATE blogger_transfer_requests SET status = 'approved', handled_at = CURRENT_TIMESTAMP WHERE id = ?`, req.RequestID)

		database.DB.Exec(`
			INSERT INTO notifications (user_id, type, title, content, redirect_url, is_read, create_time)
			VALUES (?, 'transfer_approved', '✅ 博主移交已同意', ?, ?, 0, CURRENT_TIMESTAMP)
		`, fromUserID, "您移交的博主 "+bloggerName+" 已被接收", bloggerID)

		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已同意移交"})
	} else if req.Action == "reject" {
		database.DB.Exec(`UPDATE blogger_transfer_requests SET status = 'rejected', handled_at = CURRENT_TIMESTAMP WHERE id = ?`, req.RequestID)

		database.DB.Exec(`
			INSERT INTO notifications (user_id, type, title, content, is_read, create_time)
			VALUES (?, 'transfer_rejected', '❌ 博主移交被拒绝', ?, 0, CURRENT_TIMESTAMP)
		`, fromUserID, "您移交的博主 "+bloggerName+" 被拒绝了")

		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已拒绝移交"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的操作"})
	}
}

func (h *TeamCollabHandler) GetTeamTasks(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT t.id, t.title, t.content, t.status, t.assignee_id, t.assigner_id,
			COALESCE(u1.real_name, '') as assignee_name, COALESCE(u2.real_name, '') as assigner_name,
			t.deadline, t.create_time, t.update_time,
			COALESCE(b.nickname, '') as blogger_name
		FROM team_tasks t
		LEFT JOIN users u1 ON t.assignee_id = u1.id
		LEFT JOIN users u2 ON t.assigner_id = u2.id
		LEFT JOIN blogger b ON t.blogger_id = b.id
		WHERE t.assignee_id = ? OR t.assigner_id = ? OR t.assignee_id IS NULL
		ORDER BY t.create_time DESC
		LIMIT 50
	`, userID, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取任务失败"})
		return
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	for rows.Next() {
		var id, assigneeID, assignerID int
		var title, content, reqStatus, assigneeName, assignerName, deadline, createTime, updateTime, bloggerName string
		rows.Scan(&id, &title, &content, &reqStatus, &assigneeID, &assignerID,
			&assigneeName, &assignerName, &deadline, &createTime, &updateTime, &bloggerName)

		tasks = append(tasks, map[string]interface{}{
			"id":            id,
			"title":        title,
			"content":      content,
			"status":       reqStatus,
			"assignee_id":  assigneeID,
			"assigner_id":  assignerID,
			"assignee_name": assigneeName,
			"assigner_name": assignerName,
			"deadline":     deadline,
			"blogger_name": bloggerName,
			"create_time":  createTime,
			"update_time":  updateTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tasks})
}

func (h *TeamCollabHandler) CreateTeamTask(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content"`
		AssigneeID int    `json:"assignee_id"`
		BloggerID int    `json:"blogger_id"`
		Deadline  string `json:"deadline"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO team_tasks (title, content, assignee_id, assigner_id, blogger_id, deadline, status, create_time)
		VALUES (?, ?, ?, ?, ?, ?, 'pending', CURRENT_TIMESTAMP)
	`, req.Title, req.Content, req.AssigneeID, userID, req.BloggerID, req.Deadline)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建任务失败"})
		return
	}

	taskID, _ := result.LastInsertId()

	if req.AssigneeID > 0 {
		database.DB.Exec(`
			INSERT INTO notifications (user_id, type, title, content, redirect_url, is_read, create_time)
			VALUES (?, 'task_assigned', '📋 新任务', ?, '/team/tasks', 0, CURRENT_TIMESTAMP)
		`, req.AssigneeID, "您被分配了新任务: "+req.Title)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": map[string]interface{}{
		"task_id": taskID,
	}, "message": "任务创建成功"})
}

func (h *TeamCollabHandler) UpdateTaskStatus(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		TaskID int    `json:"task_id" binding:"required"`
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	result, err := database.DB.Exec(`
		UPDATE team_tasks SET status = ?, update_time = CURRENT_TIMESTAMP
		WHERE id = ? AND (assignee_id = ? OR assigner_id = ?)
	`, req.Status, req.TaskID, userID, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新任务失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "任务不存在或无权更新"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "任务状态已更新"})
}
