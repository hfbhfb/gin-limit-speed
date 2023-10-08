package main

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func main() {
	r := gin.New()

	// 添加限流中间件
	r.Use(rateLimitMiddleware(10, time.Second)) // 10 requests per second

	// 添加路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
	})

	// 启动服务
	r.Run(":8080")
}

func rateLimitMiddleware(requestsPerSecond int, duration time.Duration) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(duration, int64(requestsPerSecond), int64(requestsPerSecond))

	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}

