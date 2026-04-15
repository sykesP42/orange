package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort    string
	DatabasePath  string
	JWTSecret     string
	UploadPath    string
	MaxUploadSize int64
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

func Load() *Config {
	loadEnvFile()

	cfg := &Config{
		ServerPort:    getEnv("SERVER_PORT", "3003"),
		DatabasePath:  getEnv("DATABASE_PATH", "./data.db"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		UploadPath:    getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10*1024*1024),
		RedisHost:     getEnv("REDIS_HOST", ""),
		RedisPort:     getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
	}
	if cfg.JWTSecret == "your-secret-key-change-in-production" {
		log.Printf("⚠️ 警告：正在使用默认 JWT 密钥，请在生产环境中设置 JWT_SECRET 环境变量！")
	}
	return cfg
}

func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	buf := make([]byte, 1024*1024)
	n, _ := file.Read(buf)
	content := string(buf[:n])

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if _, exists := os.LookupEnv(key); !exists {
				os.Setenv(key, value)
			}
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}
