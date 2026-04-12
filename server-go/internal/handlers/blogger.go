package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server-go/internal/database"
	"server-go/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type BloggerHandler struct {
	DB *sql.DB
}

func NewBloggerHandler(db *sql.DB) *BloggerHandler {
	return &BloggerHandler{DB: db}
}

func (h *BloggerHandler) GetBloggers(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	platform := c.Query("platform")
	username := c.Query("username")
	status := c.Query("status")
	teamID := c.Query("team_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	query := `SELECT b.id, b.nickname, b.category, b.product, b.contact, b.wechat, b.custom_contact, 
		b.platform, b.platform_account, b.description, b.tags, b.avatar, b.user_name, b.real_name, 
		b.status, b.status_update_time, b.status_remark, b.is_deleted, b.is_invalid, b.create_time, b.update_time 
		FROM blogger b WHERE b.is_deleted = 0`
	countQuery := `SELECT COUNT(*) FROM blogger b WHERE b.is_deleted = 0`
	args := []interface{}{}
	countArgs := []interface{}{}

	if keyword != "" {
		kw := "%" + strings.ToLower(keyword) + "%"
		query += " AND (LOWER(b.nickname) LIKE ? OR LOWER(b.product) LIKE ? OR LOWER(b.description) LIKE ?)"
		countQuery += " AND (LOWER(b.nickname) LIKE ? OR LOWER(b.product) LIKE ? OR LOWER(b.description) LIKE ?)"
		args = append(args, kw, kw, kw)
		countArgs = append(countArgs, kw, kw, kw)
	}
	if category != "" {
		query += " AND b.category = ?"
		countQuery += " AND b.category = ?"
		args = append(args, category)
		countArgs = append(countArgs, category)
	}
	if platform != "" {
		query += " AND b.platform = ?"
		countQuery += " AND b.platform = ?"
		args = append(args, platform)
		countArgs = append(countArgs, platform)
	}
	if username != "" {
		query += " AND b.user_name = ?"
		countQuery += " AND b.user_name = ?"
		args = append(args, username)
		countArgs = append(countArgs, username)
	}
	if status != "" {
		query += " AND b.status = ?"
		countQuery += " AND b.status = ?"
		args = append(args, status)
		countArgs = append(countArgs, status)
	}
	if teamID != "" {
		query += " AND EXISTS (SELECT 1 FROM users u WHERE u.username = b.user_name AND u.team_id = ? AND u.is_deleted = 0)"
		countQuery += " AND EXISTS (SELECT 1 FROM users u WHERE u.username = b.user_name AND u.team_id = ? AND u.is_deleted = 0)"
		teamIDInt, _ := strconv.Atoi(teamID)
		args = append(args, teamIDInt)
		countArgs = append(countArgs, teamIDInt)
	}

	query += " ORDER BY b.create_time DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	var total int
	h.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	defer rows.Close()

	var bloggerList []map[string]interface{}
	for rows.Next() {
		var b models.Blogger
		rows.Scan(&b.ID, &b.Nickname, &b.Category, &b.Product, &b.Contact, &b.Wechat, &b.CustomContact,
			&b.Platform, &b.PlatformAccount, &b.Description, &b.Tags, &b.Avatar, &b.UserName, &b.RealName,
			&b.Status, &b.StatusUpdateTime, &b.StatusRemark, &b.IsDeleted, &b.IsInvalid, &b.CreateTime, &b.UpdateTime)

		var userAvatar, userEmail, userPhone, userBio, userRole sql.NullString
		var userTeamID sql.NullInt64
		var teamName, teamColor sql.NullString
		h.DB.QueryRow(`SELECT u.avatar, u.email, u.phone, u.bio, u.role, u.team_id, t.name, t.color 
			FROM users u LEFT JOIN teams t ON u.team_id = t.id WHERE u.username = ?`, b.UserName).Scan(
			&userAvatar, &userEmail, &userPhone, &userBio, &userRole, &userTeamID, &teamName, &teamColor)

		bloggerMap := map[string]interface{}{
			"id":                 b.ID,
			"nickname":           b.Nickname,
			"category":           b.Category,
			"product":            b.Product,
			"contact":            b.Contact,
			"wechat":             b.Wechat,
			"custom_contact":     b.CustomContact,
			"platform":           b.Platform,
			"platform_account":   b.PlatformAccount,
			"description":        b.Description,
			"avatar":             b.Avatar,
			"user_name":          b.UserName,
			"real_name":          b.RealName,
			"status":             b.Status,
			"status_update_time": b.StatusUpdateTime,
			"status_remark":      b.StatusRemark,
			"is_deleted":         b.IsDeleted,
			"is_invalid":         b.IsInvalid,
			"create_time":        b.CreateTime,
			"update_time":        b.UpdateTime,
			"owner_username":     b.UserName,
			"owner_real_name":    b.RealName,
			"owner_avatar":       userAvatar.String,
			"email":              userEmail.String,
			"phone":              userPhone.String,
			"bio":                userBio.String,
			"role":               userRole.String,
			"team_name":          teamName.String,
			"team_color":         teamColor.String,
		}

		if b.Avatar == "" {
			bloggerMap["avatar"] = userAvatar.String
		}

		if userTeamID.Valid {
			bloggerMap["team_id"] = userTeamID.Int64
		} else {
			bloggerMap["team_id"] = nil
		}

		if b.Tags != "" {
			var tags []string
			if err := json.Unmarshal([]byte(b.Tags), &tags); err == nil {
				bloggerMap["tags"] = tags
			} else {
				bloggerMap["tags"] = []string{}
			}
		} else {
			bloggerMap["tags"] = []string{}
		}

		bloggerList = append(bloggerList, bloggerMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     bloggerList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h *BloggerHandler) GetMyBloggers(c *gin.Context) {
	userName, _ := c.Get("username")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	query := `SELECT b.id, b.nickname, b.category, b.product, b.contact, b.wechat, b.custom_contact, 
		b.platform, b.platform_account, b.description, b.tags, b.avatar, b.user_name, b.real_name, 
		b.status, b.status_update_time, b.status_remark, b.is_deleted, b.is_invalid, b.create_time, b.update_time 
		FROM blogger b WHERE b.is_deleted = 0 AND b.user_name = ?
		ORDER BY b.create_time DESC LIMIT ? OFFSET ?`

	countQuery := `SELECT COUNT(*) FROM blogger b WHERE b.is_deleted = 0 AND b.user_name = ?`

	var total int
	h.DB.QueryRow(countQuery, userName).Scan(&total)

	rows, err := h.DB.Query(query, userName, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	defer rows.Close()

	var bloggerList []map[string]interface{}
	for rows.Next() {
		var b models.Blogger
		rows.Scan(&b.ID, &b.Nickname, &b.Category, &b.Product, &b.Contact, &b.Wechat, &b.CustomContact,
			&b.Platform, &b.PlatformAccount, &b.Description, &b.Tags, &b.Avatar, &b.UserName, &b.RealName,
			&b.Status, &b.StatusUpdateTime, &b.StatusRemark, &b.IsDeleted, &b.IsInvalid, &b.CreateTime, &b.UpdateTime)

		var userAvatar sql.NullString
		h.DB.QueryRow(`SELECT avatar FROM users WHERE username = ?`, b.UserName).Scan(&userAvatar)

		bloggerMap := map[string]interface{}{
			"id":                 b.ID,
			"nickname":           b.Nickname,
			"category":           b.Category,
			"product":            b.Product,
			"contact":            b.Contact,
			"wechat":             b.Wechat,
			"custom_contact":     b.CustomContact,
			"platform":           b.Platform,
			"platform_account":   b.PlatformAccount,
			"description":        b.Description,
			"avatar":             b.Avatar,
			"user_name":          b.UserName,
			"real_name":          b.RealName,
			"status":             b.Status,
			"status_update_time": b.StatusUpdateTime,
			"status_remark":      b.StatusRemark,
			"is_deleted":         b.IsDeleted,
			"is_invalid":         b.IsInvalid,
			"create_time":        b.CreateTime,
			"update_time":        b.UpdateTime,
		}

		if b.Avatar == "" {
			bloggerMap["avatar"] = userAvatar.String
		}

		if b.Tags != "" {
			var tags []string
			if err := json.Unmarshal([]byte(b.Tags), &tags); err == nil {
				bloggerMap["tags"] = tags
			} else {
				bloggerMap["tags"] = []string{}
			}
		} else {
			bloggerMap["tags"] = []string{}
		}

		bloggerList = append(bloggerList, bloggerMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     bloggerList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func (h *BloggerHandler) GetExpiringBloggers(c *gin.Context) {
	userName, _ := c.Get("username")
	now := time.Now()

	rows, err := h.DB.Query(`SELECT b.id, b.nickname, b.create_time FROM blogger b 
		WHERE b.is_deleted = 0 AND b.is_invalid = 0 
		AND b.contact = '' AND b.wechat = '' AND b.custom_contact = ''
		AND b.user_name = ?`, userName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Blogger{}})
		return
	}
	defer rows.Close()

	var bloggers []map[string]interface{}
	for rows.Next() {
		var id int
		var nickname, createTime string
		rows.Scan(&id, &nickname, &createTime)

		ct, _ := time.Parse(time.RFC3339, createTime)
		if ct.IsZero() {
			ct, _ = time.Parse("2006-01-02 15:04:05", createTime)
		}
		diffTime := ct.Add(15*24*time.Hour).Unix() - now.Unix()
		daysLeft := int(diffTime / (24 * 60 * 60))

		if daysLeft > 0 && daysLeft <= 15 {
			bloggers = append(bloggers, map[string]interface{}{
				"id":          id,
				"nickname":    nickname,
				"days_left":   daysLeft,
				"create_time": createTime,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": bloggers})
}

func (h *BloggerHandler) CreateBlogger(c *gin.Context) {
	userNameVal, _ := c.Get("username")
	userName, _ := userNameVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}

	var req struct {
		Nickname        string   `json:"nickname" binding:"required"`
		Category        string   `json:"category" binding:"required"`
		Product         string   `json:"product"`
		Platform        string   `json:"platform" binding:"required"`
		PlatformAccount string   `json:"platform_account" binding:"required"`
		Description     string   `json:"description"`
		Tags            []string `json:"tags"`
		Avatar          string   `json:"avatar"`
		Contact         string   `json:"contact"`
		Wechat          string   `json:"wechat"`
		CustomContact   string   `json:"custom_contact"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "昵称、分类、平台和平台账号不能为空"})
		return
	}

	var existingRealName string
	err := h.DB.QueryRow(`SELECT real_name FROM blogger 
		WHERE platform = ? AND platform_account = ? AND is_deleted = 0`,
		req.Platform, req.PlatformAccount).Scan(&existingRealName)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    4001,
			"message": "该平台账号已被【" + existingRealName + "】登记，不可重复录入",
			"data":    existingRealName,
		})
		return
	}

	finalNickname := req.Nickname
	counter := 1
	for {
		var count int
		h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE nickname = ? AND is_deleted = 0", finalNickname).Scan(&count)
		if count == 0 {
			break
		}
		counter++
		finalNickname = req.Nickname + "(" + strconv.Itoa(counter) + ")"
	}

	tagsJSON, _ := json.Marshal(req.Tags)
	now := time.Now().Format(time.RFC3339)

	result, err := h.DB.Exec(`INSERT INTO blogger 
		(nickname, category, product, contact, wechat, custom_contact, platform, platform_account, 
		description, tags, avatar, user_name, real_name, status, status_update_time, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		finalNickname, req.Category, req.Product, req.Contact, req.Wechat, req.CustomContact,
		req.Platform, req.PlatformAccount, req.Description, string(tagsJSON), req.Avatar,
		userName, realName, "初次联系", now, now, now)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	h.DB.Exec(`INSERT INTO blogger_status_history (blogger_id, old_status, new_status, operator, remark, create_time)
		VALUES (?, NULL, ?, ?, ?, ?)`,
		id, "初次联系", realName, "博主初次录入", now)

	database.AddLog("新增", finalNickname, realName, "新增博主【"+finalNickname+"】("+req.Platform+": "+req.PlatformAccount+")")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "录入成功",
		"data":    gin.H{"id": id, "nickname": finalNickname},
	})
}

func (h *BloggerHandler) GetBlogger(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var b models.Blogger
	var tagsStr string
	err := h.DB.QueryRow(`SELECT id, nickname, category, product, contact, wechat, custom_contact, 
		platform, platform_account, description, tags, avatar, user_name, real_name, 
		status, status_update_time, create_time, update_time 
		FROM blogger WHERE id = ? AND is_deleted = 0`, id).Scan(
		&b.ID, &b.Nickname, &b.Category, &b.Product, &b.Contact, &b.Wechat, &b.CustomContact,
		&b.Platform, &b.PlatformAccount, &b.Description, &tagsStr, &b.Avatar, &b.UserName, &b.RealName,
		&b.Status, &b.StatusUpdateTime, &b.CreateTime, &b.UpdateTime)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	var userAvatar, userEmail, userPhone, userBio, userRole sql.NullString
	var userTeamID sql.NullInt64
	var teamName, teamColor sql.NullString
	h.DB.QueryRow(`SELECT u.avatar, u.email, u.phone, u.bio, u.role, u.team_id, t.name, t.color 
		FROM users u LEFT JOIN teams t ON u.team_id = t.id WHERE u.username = ?`, b.UserName).Scan(
		&userAvatar, &userEmail, &userPhone, &userBio, &userRole, &userTeamID, &teamName, &teamColor)

	bloggerMap := map[string]interface{}{
		"id":                 b.ID,
		"nickname":           b.Nickname,
		"category":           b.Category,
		"product":            b.Product,
		"contact":            b.Contact,
		"wechat":             b.Wechat,
		"custom_contact":     b.CustomContact,
		"platform":           b.Platform,
		"platform_account":   b.PlatformAccount,
		"description":        b.Description,
		"avatar":             b.Avatar,
		"user_name":          b.UserName,
		"real_name":          b.RealName,
		"status":             b.Status,
		"status_update_time": b.StatusUpdateTime,
		"create_time":        b.CreateTime,
		"update_time":        b.UpdateTime,
		"owner_username":     b.UserName,
		"owner_real_name":    b.RealName,
		"owner_avatar":       userAvatar.String,
		"email":              userEmail.String,
		"phone":              userPhone.String,
		"bio":                userBio.String,
		"role":               userRole.String,
		"team_name":          teamName.String,
		"team_color":         teamColor.String,
	}

	if b.Avatar == "" {
		bloggerMap["avatar"] = userAvatar.String
	}

	if userTeamID.Valid {
		bloggerMap["team_id"] = userTeamID.Int64
	} else {
		bloggerMap["team_id"] = nil
	}

	if tagsStr != "" {
		var tags []string
		if err := json.Unmarshal([]byte(tagsStr), &tags); err == nil {
			bloggerMap["tags"] = tags
		} else {
			bloggerMap["tags"] = []string{}
		}
	} else {
		bloggerMap["tags"] = []string{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": bloggerMap})
}

func (h *BloggerHandler) UpdateBlogger(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userNameVal, _ := c.Get("username")
	userName, _ := userNameVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	var b models.Blogger
	err := h.DB.QueryRow(`SELECT user_name, nickname, status FROM blogger WHERE id = ? AND is_deleted = 0`, id).Scan(
		&b.UserName, &b.Nickname, &b.Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	if role != "admin" && b.UserName != userName {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限修改此博主信息"})
		return
	}

	var req struct {
		Nickname        string   `json:"nickname"`
		Category        string   `json:"category"`
		Product         string   `json:"product"`
		Contact         string   `json:"contact"`
		Wechat          string   `json:"wechat"`
		CustomContact   string   `json:"custom_contact"`
		Platform        string   `json:"platform"`
		PlatformAccount string   `json:"platform_account"`
		Description     string   `json:"description"`
		Tags            []string `json:"tags"`
		Status          string   `json:"status"`
		StatusRemark    string   `json:"status_remark"`
		Avatar          string   `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.Nickname != "" {
		updates = append(updates, "nickname = ?")
		args = append(args, req.Nickname)
	}
	if req.Category != "" {
		updates = append(updates, "category = ?")
		args = append(args, req.Category)
	}
	if req.Product != "" {
		updates = append(updates, "product = ?")
		args = append(args, req.Product)
	}
	if req.Contact != "" {
		updates = append(updates, "contact = ?")
		args = append(args, req.Contact)
	}
	if req.Wechat != "" {
		updates = append(updates, "wechat = ?")
		args = append(args, req.Wechat)
	}
	if req.CustomContact != "" {
		updates = append(updates, "custom_contact = ?")
		args = append(args, req.CustomContact)
	}
	if req.Platform != "" {
		updates = append(updates, "platform = ?")
		args = append(args, req.Platform)
	}
	if req.PlatformAccount != "" {
		updates = append(updates, "platform_account = ?")
		args = append(args, req.PlatformAccount)
	}
	if req.Description != "" {
		updates = append(updates, "description = ?")
		args = append(args, req.Description)
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		updates = append(updates, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	if req.Avatar != "" {
		updates = append(updates, "avatar = ?")
		args = append(args, req.Avatar)
	}

	oldStatus := b.Status
	if req.Status != "" && req.Status != oldStatus {
		updates = append(updates, "status = ?, status_update_time = ?")
		now := time.Now().Format(time.RFC3339)
		args = append(args, req.Status, now)

		h.DB.Exec(`INSERT INTO blogger_status_history (blogger_id, old_status, new_status, operator, remark, create_time)
			VALUES (?, ?, ?, ?, ?, ?)`,
			id, oldStatus, req.Status, realName, req.StatusRemark, now)

		database.AddLog("状态变更", b.Nickname, realName, "博主状态从【"+oldStatus+"】变更为【"+req.Status+"】")

		if (req.Status == "失效" || req.Status == "无效") && b.UserName != "" {
			var userID int
			h.DB.QueryRow("SELECT id FROM users WHERE username = ?", b.UserName).Scan(&userID)
			if userID > 0 {
				h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, blogger_id, from_user, create_time)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
					userID, "invalid_blogger", "⚠️ 博主状态变更提醒",
					"博主【"+b.Nickname+"】的状态已变更为【"+req.Status+"】，请及时处理。",
					"important", id, realName, time.Now().Format(time.RFC3339))
			}
		}
	}

	updates = append(updates, "update_time = ?")
	args = append(args, time.Now().Format(time.RFC3339))
	args = append(args, id)

	if len(updates) > 0 {
		query := "UPDATE blogger SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	database.AddLog("编辑", b.Nickname, realName, "编辑博主信息")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *BloggerHandler) DeleteBlogger(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	if idStr != "" {
		id, _ = strconv.Atoi(idStr)
	} else {
		var req struct {
			ID int `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err == nil {
			id = req.ID
		}
	}

	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	if id == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "缺少博主ID"})
		return
	}

	var nickname string
	h.DB.QueryRow("SELECT nickname FROM blogger WHERE id = ?", id).Scan(&nickname)

	h.DB.Exec("UPDATE blogger SET is_deleted = 1, update_time = ? WHERE id = ?", time.Now().Format(time.RFC3339), id)

	database.AddLog("删除", nickname, realName, "删除博主【"+nickname+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *BloggerHandler) GetStatusHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE id = ? AND is_deleted = 0", id).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	rows, err := h.DB.Query(`SELECT id, blogger_id, old_status, new_status, operator, remark, create_time 
		FROM blogger_status_history WHERE blogger_id = ? ORDER BY create_time DESC`, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.BloggerStatusHistory{}})
		return
	}
	defer rows.Close()

	var history []models.BloggerStatusHistory
	for rows.Next() {
		var h models.BloggerStatusHistory
		rows.Scan(&h.ID, &h.BloggerID, &h.OldStatus, &h.NewStatus, &h.Operator, &h.Remark, &h.CreateTime)
		history = append(history, h)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": history})
}

func (h *BloggerHandler) GetFollowups(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE id = ? AND is_deleted = 0", id).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	rows, err := h.DB.Query(`SELECT id, blogger_id, content, type, operator, create_time 
		FROM followup WHERE blogger_id = ? ORDER BY create_time DESC`, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Followup{}})
		return
	}
	defer rows.Close()

	var followups []models.Followup
	for rows.Next() {
		var f models.Followup
		rows.Scan(&f.ID, &f.BloggerID, &f.Content, &f.Type, &f.Operator, &f.CreateTime)
		followups = append(followups, f)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": followups})
}

func (h *BloggerHandler) CreateFollowup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE id = ? AND is_deleted = 0", id).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	followupType := req.Type
	if followupType == "" {
		followupType = "跟进"
	}

	result, err := h.DB.Exec(`INSERT INTO followup (blogger_id, content, type, operator) 
		VALUES (?, ?, ?, ?)`, id, req.Content, followupType, realName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	followupID, _ := result.LastInsertId()

	database.AddLog("跟进", "博主跟进记录", realName, "添加博主跟进记录")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加成功",
		"data":    gin.H{"id": followupID},
	})
}

func (h *BloggerHandler) DeleteFollowup(c *gin.Context) {
	followupID, _ := strconv.Atoi(c.Param("followupId"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	var operator string
	h.DB.QueryRow("SELECT operator FROM followup WHERE id = ?", followupID).Scan(&operator)
	if operator == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "跟进记录不存在"})
		return
	}

	if role != "admin" && operator != realName {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此跟进记录"})
		return
	}

	h.DB.Exec("DELETE FROM followup WHERE id = ?", followupID)

	database.AddLog("删除", "跟进记录", realName, "删除跟进记录")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *BloggerHandler) GetCooperations(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE id = ? AND is_deleted = 0", id).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	rows, err := h.DB.Query(`SELECT id, blogger_id, title, date, status, product, amount, note, create_time, update_time 
		FROM cooperation WHERE blogger_id = ? ORDER BY create_time DESC`, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Cooperation{}})
		return
	}
	defer rows.Close()

	var cooperations []models.Cooperation
	for rows.Next() {
		var coop models.Cooperation
		rows.Scan(&coop.ID, &coop.BloggerID, &coop.Title, &coop.Date, &coop.Status, &coop.Product, &coop.Amount, &coop.Note, &coop.CreateTime, &coop.UpdateTime)
		cooperations = append(cooperations, coop)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": cooperations})
}

func (h *BloggerHandler) CreateCooperation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE id = ? AND is_deleted = 0", id).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	var req struct {
		Title   string  `json:"title" binding:"required"`
		Date    string  `json:"date"`
		Status  string  `json:"status"`
		Product string  `json:"product"`
		Amount  float64 `json:"amount"`
		Note    string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "标题不能为空"})
		return
	}

	status := req.Status
	if status == "" {
		status = "进行中"
	}

	now := time.Now().Format(time.RFC3339)
	result, err := h.DB.Exec(`INSERT INTO cooperation (blogger_id, title, date, status, product, amount, note, create_time, update_time) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, req.Title, req.Date, status, req.Product, req.Amount, req.Note, now, now)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	coopID, _ := result.LastInsertId()

	database.AddLog("合作", "合作记录", realName, "添加合作记录【"+req.Title+"】")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加成功",
		"data":    gin.H{"id": coopID},
	})
}

