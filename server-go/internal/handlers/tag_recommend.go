package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type TagRecommendationHandler struct {
	DB *sql.DB
}

func NewTagRecommendationHandler(db *sql.DB) *TagRecommendationHandler {
	return &TagRecommendationHandler{DB: db}
}

type TagCount struct {
	TagID    int    `json:"tag_id"`
	TagName  string `json:"tag_name"`
	TagColor string `json:"tag_color"`
	Count    int    `json:"count"`
	Score    float64 `json:"score"`
}

func (h *TagRecommendationHandler) GetRecommendations(c *gin.Context) {
	bloggerIDStr := c.Query("blogger_id")
	category := c.Query("category")
	platform := c.Query("platform")
	product := c.Query("product")
	description := c.Query("description")

	if bloggerIDStr != "" {
		var bCategory, bPlatform, bProduct, bDescription string
		h.DB.QueryRow("SELECT category, platform, product, description FROM blogger WHERE id = ? AND is_deleted = 0", bloggerIDStr).Scan(&bCategory, &bPlatform, &bProduct, &bDescription)
		if bCategory != "" {
			category = bCategory
		}
		if bPlatform != "" {
			platform = bPlatform
		}
		if bProduct != "" {
			product = bProduct
		}
		if bDescription != "" {
			description = bDescription
		}
	}

	tagCounts := make(map[int]*TagCount)

	conditions := []string{}
	args := []interface{}{}

	if category != "" {
		conditions = append(conditions, "b.category = ?")
		args = append(args, category)
	}
	if platform != "" {
		conditions = append(conditions, "b.platform = ?")
		args = append(args, platform)
	}

	if len(conditions) > 0 {
		query := `
			SELECT bt.tag_id, t.name, t.color, COUNT(*) as cnt
			FROM blogger_tag_relations bt
			JOIN blogger b ON bt.blogger_id = b.id
			JOIN blogger_tags t ON bt.tag_id = t.id
			WHERE b.is_deleted = 0 AND b.is_invalid = 0
		`

		for i, cond := range conditions {
			if i == 0 {
				query += " AND " + cond
			} else {
				query += " OR " + cond
			}
		}

		query += " GROUP BY bt.tag_id ORDER BY cnt DESC LIMIT 50"

		rows, err := h.DB.Query(query, args...)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var tc TagCount
				rows.Scan(&tc.TagID, &tc.TagName, &tc.TagColor, &tc.Count)
				tc.Score = float64(tc.Count)
				tagCounts[tc.TagID] = &tc
			}
		}
	}

	descriptionKeywords := extractKeywords(description)
	if len(descriptionKeywords) > 0 {
		for keyword, weight := range descriptionKeywords {
			query := `
				SELECT t.id, t.name, t.color, COUNT(*) as cnt
				FROM blogger_tags t
				JOIN blogger_tag_relations bt ON t.id = bt.tag_id
				JOIN blogger b ON bt.blogger_id = b.id
				WHERE b.is_deleted = 0 AND (
					LOWER(t.name) LIKE ? OR
					LOWER(t.name) LIKE ? OR
					LOWER(t.name) LIKE ?
				)
				GROUP BY t.id
				LIMIT 20
			`
			likeKeyword := "%" + keyword + "%"
			rows, err := h.DB.Query(query, likeKeyword, likeKeyword, likeKeyword)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var tc TagCount
					var name string
					rows.Scan(&tc.TagID, &name, &tc.TagColor, &tc.Count)
					tc.TagName = name
					tc.Score = float64(tc.Count) * weight * 2
					if existing, exists := tagCounts[tc.TagID]; exists {
						existing.Score += tc.Score
						existing.Count += tc.Count
					} else {
						tagCounts[tc.TagID] = &tc
					}
				}
			}
		}
	}

	var recommendations []TagCount
	for _, tc := range tagCounts {
		recommendations = append(recommendations, *tc)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	var existingTagIDs []int
	if bloggerIDStr != "" {
		rows, _ := h.DB.Query("SELECT tag_id FROM blogger_tag_relations WHERE blogger_id = ?", bloggerIDStr)
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var tid int
				rows.Scan(&tid)
				existingTagIDs = append(existingTagIDs, tid)
			}
		}
	}

	for i := range recommendations {
		recommendations[i].Score = float64(recommendations[i].Count)
		for _, existingID := range existingTagIDs {
			if recommendations[i].TagID == existingID {
				recommendations[i].Score = 0
				break
			}
		}
	}

	var filtered []TagCount
	for _, r := range recommendations {
		if r.Score > 0 {
			filtered = append(filtered, r)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"recommendations": filtered,
			"features": gin.H{
				"category":  category,
				"platform":  platform,
				"product":   product,
				"keywords":  descriptionKeywords,
			},
		},
	})
}

