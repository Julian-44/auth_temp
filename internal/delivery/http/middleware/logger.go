package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		latency := time.Since(start)
		requestID, exists := c.Get("request_id")
		reqIDStr := ""
		if exists {
			reqIDStr = requestID.(string)
		}
		log.Printf("[%s] %s %s status=%d latency=%s",
			reqIDStr, c.Request.Method, path, c.Writer.Status(), latency)
	}
}
