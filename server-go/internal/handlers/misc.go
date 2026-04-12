package handlers

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server-go/internal/config"
	"server-go/internal/database"
	"server-go/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MiscHandler struct {
	DB  *sql.DB
	Cfg *config.Config
}

func NewMiscHandler(db *sql.DB, cfg *config.Config) *MiscHandler {
	return &MiscHandler{DB: db, Cfg: cfg}
}

func (h *MiscHandler) GetCategories(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, name, color, icon, create_time, update_time FROM categories ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Category{}})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		rows.Scan(&cat.ID, &cat.Name, &cat.Color, &cat.Icon, &cat.CreateTime, &cat.UpdateTime)
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": categories})
}

func (h *MiscHandler) GetPlatforms(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, name, icon, create_time FROM platforms ORDER BY id`)
	if err != nil {
		platforms := []map[string]interface{}{
			{"id": 1, "name": "小红书"},
			{"id": 2, "name": "抖音"},
			{"id": 3, "name": "微博"},
			{"id": 4, "name": "B站"},
			{"id": 5, "name": "快手"},
			{"id": 6, "name": "微信公众号"},
			{"id": 7, "name": "视频号"},
			{"id": 8, "name": "知乎"},
			{"id": 9, "name": "淘宝直播"},
			{"id": 10, "name": "其他"},
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": platforms})
		return
	}
	defer rows.Close()

	var platforms []models.Platform
	for rows.Next() {
		var p models.Platform
		rows.Scan(&p.ID, &p.Name, &p.Icon, &p.CreateTime)
		platforms = append(platforms, p)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": platforms})
}

func (h *MiscHandler) GetProducts(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, name, create_time FROM products ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Product{}})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.Name, &p.CreateTime)
		products = append(products, p)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": products})
}

func (h *MiscHandler) AddCategory(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
		Icon  string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "分类名称不能为空"})
		return
	}

	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE name = ?", req.Name).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "分类已存在"})
		return
	}

	color := req.Color
	if color == "" {
		color = "#6b7280"
	}
	icon := req.Icon
	if icon == "" {
		icon = "🏷️"
	}

	h.DB.Exec("INSERT INTO categories (name, color, icon) VALUES (?, ?, ?)", req.Name, color, icon)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	database.AddLog("新增", "分类【"+req.Name+"】", realName, "新增分类【"+req.Name+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "添加成功"})
}

func (h *MiscHandler) UpdateCategory(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
		Icon  string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var oldName string
	h.DB.QueryRow("SELECT name FROM categories WHERE id = ?", req.ID).Scan(&oldName)
	if oldName == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "分类不存在"})
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
	if req.Icon != "" {
		updates = append(updates, "icon = ?")
		args = append(args, req.Icon)
	}

	if len(updates) > 0 {
		updates = append(updates, "update_time = CURRENT_TIMESTAMP")
		args = append(args, req.ID)
		query := "UPDATE categories SET " + strings.Join(updates, ", ") + " WHERE id = ?"
		h.DB.Exec(query, args...)
	}

	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	newName := req.Name
	if newName == "" {
		newName = oldName
	}
	database.AddLog("编辑", "分类【"+oldName+"】", realName, "编辑分类【"+oldName+"】→【"+newName+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *MiscHandler) DeleteCategory(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var name string
	h.DB.QueryRow("SELECT name FROM categories WHERE id = ?", req.ID).Scan(&name)
	if name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "分类不存在"})
		return
	}

	h.DB.Exec("DELETE FROM categories WHERE id = ?", req.ID)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	database.AddLog("删除", "分类【"+name+"】", realName, "删除分类【"+name+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *MiscHandler) AddProduct(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "产品名称不能为空"})
		return
	}

	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM products WHERE name = ?", req.Name).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "产品已存在"})
		return
	}

	h.DB.Exec("INSERT INTO products (name) VALUES (?)", req.Name)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	database.AddLog("新增", "产品【"+req.Name+"】", realName, "新增产品【"+req.Name+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "添加成功"})
}

func (h *MiscHandler) UpdateProduct(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id := c.Param("id")
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "产品名称不能为空"})
		return
	}

	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", id).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	h.DB.Exec("UPDATE products SET name = ? WHERE id = ?", req.Name, id)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	database.AddLog("更新", "产品【"+req.Name+"】", realName, "更新产品【"+req.Name+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *MiscHandler) DeleteProduct(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	id := c.Param("id")

	var name string
	h.DB.QueryRow("SELECT name FROM products WHERE id = ?", id).Scan(&name)
	if name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	h.DB.Exec("DELETE FROM products WHERE id = ?", id)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	database.AddLog("删除", "产品【"+name+"】", realName, "删除产品【"+name+"】")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *MiscHandler) GetBloggerStat(c *gin.Context) {
	var total int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE is_deleted = 0").Scan(&total)

	today := time.Now()
	todayStr := today.Format("2006-01-02")
	var todayCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE is_deleted = 0 AND DATE(create_time) = ?", todayStr).Scan(&todayCount)

	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.Local)
	monthStartStr := monthStart.Format("2006-01-02")
	var monthCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE is_deleted = 0 AND DATE(create_time) >= ?", monthStartStr).Scan(&monthCount)

	rows, _ := h.DB.Query("SELECT DISTINCT user_name FROM blogger WHERE is_deleted = 0")
	var users int
	for rows.Next() {
		users++
	}
	rows.Close()

	byUserMap := make(map[string]int)
	rows, _ = h.DB.Query("SELECT real_name FROM blogger WHERE is_deleted = 0")
	for rows.Next() {
		var realName string
		rows.Scan(&realName)
		byUserMap[realName]++
	}
	rows.Close()

	byUser := []map[string]interface{}{}
	for name, count := range byUserMap {
		percentage := 0.0
		if total > 0 {
			percentage = float64(count) / float64(total) * 100
		}
		byUser = append(byUser, map[string]interface{}{
			"real_name":  name,
			"count":      count,
			"percentage": fmt.Sprintf("%.1f%%", percentage),
		})
	}

	byCategoryMap := make(map[string]int)
	rows, _ = h.DB.Query("SELECT category FROM blogger WHERE is_deleted = 0")
	for rows.Next() {
		var category string
		rows.Scan(&category)
		if category != "" {
			byCategoryMap[category]++
		}
	}
	rows.Close()

	byCategory := []map[string]interface{}{}
	for cat, count := range byCategoryMap {
		percentage := 0.0
		if total > 0 {
			percentage = float64(count) / float64(total) * 100
		}
		byCategory = append(byCategory, map[string]interface{}{
			"category":   cat,
			"count":      count,
			"percentage": fmt.Sprintf("%.1f%%", percentage),
		})
	}

	byPlatformMap := make(map[string]int)
	rows, _ = h.DB.Query("SELECT platform FROM blogger WHERE is_deleted = 0 AND platform != ''")
	for rows.Next() {
		var platform string
		rows.Scan(&platform)
		byPlatformMap[platform]++
	}
	rows.Close()

	byPlatform := []map[string]interface{}{}
	for plat, count := range byPlatformMap {
		percentage := 0.0
		if total > 0 {
			percentage = float64(count) / float64(total) * 100
		}
		byPlatform = append(byPlatform, map[string]interface{}{
			"platform":   plat,
			"count":      count,
			"percentage": fmt.Sprintf("%.1f%%", percentage),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"stats": gin.H{
				"total": total,
				"today": todayCount,
				"month": monthCount,
				"users": users,
			},
			"byUser":     byUser,
			"byCategory": byCategory,
			"byPlatform": byPlatform,
		},
	})
}

func (h *MiscHandler) GetBloggerCharts(c *gin.Context) {
	byPlatformMap := make(map[string]int)
	rows, _ := h.DB.Query("SELECT platform FROM blogger WHERE is_deleted = 0 AND platform != ''")
	for rows.Next() {
		var platform string
		rows.Scan(&platform)
		byPlatformMap[platform]++
	}
	rows.Close()

	byPlatform := []map[string]interface{}{}
	for name, value := range byPlatformMap {
		byPlatform = append(byPlatform, map[string]interface{}{"name": name, "value": value})
	}

	byCategoryMap := make(map[string]int)
	rows, _ = h.DB.Query("SELECT category FROM blogger WHERE is_deleted = 0")
	for rows.Next() {
		var category string
		rows.Scan(&category)
		if category != "" {
			byCategoryMap[category]++
		}
	}
	rows.Close()

	byCategory := []map[string]interface{}{}
	for name, value := range byCategoryMap {
		byCategory = append(byCategory, map[string]interface{}{"name": name, "value": value})
	}

	byStatusMap := make(map[string]int)
	rows, _ = h.DB.Query("SELECT status FROM blogger WHERE is_deleted = 0")
	for rows.Next() {
		var status string
		rows.Scan(&status)
		if status == "" {
			status = "初次联系"
		}
		byStatusMap[status]++
	}
	rows.Close()

	byStatus := []map[string]interface{}{}
	for name, value := range byStatusMap {
		byStatus = append(byStatus, map[string]interface{}{"name": name, "value": value})
	}

	byUserMap := make(map[string]int)
	rows, _ = h.DB.Query("SELECT real_name FROM blogger WHERE is_deleted = 0")
	for rows.Next() {
		var realName string
		rows.Scan(&realName)
		byUserMap[realName]++
	}
	rows.Close()

	byUser := []map[string]interface{}{}
	for name, value := range byUserMap {
		byUser = append(byUser, map[string]interface{}{"name": name, "value": value})
	}

	monthlyMap := make(map[string]int)
	rows, _ = h.DB.Query("SELECT strftime('%Y-%m', create_time) as month FROM blogger WHERE is_deleted = 0")
	for rows.Next() {
		var month string
		rows.Scan(&month)
		monthlyMap[month]++
	}
	rows.Close()

	monthly := []map[string]interface{}{}
	for name, value := range monthlyMap {
		monthly = append(monthly, map[string]interface{}{"name": name, "value": value})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"byPlatform": byPlatform,
			"byCategory": byCategory,
			"byStatus":   byStatus,
			"byUser":     byUser,
			"monthly":    monthly,
		},
	})
}

func (h *MiscHandler) GetBloggerStatusList(c *gin.Context) {
	statusList := []map[string]string{
		{"id": "初次联系", "name": "初次联系"},
		{"id": "洽谈中", "name": "洽谈中"},
		{"id": "已合作", "name": "已合作"},
		{"id": "已拒绝", "name": "已拒绝"},
		{"id": "暂停合作", "name": "暂停合作"},
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": statusList})
}

func (h *MiscHandler) GetInvalidBloggerCount(c *gin.Context) {
	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM blogger WHERE is_deleted = 0 AND is_invalid = 1").Scan(&count)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": count})
}

func (h *MiscHandler) GetPublicUsers(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	testUsernames := map[string]bool{"test": true, "testuser": true}

	rows, err := h.DB.Query(`SELECT id, username, real_name, avatar, bio, role, status, team_id, create_time, approved_by, approved_time 
		FROM users WHERE is_deleted = 0`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []map[string]interface{}{}})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, realName, avatar, bio, userRole, status sql.NullString
		var teamID sql.NullInt64
		var createTime, approvedBy, approvedTime sql.NullString
		rows.Scan(&id, &username, &realName, &avatar, &bio, &userRole, &status, &teamID, &createTime, &approvedBy, &approvedTime)

		if testUsernames[strings.ToLower(username.String)] {
			continue
		}

		if role != "admin" && userRole.String == "admin" {
			continue
		}

		user := map[string]interface{}{
			"id":            id,
			"username":      username.String,
			"real_name":     realName.String,
			"avatar":        avatar.String,
			"bio":           bio.String,
			"role":          userRole.String,
			"status":        status.String,
			"team_id":       nil,
			"create_time":   createTime.String,
			"approved_by":   approvedBy.String,
			"approved_time": approvedTime.String,
		}
		if teamID.Valid {
			user["team_id"] = teamID.Int64
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": users})
}

func (h *MiscHandler) GetOperationLog(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	var total int
	h.DB.QueryRow("SELECT COUNT(*) FROM operation_log").Scan(&total)

	rows, _ := h.DB.Query(`SELECT id, action, target, operator, detail, create_time 
		FROM operation_log ORDER BY create_time DESC LIMIT ? OFFSET ?`, pageSize, offset)

	var logs []models.OperationLog
	for rows.Next() {
		var l models.OperationLog
		rows.Scan(&l.ID, &l.Action, &l.Target, &l.Operator, &l.Detail, &l.CreateTime)
		logs = append(logs, l)
	}
	rows.Close()

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

func (h *MiscHandler) GetBloggerHotData(c *gin.Context) {
	_ = c.Param("nickname")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"trend":     []interface{}{},
			"followers": []interface{}{},
			"message":   "数据暂未接入",
		},
	})
}

func (h *MiscHandler) GetBloggerNews(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": []interface{}{},
	})
}

func (h *MiscHandler) GetNotifications(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	rows, err := h.DB.Query(`SELECT id, user_id, type, title, content, priority, team_id, from_user, blogger_id, post_id, is_read, create_time 
		FROM notifications WHERE user_id = ? ORDER BY create_time DESC`, userID)

	if err != nil || rows == nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []map[string]interface{}{}})
		return
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var id, userID, isRead int
		var type_, title, content, priority, fromUser sql.NullString
		var teamID, bloggerID, postID sql.NullInt64
		var createTime string
		rows.Scan(&id, &userID, &type_, &title, &content, &priority, &teamID, &fromUser, &bloggerID, &postID, &isRead, &createTime)

		notif := map[string]interface{}{
			"id":          id,
			"user_id":     userID,
			"type":        type_.String,
			"title":       title.String,
			"content":     content.String,
			"priority":    priority.String,
			"from_user":   fromUser.String,
			"is_read":     isRead == 1,
			"read":        isRead == 1,
			"create_time": createTime,
		}
		if teamID.Valid {
			notif["team_id"] = teamID.Int64
		}
		if bloggerID.Valid {
			notif["blogger_id"] = bloggerID.Int64
		}
		if postID.Valid {
			notif["post_id"] = postID.Int64
		}
		notifications = append(notifications, notif)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": notifications})
}

func (h *MiscHandler) MarkNotificationRead(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	h.DB.Exec("UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?", id, userID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已标记为已读"})
}

func (h *MiscHandler) MarkAllNotificationsRead(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	h.DB.Exec("UPDATE notifications SET is_read = 1 WHERE user_id = ?", userID)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "已全部标记为已读"})
}

func (h *MiscHandler) DeleteNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	result, err := h.DB.Exec("DELETE FROM notifications WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "通知不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *MiscHandler) BatchDeleteNotifications(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	tx, err := h.DB.Begin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	stmt, err := tx.Prepare("DELETE FROM notifications WHERE id = ? AND user_id = ?")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	defer stmt.Close()
	for _, id := range req.IDs {
		_, err = stmt.Exec(id, userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
			return
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *MiscHandler) GetAnnouncements(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, title, content, create_time FROM announcements ORDER BY create_time DESC`)

	if err != nil || rows == nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []map[string]interface{}{}})
		return
	}
	defer rows.Close()

	var announcements []map[string]interface{}
	for rows.Next() {
		var id int
		var title, content, createTime string
		rows.Scan(&id, &title, &content, &createTime)
		announcements = append(announcements, map[string]interface{}{
			"id":          id,
			"title":       title,
			"content":     content,
			"create_time": createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": announcements})
}

func (h *MiscHandler) CreateAnnouncement(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "标题不能为空"})
		return
	}

	h.DB.Exec("INSERT INTO announcements (title, content) VALUES (?, ?)", req.Title, req.Content)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功"})
}

