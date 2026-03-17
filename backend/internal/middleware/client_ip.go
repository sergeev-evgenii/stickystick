package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

const ClientIPKey = "client_ip"

// ClientIPMiddleware извлекает реальный IP из X-Forwarded-For, X-Real-IP или RemoteAddr
// и сохраняет в контексте под ключ "client_ip".
func ClientIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.GetHeader("X-Forwarded-For")
		if ip != "" {
			// Берём первый IP из списка (клиент), остальные — прокси
			ip = strings.TrimSpace(strings.Split(ip, ",")[0])
		}
		if ip == "" {
			ip = c.GetHeader("X-Real-IP")
		}
		if ip == "" {
			addr := c.Request.RemoteAddr
			if host, _, err := net.SplitHostPort(addr); err == nil {
				ip = host
			} else {
				ip = addr
			}
		}
		c.Set(ClientIPKey, ip)
		c.Next()
	}
}

// GetClientIP возвращает IP из контекста или пустую строку.
func GetClientIP(c *gin.Context) string {
	v, _ := c.Get(ClientIPKey)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
