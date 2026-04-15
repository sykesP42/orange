package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"server-go/internal/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	Cfg *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{Cfg: cfg}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "No file uploaded"})
		return
	}

	if file.Size > h.Cfg.MaxUploadSize {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "File too large"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
	}

	if !allowedExts[ext] {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "File type not allowed"})
		return
	}

	uploadDir := h.Cfg.UploadPath
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to create upload directory"})
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "File uploaded successfully",
		"data": gin.H{
			"url":      "/uploads/" + filename,
			"filename": filename,
		},
	})
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "No image uploaded"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Image too large (max 5MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Only image files are allowed"})
		return
	}

	uploadDir := filepath.Join(h.Cfg.UploadPath, "images")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to create upload directory"})
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to open file"})
		return
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to create file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Image uploaded successfully",
		"data": gin.H{
			"url":      "/uploads/images/" + filename,
			"filename": filename,
		},
	})
}