func (h *MiscHandler) UploadBloggerAvatar(c *gin.Context) {
	log.Printf("[UploadBloggerAvatar] Request received, Content-Type: %s", c.GetHeader("Content-Type"))

	var req struct {
		Image string `json:"image" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UploadBloggerAvatar] JSON bind error: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	idx := strings.Index(req.Image, ",")
	if idx == -1 {
		log.Printf("[UploadBloggerAvatar] Invalid image format - no comma found")
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Invalid image format"})
		return
	}

	allowedTypes := map[string]bool{
		"data:image/jpeg": true,
		"data:image/png":  true,
		"data:image/gif":  true,
		"data:image/webp": true,
	}
	validType := false
	for prefix := range allowedTypes {
		if strings.HasPrefix(req.Image, prefix) {
			validType = true
			break
		}
	}
	if !validType {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "仅支持 JPG、PNG、GIF、WebP 格式的图片"})
		return
	}

	imageData := req.Image[idx+1:]
	log.Printf("[UploadBloggerAvatar] Base64 data length: %d", len(imageData))

	if len(req.Image) > 7*1024*1024 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "文件大小不能超过5MB"})
		return
	}

	ext := ".jpg"
	if strings.HasPrefix(req.Image, "data:image/png") {
		ext = ".png"
	} else if strings.HasPrefix(req.Image, "data:image/gif") {
		ext = ".gif"
	} else if strings.HasPrefix(req.Image, "data:image/webp") {
		ext = ".webp"
	}

	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		log.Printf("[UploadBloggerAvatar] Base64 decode error: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Failed to decode image"})
		return
	}

	uploadDir := filepath.Join(h.Cfg.UploadPath, "avatars")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("[UploadBloggerAvatar] Failed to create directory: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to create upload directory"})
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().UnixMilli(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
		log.Printf("[UploadBloggerAvatar] Failed to save file: %v", err)
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to save file"})
		return
	}

	log.Printf("[UploadBloggerAvatar] File saved successfully: %s", filePath)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"url": "/uploads/avatars/" + filename,
		},
		"message": "Avatar uploaded successfully",
	})
}

func (h *MiscHandler) UploadUserAvatar(c *gin.Context) {
	log.Printf("[UploadUserAvatar] Request received")

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择图片"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "图片大小不能超过5MB"})
		return
	}

	contentType := file.Header.Get("Content-Type")
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if contentType != "" && !allowedContentTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "仅支持 JPG、PNG、GIF、WebP 格式的图片"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只支持JPG、PNG、GIF、WEBP格式"})
		return
	}

	uploadDir := filepath.Join(h.Cfg.UploadPath, "user-avatars")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建上传目录失败"})
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存文件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"url": "/uploads/user-avatars/" + filename,
		},
	})
}

func (h *MiscHandler) UploadTeamLogo(c *gin.Context) {
	var req struct {
		Image string `json:"image"`
	}

	var filename string
	var uploadDir string = filepath.Join(h.Cfg.UploadPath, "team-logos")

	if err := c.ShouldBindJSON(&req); err == nil && req.Image != "" {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建上传目录失败"})
			return
		}

		dataParts := strings.SplitN(req.Image, ",", 2)
		if len(dataParts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的图片格式"})
			return
		}

		mimeType := strings.Split(dataParts[0], ";")[0]
		var ext string
		switch mimeType {
		case "data:image/jpeg":
			ext = ".jpg"
		case "data:image/png":
			ext = ".png"
		case "data:image/gif":
			ext = ".gif"
		case "data:image/webp":
			ext = ".webp"
		default:
			ext = ".png"
		}

		filename = fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		data, err := base64.StdEncoding.DecodeString(dataParts[1])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解码图片失败"})
			return
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存文件失败"})
			return
		}
	} else {
		file, err := c.FormFile("logo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择图片"})
			return
		}

		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "图片大小不能超过5MB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只支持JPG、PNG、GIF、WEBP格式"})
			return
		}

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建上传目录失败"})
			return
		}

		filename = fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存文件失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"url": "/uploads/team-logos/" + filename,
		},
	})
}

func (h *MiscHandler) UploadTeamBgImage(c *gin.Context) {
	var req struct {
		Image string `json:"image"`
	}

	var filename string
	var uploadDir string = filepath.Join(h.Cfg.UploadPath, "team-bgimages")

	if err := c.ShouldBindJSON(&req); err == nil && req.Image != "" {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建上传目录失败"})
			return
		}

		dataParts := strings.SplitN(req.Image, ",", 2)
		if len(dataParts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的图片格式"})
			return
		}

		mimeType := strings.Split(dataParts[0], ";")[0]
		var ext string
		switch mimeType {
		case "data:image/jpeg":
			ext = ".jpg"
		case "data:image/png":
			ext = ".png"
		case "data:image/gif":
			ext = ".gif"
		case "data:image/webp":
			ext = ".webp"
		default:
			ext = ".png"
		}

		filename = fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		data, err := base64.StdEncoding.DecodeString(dataParts[1])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解码图片失败"})
			return
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存文件失败"})
			return
		}
	} else {
		file, err := c.FormFile("bgImage")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择图片"})
			return
		}

		if file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "图片大小不能超过10MB"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只支持JPG、PNG、GIF、WEBP格式"})
			return
		}

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建上传目录失败"})
			return
		}

		filename = fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存文件失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"url": "/uploads/team-bgimages/" + filename,
		},
	})
}

func (h *MiscHandler) GetSnapshots(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	rows, err := h.DB.Query(`SELECT id, name, filename, size, created, mtime, create_time 
		FROM snapshots ORDER BY create_time DESC`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Snapshot{}})
		return
	}
	defer rows.Close()

	var snapshots []models.Snapshot
	for rows.Next() {
		var s models.Snapshot
		rows.Scan(&s.ID, &s.Name, &s.Filename, &s.Size, &s.Created, &s.Mtime, &s.CreateTime)
		snapshots = append(snapshots, s)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": snapshots})
}

func (h *MiscHandler) CreateSnapshot(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "名称不能为空"})
		return
	}

	now := time.Now()
	filename := fmt.Sprintf("snapshot_%s.db", now.Format("20060102_150405"))
	snapshotPath := filepath.Join(h.Cfg.UploadPath, "snapshots")
	if err := os.MkdirAll(snapshotPath, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建快照目录失败"})
		return
	}

	dbPath := h.Cfg.DatabasePath
	destPath := filepath.Join(snapshotPath, filename)

	sourceData, err := os.ReadFile(dbPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "读取数据库失败"})
		return
	}

	if err := os.WriteFile(destPath, sourceData, 0644); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "保存快照失败"})
		return
	}

	fileInfo, _ := os.Stat(destPath)
	size := int(fileInfo.Size())
	created := now.Format(time.RFC3339)
	mtime := now.Format(time.RFC3339)

	result, _ := h.DB.Exec(`INSERT INTO snapshots (name, filename, size, created, mtime) 
		VALUES (?, ?, ?, ?, ?)`, req.Name, filename, size, created, mtime)

	snapshotID, _ := result.LastInsertId()

	database.AddLog("快照", "创建快照【"+req.Name+"】", realName, "创建数据库快照")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "快照创建成功",
		"data":    gin.H{"id": snapshotID},
	})
}

func (h *MiscHandler) DeleteSnapshot(c *gin.Context) {
	idOrFilename := c.Param("id")
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var filename string
	var snapshotID int
	if idInt, err := strconv.Atoi(idOrFilename); err == nil {
		snapshotID = idInt
		h.DB.QueryRow("SELECT filename FROM snapshots WHERE id = ?", snapshotID).Scan(&filename)
	} else {
		filename = idOrFilename
		h.DB.QueryRow("SELECT id FROM snapshots WHERE filename = ?", filename).Scan(&snapshotID)
	}

	if filename == "" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "快照不存在"})
		return
	}

	snapshotPath := filepath.Join(h.Cfg.UploadPath, "snapshots", filename)
	absPath, _ := filepath.Abs(snapshotPath)
	snapshotDir, _ := filepath.Abs(filepath.Join(h.Cfg.UploadPath, "snapshots"))
	if !strings.HasPrefix(absPath, snapshotDir+string(filepath.Separator)) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "非法文件路径"})
		return
	}

	os.Remove(snapshotPath)

	h.DB.Exec("DELETE FROM snapshots WHERE filename = ?", filename)

	database.AddLog("删除快照", "快照:"+filename, realName, "删除数据库快照")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *MiscHandler) GetSnapshotSettings(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var autoBackup bool
	var backupInterval int
	var keepBackups int

	var autoBackupStr string
	err := h.DB.QueryRow(`SELECT value FROM system_settings WHERE key = 'auto_backup'`).Scan(&autoBackupStr)
	if err != nil {
		autoBackup = true
	} else {
		autoBackup = autoBackupStr == "true"
	}

	var backupIntervalStr string
	err = h.DB.QueryRow(`SELECT value FROM system_settings WHERE key = 'backup_interval'`).Scan(&backupIntervalStr)
	if err != nil {
		backupInterval = 24
	} else {
		backupInterval, _ = strconv.Atoi(backupIntervalStr)
		if backupInterval == 0 {
			backupInterval = 24
		}
	}

	var keepBackupsStr string
	err = h.DB.QueryRow(`SELECT value FROM system_settings WHERE key = 'keep_backups'`).Scan(&keepBackupsStr)
	if err != nil {
		keepBackups = 7
	} else {
		keepBackups, _ = strconv.Atoi(keepBackupsStr)
		if keepBackups == 0 {
			keepBackups = 7
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"auto_backup":     autoBackup,
			"backup_interval": backupInterval,
			"keep_backups":    keepBackups,
		},
	})
}

func (h *MiscHandler) SetSnapshotSettings(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		AutoBackup     bool `json:"auto_backup"`
		BackupInterval int  `json:"backup_interval"`
		KeepBackups    int  `json:"keep_backups"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	h.DB.Exec(`INSERT OR REPLACE INTO system_settings (key, value) VALUES (?, ?)`, "auto_backup", req.AutoBackup)
	h.DB.Exec(`INSERT OR REPLACE INTO system_settings (key, value) VALUES (?, ?)`, "backup_interval", req.BackupInterval)
	h.DB.Exec(`INSERT OR REPLACE INTO system_settings (key, value) VALUES (?, ?)`, "keep_backups", req.KeepBackups)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "设置保存成功"})
}