func (h *BloggerHandler) UpdateCooperation(c *gin.Context) {
	coopID, _ := strconv.Atoi(c.Param("coopId"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM cooperation WHERE id = ?", coopID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合作记录不存在"})
		return
	}

	var req struct {
		Title   string  `json:"title"`
		Date    string  `json:"date"`
		Status  string  `json:"status"`
		Product string  `json:"product"`
		Amount  float64 `json:"amount"`
		Note    string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.Title != "" {
		updates = append(updates, "title = ?")
		args = append(args, req.Title)
	}
	if req.Date != "" {
		updates = append(updates, "date = ?")
		args = append(args, req.Date)
	}
	if req.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, req.Status)
	}
	if req.Product != "" {
		updates = append(updates, "product = ?")
		args = append(args, req.Product)
	}
	if req.Amount != 0 {
		updates = append(updates, "amount = ?")
		args = append(args, req.Amount)
	}
	if req.Note != "" {
		updates = append(updates, "note = ?")
		args = append(args, req.Note)
	}

	if len(updates) > 0 {
		updates = append(updates, "update_time = ?")
		args = append(args, time.Now().Format(time.RFC3339), coopID)
		query := "UPDATE cooperation SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	database.AddLog("编辑", "合作记录", realName, "编辑合作记录")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *BloggerHandler) DeleteCooperation(c *gin.Context) {
	coopID, _ := strconv.Atoi(c.Param("coopId"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var title string
	h.DB.QueryRow("SELECT title FROM cooperation WHERE id = ?", coopID).Scan(&title)
	if title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合作记录不存在"})
		return
	}

	h.DB.Exec("DELETE FROM cooperation WHERE id = ?", coopID)

	database.AddLog("删除", "合作记录【"+title+"】", realName, "删除合作记录")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *BloggerHandler) BatchUpdateStatus(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		realName = "未知用户"
	}

	var req struct {
		IDs    []int  `json:"ids" binding:"required"`
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请选择要更新的博主"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	successCount := 0

	for _, id := range req.IDs {
		var oldStatus string
		h.DB.QueryRow("SELECT status FROM blogger WHERE id = ? AND is_deleted = 0", id).Scan(&oldStatus)

		if oldStatus == "" {
			oldStatus = "初次联系"
		}

		result, err := h.DB.Exec("UPDATE blogger SET status = ?, status_update_time = ?, update_time = ? WHERE id = ? AND is_deleted = 0",
			req.Status, now, now, id)
		if err != nil {
			continue
		}

		affected, _ := result.RowsAffected()
		if affected > 0 {
			successCount++

			h.DB.Exec(`INSERT INTO blogger_status_history (blogger_id, old_status, new_status, operator, remark, create_time)
				VALUES (?, ?, ?, ?, ?, ?)`,
				id, oldStatus, req.Status, realName, "批量状态更新", now)

			var nickname string
			var userID int
			h.DB.QueryRow("SELECT nickname, user_name FROM blogger WHERE id = ?", id).Scan(&nickname, &userID)
			if userID > 0 && (req.Status == "失效" || req.Status == "无效") {
				h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, blogger_id, from_user, create_time)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
					userID, "invalid_blogger", "⚠️ 博主状态变更提醒",
					"博主【"+nickname+"】的状态已变更为【"+req.Status+"】，请及时处理。",
					"important", id, realName, now)
			}
		}
	}

	database.AddLog("批量更新", "批量博主", realName, fmt.Sprintf("批量更新 %d 位博主状态为【%s】", successCount, req.Status))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("成功更新 %d 位博主", successCount),
		"data": gin.H{
			"success_count": successCount,
			"total_count":   len(req.IDs),
		},
	})
}

