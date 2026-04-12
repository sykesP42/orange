package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"server-go/internal/database"
)

type ImportHandler struct{}

func NewImportHandler() *ImportHandler {
	return &ImportHandler{}
}

type ImportPreview struct {
	TotalRows    int                   `json:"total_rows"`
	ValidRows    int                   `json:"valid_rows"`
	InvalidRows  int                   `json:"invalid_rows"`
	Errors       []ImportError         `json:"errors"`
	PreviewData  []map[string]string  `json:"preview_data"`
}

type ImportError struct {
	Row    int    `json:"row"`
	Field  string `json:"field"`
	Error  string `json:"error"`
}

type ImportResult struct {
	SuccessCount int      `json:"success_count"`
	FailCount   int      `json:"fail_count"`
	Errors      []ImportRowError `json:"errors"`
}

type ImportRowError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

func (h *ImportHandler) PreviewImport(c *gin.Context) {
	if _, exists := c.Get("userID"); !exists {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请上传文件"})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "读取文件失败"})
		return
	}

	filename := c.Request.FormValue("filename")
	var preview ImportPreview

	if strings.HasSuffix(strings.ToLower(filename), ".csv") {
		preview = h.parseCSV(content)
	} else if strings.HasSuffix(strings.ToLower(filename), ".json") {
		preview = h.parseJSON(content)
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只支持 CSV 或 JSON 格式"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": preview})
}

func (h *ImportHandler) DoImport(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		Filename string `json:"filename"`
		Data     []map[string]string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var realName string
	database.DB.QueryRow(`SELECT COALESCE(real_name, '') FROM users WHERE id = ?`, userID).Scan(&realName)

	var result ImportResult

	for i, row := range req.Data {
		err := h.validateAndInsertBlogger(row, realName, userID)
		if err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, ImportRowError{
				Row:   i + 1,
				Error: err.Error(),
			})
		} else {
			result.SuccessCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result, "message": fmt.Sprintf("成功导入 %d 条，失败 %d 条", result.SuccessCount, result.FailCount)})
}

func (h *ImportHandler) DownloadTemplate(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")

	if format == "json" {
		template := []map[string]string{
			{
				"nickname":        "示例博主1",
				"category":        "美妆",
				"platform":        "抖音",
				"platform_account": "douyin123",
				"contact":         "13800138000",
				"wechat":          "wechat123",
				"product":         "护肤套装",
				"description":     "优质美妆博主",
			},
			{
				"nickname":        "示例博主2",
				"category":        "美食",
				"platform":        "小红书",
				"platform_account": "xiaohongshu456",
				"contact":         "13900139000",
				"wechat":          "wechat456",
				"product":         "零食礼包",
				"description":     "美食探店博主",
			},
		}

		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=blogger_import_template.json")
		c.JSON(http.StatusOK, template)
	} else {
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=blogger_import_template.csv")

		csvContent := "\xEF\xBB\xBF"
		csvContent += "nickname,category,platform,platform_account,contact,wechat,product,description\n"
		csvContent += "示例博主1,美妆,抖音,douyin123,13800138000,wechat123,护肤套装,优质美妆博主\n"
		csvContent += "示例博主2,美食,小红书,xiaohongshu456,13900139000,wechat456,零食礼包,美食探店博主\n"

		c.String(http.StatusOK, csvContent)
	}
}

func (h *ImportHandler) parseCSV(content []byte) ImportPreview {
	reader := csv.NewReader(strings.NewReader(string(content)))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return ImportPreview{Errors: []ImportError{{Row: 0, Error: "CSV 解析失败"}}}
	}

	if len(records) < 2 {
		return ImportPreview{Errors: []ImportError{{Row: 0, Error: "文件内容为空或格式错误"}}}
	}

	headers := records[0]
	var preview ImportPreview
	preview.PreviewData = make([]map[string]string, 0, 10)

	for i := 1; i < len(records) && len(preview.PreviewData) < 10; i++ {
		row := records[i]
		if len(row) < len(headers) {
			preview.InvalidRows++
			preview.Errors = append(preview.Errors, ImportError{
				Row:   i + 1,
				Field: "row",
				Error: fmt.Sprintf("列数不匹配，期望 %d 列，实际 %d 列", len(headers), len(row)),
			})
			continue
		}

		record := make(map[string]string)
		valid := true
		for j, header := range headers {
			value := strings.TrimSpace(row[j])
			record[strings.ToLower(header)] = value

			if strings.ToLower(header) == "nickname" && value == "" {
				valid = false
				preview.Errors = append(preview.Errors, ImportError{
					Row:   i + 1,
					Field: "nickname",
					Error: "昵称不能为空",
				})
			}
		}

		if valid {
			preview.ValidRows++
			preview.PreviewData = append(preview.PreviewData, record)
		} else {
			preview.InvalidRows++
		}
		preview.TotalRows++
	}

	return preview
}

