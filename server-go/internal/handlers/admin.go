package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server-go/internal/database"
	"strconv"
	"time"

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

	adminName, _ := c.Get("realName")
	database.AddLog("更新注册模式", "系统设置", adminName.(string), fmt.Sprintf("注册模式变更为: %s", req.Mode))

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

func (h *AdminHandler) GetPendingTeams(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, name, color, description, logo, bg_image, announcement, creator_id, creator_name, create_time
		FROM team_requests
		WHERE status = 'pending'
		ORDER BY create_time DESC`)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []map[string]interface{}{}})
		return
	}
	defer rows.Close()

	var teams []map[string]interface{}
	for rows.Next() {
		var id, creatorID int
		var name, color, description, logo, bgImage, announcement, creatorName, createTime sql.NullString
		rows.Scan(&id, &name, &color, &description, &logo, &bgImage, &announcement, &creatorID, &creatorName, &createTime)
		teams = append(teams, map[string]interface{}{
			"id":           id,
			"name":         name.String,
			"color":        color.String,
			"description":  description.String,
			"logo":         logo.String,
			"bg_image":     bgImage.String,
			"announcement": announcement.String,
			"creator_id":   creatorID,
			"creator_name": creatorName.String,
			"create_time":  createTime.String,
			"status":       "pending",
		})
	}

	if teams == nil {
		teams = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": teams})
}

func (h *AdminHandler) ApproveTeam(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	type ApproveRequest struct {
		ID     int    `json:"id"`
		Action string `json:"action"`
		Remark string `json:"remark"`
	}
	var req ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	now := time.Now().Format(time.RFC3339)

	adminRealNameVal, _ := c.Get("realName")
	adminRealName, _ := adminRealNameVal.(string)

	var creatorID int
	var teamName, color, description, logo, bgImage, announcement string
	err := h.DB.QueryRow(`SELECT creator_id, name, color, description, logo, bg_image, announcement FROM team_requests WHERE id = ? AND status = 'pending'`,
		req.ID).Scan(&creatorID, &teamName, &color, &description, &logo, &bgImage, &announcement)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "申请不存在或已处理"})
		return
	}

	if req.Action == "reject" {
		h.DB.Exec(`UPDATE team_requests SET status = 'rejected', approve_id = ?, approve_name = ?, approve_time = ?, approve_remark = ? WHERE id = ?`,
			0, adminRealName, now, req.Remark, req.ID)

		rejectContent := "很遗憾，您申请的团队【" + teamName + "】已被管理员拒绝"
		if req.Remark != "" {
			rejectContent += "。原因：" + req.Remark
		}

		h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, from_user, create_time)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			creatorID, "team_rejected", "❌ 团队创建申请被拒绝",
			rejectContent,
			"important", adminRealName, now)

		if wsHub := h.getWSHub(); wsHub != nil {
			payload, _ := json.Marshal(map[string]interface{}{
				"type": "notification",
				"data": map[string]interface{}{
					"type":    "team_rejected",
					"title":   "❌ 团队创建申请被拒绝",
					"content": rejectContent,
				},
			})
			wsHub.SendToUser(creatorID, payload)
		}

		database.AddLog("审核", "团队申请【"+teamName+"】", adminRealName, "审核拒绝")
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已拒绝申请"})
		return
	}

	colorVal := color
	if colorVal == "" {
		colorVal = "#3b82f6"
	}

	tx, _ := h.DB.Begin()
	result, err := tx.Exec(`INSERT INTO teams (name, color, description, logo, bg_image, announcement, leader_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		teamName, colorVal, description, logo, bgImage, announcement, creatorID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建团队失败"})
		return
	}
	teamID, _ := result.LastInsertId()

	_, err = tx.Exec(`UPDATE team_requests SET status = 'approved', approve_id = ?, approve_name = ?, approve_time = ? WHERE id = ?`,
		0, adminRealName, now, req.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新申请状态失败"})
		return
	}

	_, err = tx.Exec(`UPDATE users SET team_id = ? WHERE id = ?`, teamID, creatorID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新用户团队失败"})
		return
	}

	tx.Commit()

	h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, from_user, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		creatorID, "team_approved", "🎉 团队创建申请通过",
		"恭喜！您申请的团队【"+teamName+"】已通过审核，现在可以开始使用了！",
		"important", adminRealName, now)

	if wsHub := h.getWSHub(); wsHub != nil {
		payload, _ := json.Marshal(map[string]interface{}{
			"type": "notification",
			"data": map[string]interface{}{
				"type":    "team_approved",
				"title":   "🎉 团队创建申请通过",
				"content": "恭喜！您申请的团队【" + teamName + "】已通过审核，现在可以开始使用了！",
			},
		})
		wsHub.SendToUser(creatorID, payload)
	}

	database.AddLog("审核", "团队申请【"+teamName+"】", adminRealName, "审核通过")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已通过申请，团队创建成功"})
}

func (h *AdminHandler) BatchApproveUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		UserIDs []int `json:"user_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.UserIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	adminName, _ := c.Get("realName")
	for _, id := range req.UserIDs {
		h.DB.Exec("UPDATE users SET status = 'active', update_time = CURRENT_TIMESTAMP WHERE id = ? AND status = 'pending'", id)
		database.AddLog("批量通过", "用户管理", adminName.(string), fmt.Sprintf("批量通过用户ID=%d", id))
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": fmt.Sprintf("已通过 %d 个用户", len(req.UserIDs))})
}

func (h *AdminHandler) BatchRejectUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		UserIDs []int `json:"user_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.UserIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	adminName, _ := c.Get("realName")
	for _, id := range req.UserIDs {
		h.DB.Exec("UPDATE users SET status = 'rejected', update_time = CURRENT_TIMESTAMP WHERE id = ? AND status = 'pending'", id)
		database.AddLog("批量拒绝", "用户管理", adminName.(string), fmt.Sprintf("批量拒绝用户ID=%d", id))
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": fmt.Sprintf("已拒绝 %d 个用户", len(req.UserIDs))})
}

func (h *AdminHandler) getWSHub() *WSHub {
	return Hub
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, username, real_name, role, status, avatar, email, phone, team_id, create_time 
		FROM users 
		WHERE is_deleted = 0
		ORDER BY create_time DESC`)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []map[string]interface{}{}})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, realName, roleStr, status string
		var avatar, email, phone, createTime sql.NullString
		var teamID sql.NullInt64
		err := rows.Scan(&id, &username, &realName, &roleStr, &status, &avatar, &email, &phone, &teamID, &createTime)
		if err != nil {
			continue
		}

		user := map[string]interface{}{
			"id":          id,
			"username":    username,
			"real_name":   realName,
			"role":        roleStr,
			"status":      status,
			"avatar":      avatar.String,
			"email":       email.String,
			"phone":       phone.String,
			"team_id":     nil,
			"team_name":   "",
			"team_color":  "",
			"create_time": createTime.String,
		}
		if teamID.Valid {
			user["team_id"] = teamID.Int64
			var teamName, teamColor sql.NullString
			h.DB.QueryRow("SELECT name, color FROM teams WHERE id = ?", teamID.Int64).Scan(&teamName, &teamColor)
			user["team_name"] = teamName.String
			user["team_color"] = teamColor.String
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": users,
	})
}