func (h *BloggerHandler) BatchUpdateTags(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		realName = "未知用户"
	}

	var req struct {
		IDs    []int  `json:"ids" binding:"required"`
		Tags   []int  `json:"tag_ids"`
		Action string `json:"action"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请选择要更新的博主"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	successCount := 0

	for _, bloggerID := range req.IDs {
		if req.Action == "add" && len(req.Tags) > 0 {
			for _, tagID := range req.Tags {
				var exists int
				h.DB.QueryRow("SELECT COUNT(*) FROM blogger_tags WHERE blogger_id = ? AND tag_id = ?", bloggerID, tagID).Scan(&exists)
				if exists == 0 {
					h.DB.Exec("INSERT INTO blogger_tags (blogger_id, tag_id, create_time) VALUES (?, ?, ?)",
						bloggerID, tagID, now)
				}
			}
			successCount++
		} else if req.Action == "remove" && len(req.Tags) > 0 {
			placeholders := strings.Repeat("?,", len(req.Tags))
			placeholders = placeholders[:len(placeholders)-1]
			args := make([]interface{}, 0, len(req.Tags)+1)
			args = append(args, bloggerID)
			for _, tagID := range req.Tags {
				args = append(args, tagID)
			}
			h.DB.Exec("DELETE FROM blogger_tags WHERE blogger_id = ? AND tag_id IN ("+placeholders+")", args...)
			successCount++
		} else if req.Action == "replace" {
			h.DB.Exec("DELETE FROM blogger_tags WHERE blogger_id = ?", bloggerID)
			for _, tagID := range req.Tags {
				h.DB.Exec("INSERT INTO blogger_tags (blogger_id, tag_id, create_time) VALUES (?, ?, ?)",
					bloggerID, tagID, now)
			}
			successCount++
		}
	}

	database.AddLog("批量更新", "批量博主", realName, fmt.Sprintf("批量更新 %d 位博主标签", successCount))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("成功更新 %d 位博主", successCount),
		"data": gin.H{
			"success_count": successCount,
			"total_count":   len(req.IDs),
		},
	})
}

func (h *BloggerHandler) BatchDelete(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		realName = "未知用户"
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请选择要删除的博主"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	successCount := 0

	for _, id := range req.IDs {
		var nickname string
		h.DB.QueryRow("SELECT nickname FROM blogger WHERE id = ?", id).Scan(&nickname)

		result, err := h.DB.Exec("UPDATE blogger SET is_deleted = 1, update_time = ? WHERE id = ?", now, id)
		if err != nil {
			continue
		}

		affected, _ := result.RowsAffected()
		if affected > 0 {
			successCount++
			database.AddLog("删除", nickname, realName, "批量删除博主【"+nickname+"】")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("成功删除 %d 位博主", successCount),
		"data": gin.H{
			"success_count": successCount,
			"total_count":   len(req.IDs),
		},
	})
}

func (h *BloggerHandler) SubmitEvaluation(c *gin.Context) {
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
		BloggerID int    `json:"blogger_id" binding:"required"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
		Tags      string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "博主ID不能为空"})
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		req.Rating = 5
	}

	var nickname string
	h.DB.QueryRow("SELECT nickname FROM blogger WHERE id = ? AND is_deleted = 0", req.BloggerID).Scan(&nickname)
	if nickname == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	h.DB.Exec(`INSERT OR REPLACE INTO blogger_evaluations (blogger_id, user_id, rating, comment, tags, create_time, update_time) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`, req.BloggerID, userID, req.Rating, req.Comment, req.Tags, now, now)

	database.AddLog("评价", "博主【"+nickname+"】", realName, "提交博主评价")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "评价成功"})
}

