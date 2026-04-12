package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"server-go/internal/database"
)

type FollowupTemplateHandler struct{}

func NewFollowupTemplateHandler() *FollowupTemplateHandler {
	return &FollowupTemplateHandler{}
}

func (h *FollowupTemplateHandler) GetTemplates(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, name, content, category, is_default, use_count, create_time
		FROM followup_templates
		WHERE user_id IS NULL OR user_id = ?
		ORDER BY is_default DESC, use_count DESC, id ASC
	`, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取模板失败"})
		return
	}
	defer rows.Close()

	var templates []map[string]interface{}
	for rows.Next() {
		var id, useCount int
		var name, content, category string
		var isDefault bool
		var createTime string
		rows.Scan(&id, &name, &content, &category, &isDefault, &useCount, &createTime)
		templates = append(templates, map[string]interface{}{
			"id":          id,
			"name":       name,
			"content":    content,
			"category":   category,
			"is_default": isDefault,
			"use_count": useCount,
			"create_time": createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": templates})
}

func (h *FollowupTemplateHandler) CreateTemplate(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Category == "" {
		req.Category = "通用"
	}

	result, err := database.DB.Exec(`
		INSERT INTO followup_templates (user_id, name, content, category)
		VALUES (?, ?, ?, ?)
	`, userID, req.Name, req.Content, req.Category)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建模板失败"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": map[string]interface{}{
		"id":       id,
		"name":     req.Name,
		"content":  req.Content,
		"category": req.Category,
	}, "message": "创建成功"})
}

func (h *FollowupTemplateHandler) UpdateTemplate(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Name     string `json:"name" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Category == "" {
		req.Category = "通用"
	}

	result, err := database.DB.Exec(`
		UPDATE followup_templates
		SET name = ?, content = ?, category = ?, update_time = CURRENT_TIMESTAMP
		WHERE id = ? AND (user_id IS NULL OR user_id = ?)
	`, req.Name, req.Content, req.Category, id, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新模板失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在或无权修改"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *FollowupTemplateHandler) DeleteTemplate(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	result, err := database.DB.Exec(`
		DELETE FROM followup_templates
		WHERE id = ? AND user_id = ?
	`, id, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除模板失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在或无权删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *FollowupTemplateHandler) UseTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	database.DB.Exec("UPDATE followup_templates SET use_count = use_count + 1 WHERE id = ?", id)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "记录成功"})
}

func (h *FollowupTemplateHandler) SetDefault(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	database.DB.Exec("UPDATE followup_templates SET is_default = 0 WHERE user_id IS NULL OR user_id = ?", userID)
	database.DB.Exec("UPDATE followup_templates SET is_default = 1 WHERE id = ?", id)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "设置成功"})
}
