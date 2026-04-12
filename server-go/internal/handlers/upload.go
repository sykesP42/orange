package handlers

import (
	"bytes"
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

var magicSignatures = map[string][]byte{
	".jpg":  {0xFF, 0xD8, 0xFF},
	".jpeg": {0xFF, 0xD8, 0xFF},
	".png":  {0x89, 0x50, 0x4E, 0x47},
	".gif":  {0x47, 0x49, 0x46},
	".webp": {0x52, 0x49, 0x46, 0x46},
	".pdf":  {0x25, 0x50, 0x44, 0x46},
}

var docSignatures = map[string][]byte{
	".doc":  {0xD0, 0xCF, 0x11, 0xE0},
	".docx": {0x50, 0x4B, 0x03, 0x04},
	".xls":  {0xD0, 0xCF, 0x11, 0xE0},
	".xlsx": {0x50, 0x4B, 0x03, 0x04},
}

func validateFileMagic(header []byte, ext string) bool {
	if sig, ok := magicSignatures[ext]; ok {
		return bytes.HasPrefix(header, sig)
	}
	if sig, ok := docSignatures[ext]; ok {
		return bytes.HasPrefix(header, sig)
	}
	return true
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

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to open file"})
		return
	}
	buf := make([]byte, 8)
	n, _ := src.Read(buf)
	src.Close()
	if n > 0 && !validateFileMagic(buf[:n], ext) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "File content does not match extension"})
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

	imgSrc, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "Failed to open file"})
		return
	}
	buf := make([]byte, 8)
	n, _ := imgSrc.Read(buf)
	imgSrc.Close()
	if n > 0 && !validateFileMagic(buf[:n], ext) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Image content does not match extension"})
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