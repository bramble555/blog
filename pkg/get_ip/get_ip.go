package get_ip

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	if c == nil {
		return "127.0.0.1"
	}
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP",
	}
	for _, header := range headers {
		if value := c.GetHeader(header); value != "" {
			parts := strings.Split(value, ",")
			for _, part := range parts {
				if ip := normalizeIP(part); ip != "" {
					return ip
				}
			}
		}
	}
	if ip := normalizeIP(c.ClientIP()); ip != "" {
		return ip
	}
	return "127.0.0.1"
}

func normalizeIP(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return ""
	}
	if ip := net.ParseIP(v); ip != nil {
		return ip.String()
	}
	host, _, err := net.SplitHostPort(v)
	if err == nil {
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
	}
	return ""
}