func extractKeywords(text string) map[string]float64 {
	keywords := make(map[string]float64)

	keywordList := []string{
		"美妆", "护肤", "彩妆", "化妆", "口红", "粉底", "眼影",
		"美食", "餐厅", "小吃", "甜品", "烘焙", "烹饪", "食谱",
		"穿搭", "衣服", "时尚", "搭配", "潮流", "服装",
		"旅游", "旅行", "酒店", "景点", "攻略", "打卡",
		"科技", "数码", "手机", "电脑", "测评", "评测",
		"母婴", "育儿", "宝宝", "儿童", "奶粉",
		"健身", "运动", "瑜伽", "减肥", "健康",
		"家居", "装修", "家具", "收纳", "生活",
		"宠物", "猫", "狗", "萌宠",
		"游戏", "电竞", "直播",
		"职场", "办公", "效率",
	}

	for _, kw := range keywordList {
		count := 0
		textLower := text
		for i := 0; i < len(textLower)-len(kw)+1; i++ {
			if textLower[i:i+len(kw)] == kw {
				count++
			}
		}
		if count > 0 {
			keywords[kw] = float64(count)
		}
	}

	return keywords
}

func (h *TagRecommendationHandler) GetSimilarBloggers(c *gin.Context) {
	bloggerIDStr := c.Query("blogger_id")
	limitStr := c.DefaultQuery("limit", "5")

	limit := 5
	for _, ch := range limitStr {
		if ch >= '0' && ch <= '9' {
			limit = limit*10 + int(ch-'0')
		}
	}

	var bCategory, bPlatform, bTags string
	var bID int
	if bloggerIDStr != "" {
		h.DB.QueryRow("SELECT id, category, platform, tags FROM blogger WHERE id = ? AND is_deleted = 0", bloggerIDStr).Scan(&bID, &bCategory, &bPlatform, &bTags)
	}

	if bID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []interface{}{}})
		return
	}

	conditions := []string{}
	args := []interface{}{}

	if bCategory != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, bCategory)
	}
	if bPlatform != "" {
		conditions = append(conditions, "platform = ?")
		args = append(args, bPlatform)
	}

	conditions = append(conditions, "is_deleted = 0", "is_invalid = 0", "id != ?")
	args = append(args, bID)

	query := "SELECT id, nickname, category, platform, tags, avatar FROM blogger WHERE " + joinStrings(conditions, " AND ")
	query += " ORDER BY create_time DESC LIMIT ?"
	args = append(args, limit)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []interface{}{}})
		return
	}
	defer rows.Close()

	var bloggers []map[string]interface{}
	for rows.Next() {
		var id int
		var nickname, category, platform, tags, avatar string
		rows.Scan(&id, &nickname, &category, &platform, &tags, &avatar)

		var parsedTags []string
		json.Unmarshal([]byte(tags), &parsedTags)

		var similarity float64 = 0
		if bCategory == category {
			similarity += 40
		}
		if bPlatform == platform {
			similarity += 30
		}

		var bParsedTags []string
		json.Unmarshal([]byte(bTags), &bParsedTags)
		for _, bt := range bParsedTags {
			for _, pt := range parsedTags {
				if bt == pt {
					similarity += 15
				}
			}
		}

		bloggers = append(bloggers, map[string]interface{}{
			"id":         id,
			"nickname":   nickname,
			"category":   category,
			"platform":   platform,
			"tags":       parsedTags,
			"avatar":      avatar,
			"similarity": similarity,
		})
	}

	sort.Slice(bloggers, func(i, j int) bool {
		return bloggers[i]["similarity"].(float64) > bloggers[j]["similarity"].(float64)
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"bloggers":     bloggers,
			"reference_id": bID,
			"reference": gin.H{
				"category": bCategory,
				"platform": bPlatform,
			},
		},
	})
}