package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"server-go/internal/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	DB *sql.DB
}

func NewWorkflowHandler(db *sql.DB) *WorkflowHandler {
	return &WorkflowHandler{DB: db}
}

type WorkflowRule struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	TriggerType  string `json:"trigger_type"`
	TriggerValue string `json:"trigger_value"`
	ActionType   string `json:"action_type"`
	ActionConfig string `json:"action_config"`
	Enabled      int    `json:"enabled"`
	Priority     int    `json:"priority"`
	CreateTime   string `json:"create_time"`
}

func (h *WorkflowHandler) GetRules(c *gin.Context) {
	rows, err := h.DB.Query(`SELECT id, name, trigger_type, trigger_value, action_type, action_config, enabled, priority, create_time
		FROM workflow_rules ORDER BY priority DESC, id ASC`)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []WorkflowRule{}})
		return
	}
	defer rows.Close()

	var rules []WorkflowRule
	for rows.Next() {
		var r WorkflowRule
		rows.Scan(&r.ID, &r.Name, &r.TriggerType, &r.TriggerValue, &r.ActionType, &r.ActionConfig, &r.Enabled, &r.Priority, &r.CreateTime)
		rules = append(rules, r)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": rules})
}

func (h *WorkflowHandler) CreateRule(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, _ := realNameVal.(string)

	var req struct {
		Name         string `json:"name" binding:"required"`
		TriggerType  string `json:"trigger_type" binding:"required"`
		TriggerValue string `json:"trigger_value" binding:"required"`
		ActionType   string `json:"action_type" binding:"required"`
		ActionConfig string `json:"action_config"`
		Enabled      bool   `json:"enabled"`
		Priority     int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	enabled := 0
	if req.Enabled {
		enabled = 1
	}

	result, err := h.DB.Exec(`INSERT INTO workflow_rules (name, trigger_type, trigger_value, action_type, action_config, enabled, priority, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.Name, req.TriggerType, req.TriggerValue, req.ActionType, req.ActionConfig, enabled, req.Priority, now)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	database.AddLog("工作流", "创建规则", realName, "创建工作流规则【"+req.Name+"】")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    gin.H{"id": id},
	})
}

func (h *WorkflowHandler) UpdateRule(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, _ := realNameVal.(string)

	ruleID, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Name         string `json:"name"`
		TriggerType  string `json:"trigger_type"`
		TriggerValue string `json:"trigger_value"`
		ActionType   string `json:"action_type"`
		ActionConfig string `json:"action_config"`
		Enabled      *bool  `json:"enabled"`
		Priority     int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, req.Name)
	}
	if req.TriggerType != "" {
		updates = append(updates, "trigger_type = ?")
		args = append(args, req.TriggerType)
	}
	if req.TriggerValue != "" {
		updates = append(updates, "trigger_value = ?")
		args = append(args, req.TriggerValue)
	}
	if req.ActionType != "" {
		updates = append(updates, "action_type = ?")
		args = append(args, req.ActionType)
	}
	if req.ActionConfig != "" {
		updates = append(updates, "action_config = ?")
		args = append(args, req.ActionConfig)
	}
	if req.Enabled != nil {
		enabled := 0
		if *req.Enabled {
			enabled = 1
		}
		updates = append(updates, "enabled = ?")
		args = append(args, enabled)
	}
	if req.Priority > 0 {
		updates = append(updates, "priority = ?")
		args = append(args, req.Priority)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "没有要更新的字段"})
		return
	}

	args = append(args, ruleID)
	query := "UPDATE workflow_rules SET " + joinStrings(updates, ", ") + " WHERE id = ?"
	h.DB.Exec(query, args...)

	database.AddLog("工作流", "更新规则", realName, "更新工作流规则ID:"+strconv.Itoa(ruleID))

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *WorkflowHandler) DeleteRule(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, _ := realNameVal.(string)

	ruleID, _ := strconv.Atoi(c.Param("id"))

	h.DB.Exec("DELETE FROM workflow_rules WHERE id = ?", ruleID)
	database.AddLog("工作流", "删除规则", realName, "删除工作流规则ID:"+strconv.Itoa(ruleID))

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *WorkflowHandler) ToggleRule(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, _ := realNameVal.(string)

	ruleID, _ := strconv.Atoi(c.Param("id"))

	var enabled int
	h.DB.QueryRow("SELECT enabled FROM workflow_rules WHERE id = ?", ruleID).Scan(&enabled)

	newEnabled := 0
	if enabled == 0 {
		newEnabled = 1
	}

	h.DB.Exec("UPDATE workflow_rules SET enabled = ? WHERE id = ?", newEnabled, ruleID)

	action := "启用"
	if newEnabled == 0 {
		action = "禁用"
	}
	database.AddLog("工作流", action+"规则", realName, action+"工作流规则ID:"+strconv.Itoa(ruleID))

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": action + "成功"})
}

func (h *WorkflowHandler) TriggerWorkflow(c *gin.Context) {
	realNameVal, _ := c.Get("realName")
	realName, _ := realNameVal.(string)

	var req struct {
		TriggerType  string `json:"trigger_type" binding:"required"`
		TriggerValue string `json:"trigger_value" binding:"required"`
		TargetID     int    `json:"target_id"`
		TargetType   string `json:"target_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	rows, err := h.DB.Query(`SELECT id, name, action_type, action_config FROM workflow_rules
		WHERE trigger_type = ? AND trigger_value = ? AND enabled = 1`, req.TriggerType, req.TriggerValue)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"triggered": 0}})
		return
	}
	defer rows.Close()

	triggeredCount := 0
	for rows.Next() {
		var ruleID int
		var ruleName, actionType, actionConfig string
		rows.Scan(&ruleID, &ruleName, &actionType, &actionConfig)

		h.executeAction(actionType, actionConfig, req.TargetID, req.TargetType, realName)
		triggeredCount++

		database.AddLog("工作流", "触发", realName, "触发工作流规则【"+ruleName+"】")
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "触发成功",
		"data": gin.H{
			"triggered":     triggeredCount,
			"trigger_type":  req.TriggerType,
			"trigger_value": req.TriggerValue,
		},
	})
}

func (h *WorkflowHandler) executeAction(actionType, actionConfig string, targetID int, targetType, operator string) {
	now := time.Now().Format(time.RFC3339)

	switch actionType {
	case "notification":
		var config struct {
			UserID   int    `json:"user_id"`
			Title    string `json:"title"`
			Content  string `json:"content"`
			Priority string `json:"priority"`
		}
		json.Unmarshal([]byte(actionConfig), &config)

		if config.UserID > 0 && config.Title != "" {
			h.DB.Exec(`INSERT INTO notifications (user_id, type, title, content, priority, blogger_id, from_user, create_time)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				config.UserID, "workflow", config.Title, config.Content, config.Priority, targetID, operator, now)
		}

	case "email":
		var config struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}
		json.Unmarshal([]byte(actionConfig), &config)
		database.AddLog("工作流", "邮件", operator, "发送邮件到:"+config.To)

	case "webhook":
		var config struct {
			URL     string            `json:"url"`
			Method  string            `json:"method"`
			Headers map[string]string `json:"headers"`
			Body    map[string]interface{} `json:"body"`
		}
		json.Unmarshal([]byte(actionConfig), &config)
		database.AddLog("工作流", "Webhook", operator, "调用:"+config.URL)

	case "tag_add":
		var config struct {
			TagIDs []int `json:"tag_ids"`
		}
		json.Unmarshal([]byte(actionConfig), &config)

		for _, tagID := range config.TagIDs {
			var exists int
			h.DB.QueryRow("SELECT COUNT(*) FROM blogger_tags WHERE blogger_id = ? AND tag_id = ?", targetID, tagID).Scan(&exists)
			if exists == 0 {
				h.DB.Exec("INSERT INTO blogger_tags (blogger_id, tag_id, create_time) VALUES (?, ?, ?)", targetID, tagID, now)
			}
		}

	case "status_change":
		var config struct {
			NewStatus string `json:"new_status"`
		}
		json.Unmarshal([]byte(actionConfig), &config)

		if config.NewStatus != "" && targetID > 0 {
			var oldStatus string
			h.DB.QueryRow("SELECT status FROM blogger WHERE id = ?", targetID).Scan(&oldStatus)

			h.DB.Exec("UPDATE blogger SET status = ?, status_update_time = ?, update_time = ? WHERE id = ?",
				config.NewStatus, now, now, targetID)

			h.DB.Exec(`INSERT INTO blogger_status_history (blogger_id, old_status, new_status, operator, remark, create_time)
				VALUES (?, ?, ?, ?, ?, ?)`,
				targetID, oldStatus, config.NewStatus, operator, "工作流自动变更", now)
		}
	}
}