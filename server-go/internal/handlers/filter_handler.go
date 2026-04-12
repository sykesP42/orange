package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server-go/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FilterHandler struct {
	DB *sql.DB
}

func NewFilterHandler(db *sql.DB) *FilterHandler {
	return &FilterHandler{DB: db}
}

func (h *FilterHandler) GetSavedFilters(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	rows, err := h.DB.Query(`SELECT id, user_id, name, filters, is_default, create_time
		FROM saved_filters WHERE user_id = ? ORDER BY is_default DESC, create_time DESC`, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.SavedFilter{}})
		return
	}
	defer rows.Close()

	var filters []models.SavedFilter
	for rows.Next() {
		var f models.SavedFilter
		rows.Scan(&f.ID, &f.UserID, &f.Name, &f.Filters, &f.IsDefault, &f.CreateTime)
		filters = append(filters, f)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": filters})
}

func (h *FilterHandler) CreateSavedFilter(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		Filters struct {
			Keyword  string `json:"keyword"`
			Category string `json:"category"`
			Platform string `json:"platform"`
			Status   string `json:"status"`
			TeamID   string `json:"team_id"`
			Tag      string `json:"tag"`
		} `json:"filters"`
		IsDefault bool `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "名称不能为空"})
		return
	}

	filtersJSON, _ := json.Marshal(req.Filters)
	now := time.Now().Format(time.RFC3339)

	if req.IsDefault {
		h.DB.Exec("UPDATE saved_filters SET is_default = 0 WHERE user_id = ?", userID)
	}

	result, err := h.DB.Exec(`INSERT INTO saved_filters (user_id, name, filters, is_default, create_time)
		VALUES (?, ?, ?, ?, ?)`,
		userID, req.Name, string(filtersJSON), boolToInt(req.IsDefault), now)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
		"data":    gin.H{"id": id},
	})
}

func (h *FilterHandler) UpdateSavedFilter(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	filterID, _ := strconv.Atoi(c.Param("id"))

	var existingUserID int
	h.DB.QueryRow("SELECT user_id FROM saved_filters WHERE id = ?", filterID).Scan(&existingUserID)
	if existingUserID != userID {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限修改此筛选器"})
		return
	}

	var req struct {
		Name    string `json:"name"`
		Filters struct {
			Keyword  string `json:"keyword"`
			Category string `json:"category"`
			Platform string `json:"platform"`
			Status   string `json:"status"`
			TeamID   string `json:"team_id"`
			Tag      string `json:"tag"`
		} `json:"filters"`
		IsDefault bool `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if req.IsDefault {
		h.DB.Exec("UPDATE saved_filters SET is_default = 0 WHERE user_id = ?", userID)
	}

	filtersJSON, _ := json.Marshal(req.Filters)
	now := time.Now().Format(time.RFC3339)

	updates := []string{}
	args := []interface{}{}

	if req.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, req.Name)
	}
	updates = append(updates, "filters = ?")
	args = append(args, string(filtersJSON))
	updates = append(updates, "is_default = ?")
	args = append(args, boolToInt(req.IsDefault))
	updates = append(updates, "create_time = ?")
	args = append(args, now)
	args = append(args, filterID)

	query := "UPDATE saved_filters SET " + joinStrings(updates, ", ") + " WHERE id = ?"
	h.DB.Exec(query, args...)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *FilterHandler) DeleteSavedFilter(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	filterID, _ := strconv.Atoi(c.Param("id"))

	var existingUserID int
	h.DB.QueryRow("SELECT user_id FROM saved_filters WHERE id = ?", filterID).Scan(&existingUserID)
	if existingUserID != userID {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限删除此筛选器"})
		return
	}

	h.DB.Exec("DELETE FROM saved_filters WHERE id = ?", filterID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *FilterHandler) SetDefaultFilter(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	filterID, _ := strconv.Atoi(c.Param("id"))

	var existingUserID int
	h.DB.QueryRow("SELECT user_id FROM saved_filters WHERE id = ?", filterID).Scan(&existingUserID)
	if existingUserID != userID {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权限修改此筛选器"})
		return
	}

	h.DB.Exec("UPDATE saved_filters SET is_default = 0 WHERE user_id = ?", userID)
	h.DB.Exec("UPDATE saved_filters SET is_default = 1 WHERE id = ?", filterID)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "设置成功"})
}