func (h *BloggerHandler) GetEvaluation(c *gin.Context) {
	bloggerIDStr := c.Param("blogger_id")
	bloggerID, _ := strconv.Atoi(bloggerIDStr)
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var evaluation struct {
		ID         int    `json:"id"`
		BloggerID  int    `json:"blogger_id"`
		UserID     int    `json:"user_id"`
		Rating     int    `json:"rating"`
		Comment    string `json:"comment"`
		Tags       string `json:"tags"`
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
	}

	err := h.DB.QueryRow(`SELECT id, blogger_id, user_id, rating, comment, tags, create_time, update_time 
		FROM blogger_evaluations WHERE blogger_id = ? AND user_id = ?`, bloggerID, userID).Scan(
		&evaluation.ID, &evaluation.BloggerID, &evaluation.UserID, &evaluation.Rating,
		&evaluation.Comment, &evaluation.Tags, &evaluation.CreateTime, &evaluation.UpdateTime)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": evaluation})
}

func (h *BloggerHandler) GetInvalidBloggers(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	query := `SELECT id, nickname, category, product, contact, wechat, platform, platform_account, 
		description, tags, avatar, user_name, real_name, status, status_update_time, status_remark, 
		create_time, update_time 
		FROM blogger WHERE is_deleted = 0 AND is_invalid = 1`

	args := []interface{}{}

	if role != "admin" {
		query += " AND user_name = (SELECT username FROM users WHERE id = ?)"
		args = append(args, userID)
	}

	query += " ORDER BY update_time DESC"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Blogger{}})
		return
	}
	defer rows.Close()

	var bloggers []models.Blogger
	for rows.Next() {
		var b models.Blogger
		rows.Scan(&b.ID, &b.Nickname, &b.Category, &b.Product, &b.Contact, &b.Wechat,
			&b.Platform, &b.PlatformAccount, &b.Description, &b.Tags, &b.Avatar,
			&b.UserName, &b.RealName, &b.Status, &b.StatusUpdateTime, &b.StatusRemark,
			&b.CreateTime, &b.UpdateTime)
		bloggers = append(bloggers, b)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": bloggers})
}

