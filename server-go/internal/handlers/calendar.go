package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"server-go/internal/database"
)

type CalendarHandler struct{}

func NewCalendarHandler() *CalendarHandler {
	return &CalendarHandler{}
}

type CalendarEvent struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Start     string `json:"start"`
	End       string `json:"end,omitempty"`
	Color     string `json:"color"`
	BlogID    int    `json:"blogger_id,omitempty"`
	BlogName  string `json:"blogger_name,omitempty"`
	Status    string `json:"status,omitempty"`
	Remark    string `json:"remark,omitempty"`
}

func (h *CalendarHandler) GetCalendarEvents(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	startDate := c.Query("start")
	endDate := c.Query("end")

	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	}

	var events []CalendarEvent

	rows, err := database.DB.Query(`
		SELECT b.id, b.nickname, b.status, b.next_follow_time, b.contact, b.wechat
		FROM blogger b
		WHERE b.user_name = ? AND b.is_deleted = 0 AND b.is_invalid = 0
		AND b.next_follow_time IS NOT NULL AND b.next_follow_time != ''
		AND b.next_follow_time >= ? AND b.next_follow_time <= ?
		ORDER BY b.next_follow_time ASC
	`, userID, startDate, endDate)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var nickname, status, nextFollowTime, contact, wechat string
			rows.Scan(&id, &nickname, &status, &nextFollowTime, &contact, &wechat)
			events = append(events, CalendarEvent{
				ID:        id,
				Type:      "followup",
				Title:     "📅 跟进: " + nickname,
				Start:     nextFollowTime,
				Color:     getStatusColor(status),
				BlogID:    id,
				BlogName:  nickname,
				Status:    status,
				Remark:    buildFollowupRemark(contact, wechat),
			})
		}
	}

	rows2, err := database.DB.Query(`
		SELECT f.id, f.blogger_id, b.nickname, f.content, f.next_follow_time, f.create_time
		FROM followup f
		INNER JOIN blogger b ON f.blogger_id = b.id
		WHERE b.user_name = ? AND f.is_deleted = 0
		AND f.next_follow_time IS NOT NULL AND f.next_follow_time != ''
		AND f.next_follow_time >= ? AND f.next_follow_time <= ?
		ORDER BY f.next_follow_time ASC
	`, userID, startDate, endDate)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var id, bloggerID int
			var nickname, content, nextFollowTime, createTime string
			rows2.Scan(&id, &bloggerID, &nickname, &content, &nextFollowTime, &createTime)
			events = append(events, CalendarEvent{
				ID:        id,
				Type:      "followup_record",
				Title:     "📋 跟进记录: " + nickname,
				Start:     nextFollowTime,
				Color:     "#3b82f6",
				BlogID:    bloggerID,
				BlogName:  nickname,
				Remark:    truncateString(content, 50),
			})
		}
	}

	rows3, err := database.DB.Query(`
		SELECT c.id, c.blogger_id, b.nickname, c.product_name, c.cooperation_fee, c.status, c.start_time, c.end_time
		FROM cooperation c
		INNER JOIN blogger b ON c.blogger_id = b.id
		WHERE b.user_name = ? AND c.is_deleted = 0
		AND (
			(c.start_time IS NOT NULL AND c.start_time >= ? AND c.start_time <= ?)
			OR
			(c.end_time IS NOT NULL AND c.end_time >= ? AND c.end_time <= ?)
		)
		ORDER BY c.start_time ASC
	`, userID, startDate, endDate, startDate, endDate)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var id, bloggerID int
			var nickname, productName, status, startTime, endTime string
			var fee float64
			rows3.Scan(&id, &bloggerID, &nickname, &productName, &fee, &status, &startTime, &endTime)
			if startTime != "" && startTime >= startDate && startTime <= endDate {
				events = append(events, CalendarEvent{
					ID:        id,
					Type:      "cooperation_start",
					Title:     "🤝 合作开始: " + nickname + " - " + productName,
					Start:     startTime,
					Color:     "#10b981",
					BlogID:    bloggerID,
					BlogName:  nickname,
					Status:    status,
					Remark:    "费用: ¥" + formatFee(fee),
				})
			}
			if endTime != "" && endTime >= startDate && endTime <= endDate {
				events = append(events, CalendarEvent{
					ID:        id,
					Type:      "cooperation_end",
					Title:     "📝 合作结束: " + nickname + " - " + productName,
					Start:     endTime,
					Color:     "#f59e0b",
					BlogID:    bloggerID,
					BlogName:  nickname,
					Status:    status,
					Remark:    "费用: ¥" + formatFee(fee),
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": events})
}

func getStatusColor(status string) string {
	colorMap := map[string]string{
		"初次联系": "#6b7280",
		"洽谈中":   "#3b82f6",
		"已合作":   "#10b981",
		"已拒绝":   "#ef4444",
		"暂停合作": "#f59e0b",
	}
	if color, ok := colorMap[status]; ok {
		return color
	}
	return "#6b7280"
}

func buildFollowupRemark(contact, wechat string) string {
	remark := ""
	if contact != "" {
		remark += "📞 " + contact
	}
	if wechat != "" {
		if remark != "" {
			remark += " | "
		}
		remark += "💬 " + wechat
	}
	return remark
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

func formatFee(fee float64) string {
	if fee >= 10000 {
		return fmt.Sprintf("%.1f万", fee/10000)
	}
	return fmt.Sprintf("%.0f", fee)
}
