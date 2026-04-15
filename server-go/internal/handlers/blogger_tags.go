package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"server-go/internal/database"
)

type TagHandler struct{}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

func (h *TagHandler) GetTags(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, color, is_system, create_time FROM blogger_tags ORDER BY is_system DESC, id ASC")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取标签失败"})
		return
	}
	defer rows.Close()

	var tags []map[string]interface{}
	for rows.Next() {
		var id int
		var name, color string
		var isSystem int
		var createTime string
		rows.Scan(&id, &name, &color, &isSystem, &createTime)
		tags = append(tags, map[string]interface{}{
			"id":         id,
			"name":       name,
			"color":      color,
			"is_system":  isSystem == 1,
			"create_time": createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tags})
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Color == "" {
		req.Color = "#6b7280"
	}

	result, err := database.DB.Exec("INSERT INTO blogger_tags (name, color, is_system) VALUES (?, ?, 0)", req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建标签失败，名称可能已存在"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": map[string]interface{}{
		"id":    id,
		"name":  req.Name,
		"color": req.Color,
	}, "message": "创建成功"})
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Color == "" {
		req.Color = "#6b7280"
	}

	result, err := database.DB.Exec("UPDATE blogger_tags SET name = ?, color = ? WHERE id = ? AND is_system = 0", req.Name, req.Color, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新标签失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "标签不存在或为系统标签无法修改"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	result, err := database.DB.Exec("DELETE FROM blogger_tags WHERE id = ? AND is_system = 0", id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除标签失败"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "标签不存在或为系统标签无法删除"})
		return
	}

	database.DB.Exec("DELETE FROM blogger_tag_relations WHERE tag_id = ?", id)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *TagHandler) GetBloggerTags(c *gin.Context) {
	bloggerIDStr := c.Param("id")
	bloggerID, _ := strconv.Atoi(bloggerIDStr)

	rows, err := database.DB.Query(`
		SELECT t.id, t.name, t.color, t.is_system
		FROM blogger_tags t
		INNER JOIN blogger_tag_relations r ON t.id = r.tag_id
		WHERE r.blogger_id = ?
		ORDER BY t.is_system DESC, t.id ASC
	`, bloggerID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取博主标签失败"})
		return
	}
	defer rows.Close()

	var tags []map[string]interface{}
	for rows.Next() {
		var id int
		var name, color string
		var isSystem bool
		rows.Scan(&id, &name, &color, &isSystem)
		tags = append(tags, map[string]interface{}{
			"id":        id,
			"name":      name,
			"color":     color,
			"is_system": isSystem,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tags})
}

func (h *TagHandler) SetBloggerTags(c *gin.Context) {
	bloggerIDStr := c.Param("id")
	bloggerID, _ := strconv.Atoi(bloggerIDStr)

	var req struct {
		TagIDs []int `json:"tag_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM blogger_tag_relations WHERE blogger_id = ?", bloggerID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	stmt, err := tx.Prepare("INSERT INTO blogger_tag_relations (blogger_id, tag_id) VALUES (?, ?)")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	defer stmt.Close()

	for _, tagID := range req.TagIDs {
		_, err = stmt.Exec(bloggerID, tagID)
		if err != nil {
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "设置成功"})
}

func (h *TagHandler) GetBloggersByTag(c *gin.Context) {
	tagIDStr := c.Query("tag_id")
	tagID, _ := strconv.Atoi(tagIDStr)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	offset := (page - 1) * pageSize

	var total int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM blogger b
		INNER JOIN blogger_tag_relations r ON b.id = r.blogger_id
		WHERE r.tag_id = ? AND b.is_deleted = 0
	`, tagID).Scan(&total)

	rows, err := database.DB.Query(`
		SELECT b.id, b.nickname, b.category, b.platform, b.status, b.avatar, b.user_name
		FROM blogger b
		INNER JOIN blogger_tag_relations r ON b.id = r.blogger_id
		WHERE r.tag_id = ? AND b.is_deleted = 0
		ORDER BY b.id DESC
		LIMIT ? OFFSET ?
	`, tagID, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": []interface{}{}, "total": 0}})
		return
	}
	defer rows.Close()

	var bloggers []map[string]interface{}
	for rows.Next() {
		var id int
		var nickname, category, platform, status, avatar, userName sql.NullString
		rows.Scan(&id, &nickname, &category, &platform, &status, &avatar, &userName)
		bloggers = append(bloggers, map[string]interface{}{
			"id":       id,
			"nickname": nickname.String,
			"category": category.String,
			"platform": platform.String,
			"status":   status.String,
			"avatar":   avatar.String,
			"user_name": userName.String,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{
		"list":    bloggers,
		"total":   total,
		"page":    page,
		"pageSize": pageSize,
	}})
}