func (h *BloggerHandler) SetInvalid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var nickname string
	h.DB.QueryRow("SELECT nickname FROM blogger WHERE id = ?", id).Scan(&nickname)
	if nickname == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	h.DB.Exec("UPDATE blogger SET is_invalid = 1, update_time = ? WHERE id = ?", time.Now().Format(time.RFC3339), id)

	database.AddLog("标记失效", "博主【"+nickname+"】", realName, "标记博主为失效，原因："+req.Reason)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "标记成功"})
}

func (h *BloggerHandler) SetValid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var nickname string
	h.DB.QueryRow("SELECT nickname FROM blogger WHERE id = ?", id).Scan(&nickname)
	if nickname == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	h.DB.Exec("UPDATE blogger SET is_invalid = 0, update_time = ? WHERE id = ?", time.Now().Format(time.RFC3339), id)

	database.AddLog("标记有效", "博主【"+nickname+"】", realName, "标记博主为有效")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "标记成功"})
}

func (h *BloggerHandler) CreateTransferRequest(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	userNameVal, _ := c.Get("username")
	userName, _ := userNameVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var b models.Blogger
	err := h.DB.QueryRow(`SELECT id, nickname, user_name FROM blogger WHERE id = ? AND is_deleted = 0`, id).Scan(
		&b.ID, &b.Nickname, &b.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "博主不存在"})
		return
	}

	if b.UserName != userName {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "只能转移自己的博主"})
		return
	}

	var req struct {
		ToUserID int    `json:"to_user_id" binding:"required"`
		Reason   string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请选择接收用户"})
		return
	}

	if req.ToUserID == userID {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不能转移给自己"})
		return
	}

	var toUserName, toRealName string
	h.DB.QueryRow("SELECT username, real_name FROM users WHERE id = ? AND is_deleted = 0", req.ToUserID).Scan(&toUserName, &toRealName)
	if toUserName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "接收用户不存在"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	result, _ := h.DB.Exec(`INSERT INTO blogger_transfer_requests 
		(blogger_id, blogger_name, from_user_id, from_user_name, to_user_id, to_user_name, reason, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, b.Nickname, userID, userName, req.ToUserID, toUserName, req.Reason, now)

	transferID, _ := result.LastInsertId()

	h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, blogger_id, from_user, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ToUserID, "transfer_request", "🔄 博主转移申请",
		"【"+realName+"】申请将博主【"+b.Nickname+"】转移给您，请处理。",
		"important", id, realName, now)

	database.AddLog("转移申请", "博主【"+b.Nickname+"】", realName, "申请转移博主给【"+toRealName+"】")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "申请已提交",
		"data":    gin.H{"id": transferID},
	})
}

func (h *BloggerHandler) GetTransferRequests(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	var query string
	var args []interface{}

	if role == "admin" {
		query = `SELECT id, blogger_id, blogger_name, from_user_id, from_user_name, to_user_id, to_user_name, 
			reason, status, from_confirmed, to_confirmed, admin_confirmed, create_time 
			FROM blogger_transfer_requests ORDER BY create_time DESC`
	} else {
		query = `SELECT id, blogger_id, blogger_name, from_user_id, from_user_name, to_user_id, to_user_name, 
			reason, status, from_confirmed, to_confirmed, admin_confirmed, create_time 
			FROM blogger_transfer_requests WHERE from_user_id = ? OR to_user_id = ? ORDER BY create_time DESC`
		args = append(args, userID, userID)
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.BloggerTransferRequest{}})
		return
	}
	defer rows.Close()

	var requests []models.BloggerTransferRequest
	for rows.Next() {
		var r models.BloggerTransferRequest
		rows.Scan(&r.ID, &r.BloggerID, &r.BloggerName, &r.FromUserID, &r.FromUserName, &r.ToUserID, &r.ToUserName,
			&r.Reason, &r.Status, &r.FromConfirmed, &r.ToConfirmed, &r.AdminConfirmed, &r.CreateTime)
		requests = append(requests, r)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": requests})
}

func (h *BloggerHandler) ConfirmTransferRequest(c *gin.Context) {
	transferID, _ := strconv.Atoi(c.Param("id"))
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
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	var req struct {
		Action string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "操作类型不能为空"})
		return
	}

	var r models.BloggerTransferRequest
	err := h.DB.QueryRow(`SELECT id, blogger_id, blogger_name, from_user_id, from_user_name, to_user_id, to_user_name, 
		status, from_confirmed, to_confirmed, admin_confirmed 
		FROM blogger_transfer_requests WHERE id = ?`, transferID).Scan(
		&r.ID, &r.BloggerID, &r.BloggerName, &r.FromUserID, &r.FromUserName, &r.ToUserID, &r.ToUserName,
		&r.Status, &r.FromConfirmed, &r.ToConfirmed, &r.AdminConfirmed)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "转移申请不存在"})
		return
	}

	if r.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "申请已处理"})
		return
	}

	if req.Action == "reject" {
		h.DB.Exec("UPDATE blogger_transfer_requests SET status = 'rejected' WHERE id = ?", transferID)
		database.AddLog("拒绝转移", "博主【"+r.BloggerName+"】", realName, "拒绝博主转移申请")
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已拒绝"})
		return
	}

	if req.Action != "confirm" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效操作"})
		return
	}

	if role == "admin" {
		h.DB.Exec("UPDATE blogger_transfer_requests SET admin_confirmed = 1, status = 'completed' WHERE id = ?", transferID)
		var toRealName string
		h.DB.QueryRow("SELECT real_name FROM users WHERE id = ?", r.ToUserID).Scan(&toRealName)
		h.DB.Exec(`UPDATE blogger SET user_name = ?, real_name = ?, update_time = ? WHERE id = ?`,
			r.ToUserName, toRealName, time.Now().Format(time.RFC3339), r.BloggerID)
		database.AddLog("确认转移", "博主【"+r.BloggerName+"】", realName, "管理员确认博主转移")
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "转移成功"})
		return
	}

	if userID == r.ToUserID {
		h.DB.Exec("UPDATE blogger_transfer_requests SET to_confirmed = 1 WHERE id = ?", transferID)
		var fromConfirmed, toConfirmed int
		h.DB.QueryRow("SELECT from_confirmed, to_confirmed FROM blogger_transfer_requests WHERE id = ?", transferID).Scan(&fromConfirmed, &toConfirmed)
		if fromConfirmed == 1 && toConfirmed == 1 {
			h.DB.Exec("UPDATE blogger_transfer_requests SET status = 'completed' WHERE id = ?", transferID)
			var toRealName string
			h.DB.QueryRow("SELECT real_name FROM users WHERE id = ?", r.ToUserID).Scan(&toRealName)
			h.DB.Exec(`UPDATE blogger SET user_name = ?, real_name = ?, update_time = ? WHERE id = ?`,
				r.ToUserName, toRealName, time.Now().Format(time.RFC3339), r.BloggerID)
			database.AddLog("完成转移", "博主【"+r.BloggerName+"】", realName, "博主转移完成")
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已确认"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限操作"})
}

func (h *BloggerHandler) GetFollowupList(c *gin.Context) {
	bloggerIDStr := c.Query("blogger_id")
	if bloggerIDStr == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Followup{}})
		return
	}

	bloggerID, _ := strconv.Atoi(bloggerIDStr)
	rows, err := h.DB.Query(`SELECT id, blogger_id, content, type, operator, create_time 
		FROM followup WHERE blogger_id = ? ORDER BY create_time DESC`, bloggerID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Followup{}})
		return
	}
	defer rows.Close()

	var followups []models.Followup
	for rows.Next() {
		var f models.Followup
		rows.Scan(&f.ID, &f.BloggerID, &f.Content, &f.Type, &f.Operator, &f.CreateTime)
		followups = append(followups, f)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": followups})
}

func (h *BloggerHandler) CreateFollowupV2(c *gin.Context) {
	var req struct {
		BloggerID int    `json:"blogger_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Type      string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	followupType := req.Type
	if followupType == "" {
		followupType = "跟进"
	}

	result, err := h.DB.Exec(`INSERT INTO followup (blogger_id, content, type, operator) 
		VALUES (?, ?, ?, ?)`, req.BloggerID, req.Content, followupType, realName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	followupID, _ := result.LastInsertId()

	database.AddLog("跟进", "博主跟进记录", realName, "添加博主跟进记录")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加成功",
		"data":    gin.H{"id": followupID},
	})
}

