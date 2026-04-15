package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"server-go/internal/database"
)

type AnalyticsHandler struct{}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{}
}

type BloggerOverview struct {
	TotalBloggers      int                    `json:"total_bloggers"`
	ActiveBloggers     int                    `json:"active_bloggers"`
	CooperationCount   int                    `json:"cooperation_count"`
	TotalExposure      int64                  `json:"total_exposure"`
	AvgFollowerCount   float64                `json:"avg_follower_count"`
	TopCategories      []CategoryStat         `json:"top_categories"`
	TopPlatforms       []PlatformStat         `json:"top_platforms"`
	MonthlyTrend       []MonthlyStat          `json:"monthly_trend"`
	StatusDistribution []StatusStat           `json:"status_distribution"`
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
	Percent  float64 `json:"percent"`
}

type PlatformStat struct {
	Platform string `json:"platform"`
	Count    int    `json:"count"`
	Percent  float64 `json:"percent"`
}

type MonthlyStat struct {
	Month   string `json:"month"`
	Added   int    `json:"added"`
	Cooperated int `json:"cooperated"`
}

type StatusStat struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	Percent float64 `json:"percent"`
}

type BloggerDetailedAnalytics struct {
	BloggerID      int                `json:"blogger_id"`
	Overview       BloggerOverview   `json:"overview"`
	FollowupStats FollowupStats     `json:"followup_stats"`
	CoopStats      CoopStats         `json:"coop_stats"`
	EngagementRate float64          `json:"engagement_rate"`
}

type FollowupStats struct {
	TotalFollowups    int     `json:"total_followups"`
	LastFollowup      string  `json:"last_followup"`
	AvgDaysBetween    float64 `json:"avg_days_between"`
	ResponseRate      float64 `json:"response_rate"`
}

type CoopStats struct {
	TotalCooperations int     `json:"total_cooperations"`
	ActiveCoops       int     `json:"active_coops"`
	CompletedCoops    int     `json:"completed_coops"`
	TotalInvestment   float64 `json:"total_investment"`
	AvgCoopFee        float64 `json:"avg_coop_fee"`
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var overview BloggerOverview

	database.DB.QueryRow(`SELECT COUNT(*) FROM blogger WHERE user_name = ? AND is_deleted = 0`, userID).Scan(&overview.TotalBloggers)

	database.DB.QueryRow(`
		SELECT COUNT(DISTINCT b.id) FROM blogger b
		WHERE b.user_name = ? AND b.is_deleted = 0
		AND b.status = '已合作'`, userID).Scan(&overview.ActiveBloggers)

	database.DB.QueryRow(`
		SELECT COUNT(*) FROM cooperation c
		INNER JOIN blogger b ON c.blogger_id = b.id
		WHERE b.user_name = ? AND c.is_deleted = 0`, userID).Scan(&overview.CooperationCount)

	database.DB.QueryRow(`
		SELECT COALESCE(SUM(follower_count), 0) FROM blogger
		WHERE user_name = ? AND is_deleted = 0`, userID).Scan(&overview.TotalExposure)

	database.DB.QueryRow(`
		SELECT COALESCE(AVG(follower_count), 0) FROM blogger
		WHERE user_name = ? AND is_deleted = 0`, userID).Scan(&overview.AvgFollowerCount)

	categoryRows, _ := database.DB.Query(`
		SELECT COALESCE(category, '未分类') as cat, COUNT(*) as cnt
		FROM blogger WHERE user_name = ? AND is_deleted = 0
		GROUP BY cat ORDER BY cnt DESC LIMIT 5`, userID)
	var catTotal int
	for categoryRows.Next() {
		var cat string
		var cnt int
		categoryRows.Scan(&cat, &cnt)
		overview.TopCategories = append(overview.TopCategories, CategoryStat{
			Category: cat,
			Count:    cnt,
		})
		catTotal += cnt
	}
	for i := range overview.TopCategories {
		if catTotal > 0 {
			overview.TopCategories[i].Percent = float64(overview.TopCategories[i].Count) / float64(catTotal) * 100
		}
	}
	categoryRows.Close()

	platformRows, _ := database.DB.Query(`
		SELECT COALESCE(platform, '未知') as plt, COUNT(*) as cnt
		FROM blogger WHERE user_name = ? AND is_deleted = 0
		GROUP BY plt ORDER BY cnt DESC LIMIT 5`, userID)
	var pltTotal int
	for platformRows.Next() {
		var plt string
		var cnt int
		platformRows.Scan(&plt, &cnt)
		overview.TopPlatforms = append(overview.TopPlatforms, PlatformStat{
			Platform: plt,
			Count:    cnt,
		})
		pltTotal += cnt
	}
	for i := range overview.TopPlatforms {
		if pltTotal > 0 {
			overview.TopPlatforms[i].Percent = float64(overview.TopPlatforms[i].Count) / float64(pltTotal) * 100
		}
	}
	platformRows.Close()

	now := time.Now()
	for i := 5; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		monthStr := month.Format("2006-01")

		var added, cooperated int
		database.DB.QueryRow(`
			SELECT COUNT(*) FROM blogger
			WHERE user_name = ? AND is_deleted = 0
			AND strftime('%Y-%m', create_time) = ?`, userID, monthStr).Scan(&added)

		database.DB.QueryRow(`
			SELECT COUNT(*) FROM cooperation c
			INNER JOIN blogger b ON c.blogger_id = b.id
			WHERE b.user_name = ? AND c.is_deleted = 0
			AND strftime('%Y-%m', c.create_time) = ?`, userID, monthStr).Scan(&cooperated)

		overview.MonthlyTrend = append(overview.MonthlyTrend, MonthlyStat{
			Month:      monthStr,
			Added:      added,
			Cooperated: cooperated,
		})
	}

	statusRows, _ := database.DB.Query(`
		SELECT COALESCE(status, '未知') as st, COUNT(*) as cnt
		FROM blogger WHERE user_name = ? AND is_deleted = 0
		GROUP BY st ORDER BY cnt DESC`, userID)
	var stTotal int
	for statusRows.Next() {
		var st string
		var cnt int
		statusRows.Scan(&st, &cnt)
		overview.StatusDistribution = append(overview.StatusDistribution, StatusStat{
			Status: st,
			Count:  cnt,
		})
		stTotal += cnt
	}
	for i := range overview.StatusDistribution {
		if stTotal > 0 {
			overview.StatusDistribution[i].Percent = float64(overview.StatusDistribution[i].Count) / float64(stTotal) * 100
		}
	}
	statusRows.Close()

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": overview})
}

