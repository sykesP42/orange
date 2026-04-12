package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	csrfCookieName = "XSRF-TOKEN"
	csrfHeaderName = "X-XSRF-TOKEN"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {
			setCSRFCookie(c)
			c.Next()
			return
		}

		cookie, err := c.Cookie(csrfCookieName)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "CSRF token 缺失",
			})
			c.Abort()
			return
		}

		headerToken := c.GetHeader(csrfHeaderName)
		if headerToken == "" {
			headerToken = c.PostForm("_csrf")
		}

		if headerToken != cookie {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "CSRF token 无效",
			})
			c.Abort()
			return
		}

		setCSRFCookie(c)
		c.Next()
	}
}

func setCSRFCookie(c *gin.Context) {
	token := uuid.New().String()
	c.SetCookie(
		csrfCookieName,
		token,
		3600,
		"/",
		"",
		false,
		false,
	)
	c.Set("csrf_token", token)
}