func (h *BloggerHandler) UpdateFollowupV2(c *gin.Context) {
	followupID, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	var operator string
	h.DB.QueryRow("SELECT operator FROM followup WHERE id = ?", followupID).Scan(&operator)
	if operator == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "跟进记录不存在"})
		return
	}

	if role != "admin" && operator != realName {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限修改此跟进记录"})
		return
	}

	var req struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.Content != "" {
		updates = append(updates, "content = ?")
		args = append(args, req.Content)
	}
	if req.Type != "" {
		updates = append(updates, "type = ?")
		args = append(args, req.Type)
	}

	if len(updates) > 0 {
		args = append(args, followupID)
		query := "UPDATE followup SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *BloggerHandler) DeleteFollowupV2(c *gin.Context) {
	followupID, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	var operator string
	h.DB.QueryRow("SELECT operator FROM followup WHERE id = ?", followupID).Scan(&operator)
	if operator == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "跟进记录不存在"})
		return
	}

	if role != "admin" && operator != realName {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此跟进记录"})
		return
	}

	h.DB.Exec("DELETE FROM followup WHERE id = ?", followupID)

	database.AddLog("删除", "跟进记录", realName, "删除跟进记录")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *BloggerHandler) GetCooperationList(c *gin.Context) {
	bloggerIDStr := c.Query("blogger_id")
	if bloggerIDStr == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Cooperation{}})
		return
	}

	bloggerID, _ := strconv.Atoi(bloggerIDStr)
	rows, err := h.DB.Query(`SELECT id, blogger_id, title, date, status, product, amount, note, create_time, update_time 
		FROM cooperation WHERE blogger_id = ? ORDER BY create_time DESC`, bloggerID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Cooperation{}})
		return
	}
	defer rows.Close()

	var cooperations []models.Cooperation
	for rows.Next() {
		var coop models.Cooperation
		rows.Scan(&coop.ID, &coop.BloggerID, &coop.Title, &coop.Date, &coop.Status, &coop.Product, &coop.Amount, &coop.Note, &coop.CreateTime, &coop.UpdateTime)
		cooperations = append(cooperations, coop)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": cooperations})
}