func (h *MiscHandler) CreateBackup(c *gin.Context) {
	h.CreateSnapshot(c)
}

func (h *MiscHandler) PurgeData(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	var req struct {
		Confirm bool `json:"confirm"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || !req.Confirm {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "需要确认操作"})
		return
	}

	now := time.Now()
	h.DB.Exec("DELETE FROM operation_log WHERE create_time < ?", now.AddDate(0, -3, 0))
	h.DB.Exec("DELETE FROM notifications WHERE is_read = 1 AND create_time < ?", now.AddDate(0, -1, 0))

	database.AddLog("数据清理", "清理旧数据", realName, "清理3个月前的操作日志和1个月前的已读通知")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "数据清理完成"})
}

func sanitizeSnapshotFilename(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("文件名不能为空")
	}
	if strings.Contains(filename, "..") || strings.ContainsAny(filename, "/\\") {
		return "", fmt.Errorf("非法文件名")
	}
	if !strings.HasPrefix(filename, "snapshot_") || !strings.HasSuffix(filename, ".db") {
		return "", fmt.Errorf("非法快照文件名")
	}
	return filename, nil
}

func (h *MiscHandler) DownloadSnapshot(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	filename := c.Param("filename")
	filename, err := sanitizeSnapshotFilename(filename)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	snapshotPath := filepath.Join(h.Cfg.UploadPath, "snapshots", filename)
	absPath, _ := filepath.Abs(snapshotPath)
	snapshotDir, _ := filepath.Abs(filepath.Join(h.Cfg.UploadPath, "snapshots"))
	if !strings.HasPrefix(absPath, snapshotDir+string(filepath.Separator)) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "非法文件路径"})
		return
	}

	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "快照不存在"})
		return
	}

	c.FileAttachment(snapshotPath, filename)
}

func (h *MiscHandler) RestoreSnapshot(c *gin.Context) {
	filename := c.Param("filename")
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	filename, err := sanitizeSnapshotFilename(filename)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	snapshotPath := filepath.Join(h.Cfg.UploadPath, "snapshots", filename)
	absPath, _ := filepath.Abs(snapshotPath)
	snapshotDir, _ := filepath.Abs(filepath.Join(h.Cfg.UploadPath, "snapshots"))
	if !strings.HasPrefix(absPath, snapshotDir+string(filepath.Separator)) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "非法文件路径"})
		return
	}

	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "快照不存在"})
		return
	}

	sourceData, err := os.ReadFile(snapshotPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "读取快照失败"})
		return
	}

	if err := os.WriteFile(h.Cfg.DatabasePath, sourceData, 0644); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "还原数据库失败"})
		return
	}

	database.AddLog("快照", "还原快照【"+filename+"】", realName, "还原数据库快照")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "数据库已还原，请重启服务以使更改生效"})
}

func (h *MiscHandler) ClearData(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	realNameVal, _ := c.Get("realName")
	realName, ok := realNameVal.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}
	if realName == "" {
		realName = "未知用户"
	}

	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "开启事务失败"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	statements := []string{
		"DELETE FROM blogger",
		"DELETE FROM followup",
		"DELETE FROM cooperation",
		"DELETE FROM operation_log",
		"DELETE FROM notifications",
		"DELETE FROM teams",
		"DELETE FROM team_posts",
		"DELETE FROM team_post_comments",
		"DELETE FROM team_post_likes",
		"DELETE FROM team_post_collects",
		"DELETE FROM team_messages",
		"DELETE FROM public_posts",
		"DELETE FROM public_post_comments",
		"DELETE FROM public_post_likes",
		"DELETE FROM public_post_collects",
		"DELETE FROM announcements",
		"DELETE FROM blogger_transfer_requests",
		"DELETE FROM blogger_evaluations",
		"DELETE FROM blogger_status_history",
		"DELETE FROM private_messages",
		"DELETE FROM snapshots",
		"UPDATE users SET team_id = NULL, team_name = NULL, team_color = NULL",
	}

	for _, stmt := range statements {
		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"code": 500, "message": "数据清空失败"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "提交事务失败"})
		return
	}

	database.AddLog("清空数据", "清空所有数据", realName, "清空博主、跟进，合作、团队、帖子、通知等所有业务数据")

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "数据清空完成"})
}

func (h *MiscHandler) GetTeamBloggerStat(c *gin.Context) {
	teamIDStr := c.Query("team_id")
	if teamIDStr == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"stats": gin.H{}}})
		return
	}

	teamID, _ := strconv.Atoi(teamIDStr)

	var total int
	h.DB.QueryRow(`SELECT COUNT(*) FROM blogger b 
		INNER JOIN users u ON b.user_name = u.username 
		WHERE b.is_deleted = 0 AND u.team_id = ?`, teamID).Scan(&total)

	today := time.Now().Format("2006-01-02")
	var todayCount int
	h.DB.QueryRow(`SELECT COUNT(*) FROM blogger b 
		INNER JOIN users u ON b.user_name = u.username 
		WHERE b.is_deleted = 0 AND u.team_id = ? AND DATE(b.create_time) = ?`, teamID, today).Scan(&todayCount)

	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	var monthCount int
	h.DB.QueryRow(`SELECT COUNT(*) FROM blogger b 
		INNER JOIN users u ON b.user_name = u.username 
		WHERE b.is_deleted = 0 AND u.team_id = ? AND DATE(b.create_time) >= ?`, teamID, monthStart).Scan(&monthCount)

	rows, _ := h.DB.Query(`SELECT DISTINCT b.user_name FROM blogger b 
		INNER JOIN users u ON b.user_name = u.username 
		WHERE b.is_deleted = 0 AND u.team_id = ?`, teamID)
	var users int
	for rows.Next() {
		users++
	}
	rows.Close()

	byUserMap := make(map[string]int)
	rows, _ = h.DB.Query(`SELECT u.real_name FROM blogger b 
		INNER JOIN users u ON b.user_name = u.username 
		WHERE b.is_deleted = 0 AND u.team_id = ?`, teamID)
	for rows.Next() {
		var realName string
		rows.Scan(&realName)
		byUserMap[realName]++
	}
	rows.Close()

	byUser := []map[string]interface{}{}
	for name, count := range byUserMap {
		byUser = append(byUser, map[string]interface{}{
			"real_name": name,
			"count":     count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"stats": gin.H{
				"total": total,
				"today": todayCount,
				"month": monthCount,
				"users": users,
			},
			"byUser": byUser,
		},
	})
}

func (h *MiscHandler) GetTeamBloggerCharts(c *gin.Context) {
	teamIDStr := c.Query("team_id")
	if teamIDStr == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{}})
		return
	}

	teamID, _ := strconv.Atoi(teamIDStr)

	byUserMap := make(map[string]int)
	rows, _ := h.DB.Query(`SELECT u.real_name FROM blogger b 
		INNER JOIN users u ON b.user_name = u.username 
		WHERE b.is_deleted = 0 AND u.team_id = ?`, teamID)
	for rows.Next() {
		var realName string
		rows.Scan(&realName)
		byUserMap[realName]++
	}
	rows.Close()

	byMember := []map[string]interface{}{}
	bloggerList := []map[string]interface{}{}

	for name, count := range byUserMap {
		var userBloggers []map[string]interface{}
		bloggerRows, _ := h.DB.Query(`SELECT b.id, b.nickname, b.category, b.platform, b.status, b.create_time 
			FROM blogger b 
			INNER JOIN users u ON b.user_name = u.username 
			WHERE b.is_deleted = 0 AND u.team_id = ? AND u.real_name = ? 
			ORDER BY b.create_time DESC`, teamID, name)

		for bloggerRows.Next() {
			var id int
			var nickname, category, platform, status, createTime string
			bloggerRows.Scan(&id, &nickname, &category, &platform, &status, &createTime)
			blogger := map[string]interface{}{
				"id":          id,
				"nickname":    nickname,
				"category":    category,
				"platform":    platform,
				"status":      status,
				"create_time": createTime,
				"real_name":   name,
			}
			userBloggers = append(userBloggers, blogger)
			bloggerList = append(bloggerList, blogger)
		}
		bloggerRows.Close()

		byMember = append(byMember, map[string]interface{}{
			"name":     name,
			"count":    count,
			"bloggers": userBloggers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"byMember":    byMember,
			"bloggerList": bloggerList,
		},
	})
}

func (h *MiscHandler) ExportData(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	exportType := c.DefaultQuery("type", "csv")
	fieldsParam := c.DefaultQuery("fields", "")

	availableFields := map[string]string{
		"id":                "ID",
		"nickname":          "昵称",
		"category":          "分类",
		"product":           "对接产品",
		"contact":           "联系方式",
		"wechat":            "微信",
		"custom_contact":    "自定义联系方式",
		"platform":          "平台",
		"platform_account":  "平台账号",
		"description":       "简介",
		"tags":              "标签",
		"blogger_tags":      "博主标签",
		"user_name":         "用户名",
		"real_name":         "真实姓名",
		"status":            "状态",
		"followup_count":    "跟进次数",
		"cooperation_count": "合作次数",
		"create_time":       "创建时间",
		"update_time":       "更新时间",
	}

	var selectedFields []string
	if fieldsParam == "" {
		selectedFields = []string{"id", "nickname", "category", "product", "contact", "wechat", "platform", "status", "user_name", "create_time"}
	} else {
		for _, f := range strings.Split(fieldsParam, ",") {
			f = strings.TrimSpace(f)
			if _, ok := availableFields[f]; ok {
				selectedFields = append(selectedFields, f)
			}
		}
	}

	if len(selectedFields) == 0 {
		selectedFields = []string{"id", "nickname", "category", "product", "contact", "wechat", "platform", "status", "user_name", "create_time"}
	}

	headerNames := make([]string, len(selectedFields))
	for i, f := range selectedFields {
		headerNames[i] = availableFields[f]
	}

	bloggerRows, err := h.DB.Query(`
		SELECT b.id, b.nickname, b.category, b.product, b.contact, b.wechat, b.custom_contact,
			b.platform, b.platform_account, b.description, b.tags, b.avatar, b.user_name, b.real_name,
			b.status, b.create_time, b.update_time,
			(SELECT COUNT(*) FROM followup WHERE blogger_id = b.id AND is_deleted = 0) as followup_count,
			(SELECT COUNT(*) FROM cooperation WHERE blogger_id = b.id AND is_deleted = 0) as cooperation_count
		FROM blogger b WHERE b.is_deleted = 0 ORDER BY b.create_time DESC`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	defer bloggerRows.Close()

	var bloggers []map[string]interface{}
	bloggerMap := make(map[int]map[string]interface{})

	for bloggerRows.Next() {
		var id, followupCount, cooperationCount int
		var nickname, category, product, contact, wechat, customContact sql.NullString
		var platform, platformAccount, description, tags, avatar, userName, realName, status, createTime, updateTime sql.NullString

		bloggerRows.Scan(&id, &nickname, &category, &product, &contact, &wechat, &customContact,
			&platform, &platformAccount, &description, &tags, &avatar, &userName, &realName,
			&status, &createTime, &updateTime, &followupCount, &cooperationCount)

		blogMap := map[string]interface{}{
			"id":                id,
			"nickname":          nickname.String,
			"category":          category.String,
			"product":           product.String,
			"contact":           contact.String,
			"wechat":            wechat.String,
			"custom_contact":    customContact.String,
			"platform":          platform.String,
			"platform_account":  platformAccount.String,
			"description":       description.String,
			"tags":              tags.String,
			"avatar":            avatar.String,
			"user_name":         userName.String,
			"real_name":         realName.String,
			"status":            status.String,
			"followup_count":    followupCount,
			"cooperation_count": cooperationCount,
			"create_time":       createTime.String,
			"update_time":       updateTime.String,
		}

		tagRows, _ := h.DB.Query(`
			SELECT t.name FROM blogger_tags t
			INNER JOIN blogger_tag_relations r ON t.id = r.tag_id
			WHERE r.blogger_id = ?`, id)
		var tagNames []string
		for tagRows.Next() {
			var tagName string
			tagRows.Scan(&tagName)
			tagNames = append(tagNames, tagName)
		}
		tagRows.Close()
		blogMap["blogger_tags"] = strings.Join(tagNames, ", ")

		bloggers = append(bloggers, blogMap)
		bloggerMap[id] = blogMap
	}

	if exportType == "json" {
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=bloggers.json")
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": bloggers, "fields": selectedFields})
	} else {
		var csvContent bytes.Buffer
		csvContent.WriteString("\xEF\xBB\xBF")
		csvContent.WriteString(strings.Join(headerNames, ",") + "\n")

		for _, blogger := range bloggers {
			var row []string
			for _, field := range selectedFields {
				val := blogger[field]
				if val == nil {
					val = ""
				}
				row = append(row, escapeCSV(fmt.Sprintf("%v", val)))
			}
			csvContent.WriteString(strings.Join(row, ",") + "\n")
		}

		filename := "bloggers.csv"
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.String(http.StatusOK, csvContent.String())
	}
}

func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\\\n\r") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return "\"" + s + "\""
	}
	return s
}