func (h *AnalyticsHandler) GetBloggerAnalytics(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var bloggerID int
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &bloggerID); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的博主ID"})
		return
	}

	var realName string
	err := database.DB.QueryRow(`SELECT real_name FROM users WHERE id = ?`, userID).Scan(&realName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权访问"})
		return
	}

	var ownerName string
	database.DB.QueryRow(`SELECT user_name FROM blogger WHERE id = ?`, bloggerID).Scan(&ownerName)
	if ownerName != realName {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "无权访问此博主"})
		return
	}

	var analytics BloggerDetailedAnalytics
	analytics.BloggerID = bloggerID

	var totalFollowups sql.NullInt64
	var lastFollowupTime sql.NullString
	var avgDaysBetween float64
	database.DB.QueryRow(`
		SELECT COUNT(*), MAX(create_time),
			COALESCE(AVG(julianday(create_time) - julianday(LAG(create_time) OVER (ORDER BY create_time)), 0)
		FROM followup WHERE blogger_id = ? AND is_deleted = 0`, bloggerID).Scan(&totalFollowups, &lastFollowupTime, &avgDaysBetween)

	analytics.FollowupStats = FollowupStats{
		AvgDaysBetween: avgDaysBetween,
	}
	if totalFollowups.Valid {
		analytics.FollowupStats.TotalFollowups = int(totalFollowups.Int64)
	}
	if lastFollowupTime.Valid {
		analytics.FollowupStats.LastFollowup = lastFollowupTime.String
	}

	var totalCoops, activeCoops, completedCoops int
	var totalInvestment, avgFee float64
	database.DB.QueryRow(`
		SELECT COUNT(*),
			SUM(CASE WHEN status = '进行中' THEN 1 ELSE 0 END),
			SUM(CASE WHEN status = '已完成' THEN 1 ELSE 0 END),
			COALESCE(SUM(cooperation_fee), 0),
			COALESCE(AVG(cooperation_fee), 0)
		FROM cooperation WHERE blogger_id = ? AND is_deleted = 0`, bloggerID).Scan(
		&totalCoops, &activeCoops, &completedCoops, &totalInvestment, &avgFee)

	analytics.CoopStats = CoopStats{
		TotalCooperations: totalCoops,
		ActiveCoops:       activeCoops,
		CompletedCoops:    completedCoops,
		TotalInvestment:   totalInvestment,
		AvgCoopFee:        avgFee,
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": analytics})
}