func (h *BloggerHandler) CreateCooperationV2(c *gin.Context) {
	var req struct {
		BloggerID int     `json:"blogger_id" binding:"required"`
		Title     string  `json:"title" binding:"required"`
		Date      string  `json:"date"`
		Status    string  `json:"status"`
		Product   string  `json:"product"`
		Amount    float64 `json:"amount"`
		Note      string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "标题不能为空"})
		return
	}

	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	status := req.Status
	if status == "" {
		status = "进行中"
	}

	now := time.Now().Format(time.RFC3339)
	result, err := h.DB.Exec(`INSERT INTO cooperation (blogger_id, title, date, status, product, amount, note, create_time, update_time) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, req.BloggerID, req.Title, req.Date, status, req.Product, req.Amount, req.Note, now, now)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	coopID, _ := result.LastInsertId()

	database.AddLog("合作", "合作记录", realName, "添加合作记录【"+req.Title+"】")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "添加成功",
		"data":    gin.H{"id": coopID},
	})
}

func (h *BloggerHandler) UpdateCooperationV2(c *gin.Context) {
	coopID, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var exists int
	h.DB.QueryRow("SELECT COUNT(*) FROM cooperation WHERE id = ?", coopID).Scan(&exists)
	if exists == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合作记录不存在"})
		return
	}

	var req struct {
		Title   string  `json:"title"`
		Date    string  `json:"date"`
		Status  string  `json:"status"`
		Product string  `json:"product"`
		Amount  float64 `json:"amount"`
		Note    string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.Title != "" {
		updates = append(updates, "title = ?")
		args = append(args, req.Title)
	}
	if req.Date != "" {
		updates = append(updates, "date = ?")
		args = append(args, req.Date)
	}
	if req.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, req.Status)
	}
	if req.Product != "" {
		updates = append(updates, "product = ?")
		args = append(args, req.Product)
	}
	if req.Amount != 0 {
		updates = append(updates, "amount = ?")
		args = append(args, req.Amount)
	}
	if req.Note != "" {
		updates = append(updates, "note = ?")
		args = append(args, req.Note)
	}

	if len(updates) > 0 {
		updates = append(updates, "update_time = ?")
		args = append(args, time.Now().Format(time.RFC3339), coopID)
		query := "UPDATE cooperation SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	database.AddLog("编辑", "合作记录", realName, "编辑合作记录")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *BloggerHandler) DeleteCooperationV2(c *gin.Context) {
	coopID, _ := strconv.Atoi(c.Param("id"))
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var title string
	h.DB.QueryRow("SELECT title FROM cooperation WHERE id = ?", coopID).Scan(&title)
	if title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "合作记录不存在"})
		return
	}

	h.DB.Exec("DELETE FROM cooperation WHERE id = ?", coopID)

	database.AddLog("删除", "合作记录【"+title+"】", realName, "删除合作记录")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
