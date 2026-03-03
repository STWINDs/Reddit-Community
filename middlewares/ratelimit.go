package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware 令牌桶限流中间件
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	// 每 fillInterval 时间放入一个令牌，桶容量为 cap
	limiter := rate.NewLimiter(rate.Every(fillInterval), int(cap))
	return func(c *gin.Context) {
		// 如果拿不到令牌，说明请求过快，直接返回 429 状态码
		if !limiter.Allow() {
			c.String(http.StatusTooManyRequests, "Too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
