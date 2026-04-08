package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	DB *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) GetRegistrationMode(c *gin.Context) {
	var mode string
	err := h.DB.QueryRow(`SELECT value FROM system_settings WHERE key = 'registration_mode'`).Scan(&mode)
	if err != nil {
		mode = "normal"
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"mode": mode}})
}

func (h *AdminHandler) SetRegistrationMode(c *gin.Context) {
	var req struct {
		Mode string `json:"mode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的请求"})
		return
	}

	if req.Mode != "normal" && req.Mode != "open" && req.Mode != "closed" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的注册模式"})
		return
	}

	_, err := h.DB.Exec(`
		INSERT INTO system_settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = ?, update_time = CURRENT_TIMESTAMP
	`, "registration_mode", req.Mode, req.Mode)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "注册模式已更新", "data": gin.H{"mode": req.Mode}})
}

func (h *AdminHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	var total int
	h.DB.QueryRow(`SELECT COUNT(*) FROM operation_log`).Scan(&total)

	rows, err := h.DB.Query(`
		SELECT id, action, target, operator, detail, create_time 
		FROM operation_log 
		ORDER BY create_time DESC 
		LIMIT ? OFFSET ?`, pageSize, offset)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int
		var action, target, operator, detail, createTime string
		rows.Scan(&id, &action, &target, &operator, &detail, &createTime)
		logs = append(logs, map[string]interface{}{
			"id":          id,
			"action":      action,
			"target":      target,
			"operator":    operator,
			"detail":      detail,
			"create_time": createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     logs,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h *AdminHandler) ClearLogs(c *gin.Context) {
	_, err := h.DB.Exec(`DELETE FROM operation_log`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "清空成功"})
}

func (h *AdminHandler) ClearEarlyLogs(c *gin.Context) {
	countStr := c.Param("count")
	count, _ := strconv.Atoi(countStr)
	if count <= 0 || count > 100 {
		count = 100
	}

	h.DB.Exec(`DELETE FROM operation_log WHERE id IN (SELECT id FROM operation_log ORDER BY create_time ASC LIMIT ?)`, count)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已删除日志"})
}

func (h *AdminHandler) GetPendingUsers(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, username, real_name, role, create_time 
		FROM users 
		WHERE status = 'pending' 
		ORDER BY create_time DESC`)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []map[string]interface{}{}})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, realName, role, createTime string
		rows.Scan(&id, &username, &realName, &role, &createTime)
		users = append(users, map[string]interface{}{
			"id":          id,
			"username":    username,
			"real_name":   realName,
			"role":        role,
			"create_time": createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": users})
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	var total int
	h.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)

	rows, err := h.DB.Query(`
		SELECT id, username, real_name, role, status, avatar, email, phone, team_id, create_time 
		FROM users 
		ORDER BY create_time DESC 
		LIMIT ? OFFSET ?`, pageSize, offset)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, realName, roleStr, status, avatar, email, phone, createTime string
		var teamID sql.NullInt64
		rows.Scan(&id, &username, &realName, &roleStr, &status, &avatar, &email, &phone, &teamID, &createTime)
		users = append(users, map[string]interface{}{
			"id":          id,
			"username":    username,
			"real_name":   realName,
			"role":        roleStr,
			"status":      status,
			"avatar":      avatar,
			"email":       email,
			"phone":       phone,
			"team_id":     teamID.Int64,
			"create_time": createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