func (h *ImportHandler) parseJSON(content []byte) ImportPreview {
	var data []map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return ImportPreview{Errors: []ImportError{{Row: 0, Error: "JSON 解析失败"}}}
	}

	var preview ImportPreview
	preview.PreviewData = make([]map[string]string, 0, 10)

	for i, item := range data {
		if i >= 10 {
			break
		}

		record := make(map[string]string)
		valid := true

		for key, value := range item {
			strValue := fmt.Sprintf("%v", value)
			record[strings.ToLower(key)] = strValue

			if strings.ToLower(key) == "nickname" && strValue == "" {
				valid = false
				preview.Errors = append(preview.Errors, ImportError{
					Row:   i + 1,
					Field: "nickname",
					Error: "昵称不能为空",
				})
			}
		}

		if valid {
			preview.ValidRows++
			preview.PreviewData = append(preview.PreviewData, record)
		} else {
			preview.InvalidRows++
		}
		preview.TotalRows++
	}

	return preview
}

func (h *ImportHandler) validateAndInsertBlogger(row map[string]string, realName string, userID int) error {
	nickname := strings.TrimSpace(row["nickname"])
	if nickname == "" {
		return fmt.Errorf("昵称不能为空")
	}

	category := strings.TrimSpace(row["category"])
	platform := strings.TrimSpace(row["platform"])
	platformAccount := strings.TrimSpace(row["platform_account"])
	contact := strings.TrimSpace(row["contact"])
	wechat := strings.TrimSpace(row["wechat"])
	product := strings.TrimSpace(row["product"])
	description := strings.TrimSpace(row["description"])

	if category == "" {
		category = "未分类"
	}
	if platform == "" {
		platform = "未知"
	}

	result, err := database.DB.Exec(`
		INSERT INTO blogger (nickname, category, platform, platform_account, contact, wechat, product, description, user_name, real_name, status, create_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '初次联系', CURRENT_TIMESTAMP)
	`, nickname, category, platform, platformAccount, contact, wechat, product, description, realName, realName)

	if err != nil {
		return fmt.Errorf("插入失败: %v", err)
	}

	id, _ := result.LastInsertId()

	if contact != "" || wechat != "" {
		database.DB.Exec(`
			INSERT INTO followup (blogger_id, content, type, user_name, create_time)
			VALUES (?, '通过批量导入添加', '系统', ?, CURRENT_TIMESTAMP)
		`, id, realName)
	}

	return nil
}

func (h *ImportHandler) GetImportHistory(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var realName string
	database.DB.QueryRow(`SELECT COALESCE(real_name, '') FROM users WHERE id = ?`, userID).Scan(&realName)

	rows, err := database.DB.Query(`
		SELECT COUNT(*) as import_count, DATE(create_time) as import_date
		FROM blogger
		WHERE user_name = ? AND DATE(create_time) >= DATE('now', '-30 days')
		GROUP BY DATE(create_time)
		ORDER BY import_date DESC
		LIMIT 30
	`, realName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取导入历史失败"})
		return
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var count int
		var date string
		rows.Scan(&count, &date)
		history = append(history, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": history})
}

func (h *ImportHandler) ExportFailedRows(c *gin.Context) {
	var req struct {
		Errors []ImportRowError `json:"errors"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var csvContent strings.Builder
	csvContent.WriteString("\xEF\xBB\xBF")
	csvContent.WriteString("行号,错误信息\n")

	for _, err := range req.Errors {
		csvContent.WriteString(strconv.Itoa(err.Row))
		csvContent.WriteString(",")
		csvContent.WriteString(err.Error)
		csvContent.WriteString("\n")
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=import_errors.csv")
	c.String(http.StatusOK, csvContent.String())
}
