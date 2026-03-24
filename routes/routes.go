package routes

import (
	"time"

	"github.com/Astrasv/go-gully-backend/auth"
	"github.com/Astrasv/go-gully-backend/handlers"
	"github.com/Astrasv/go-gully-backend/middleware"
	ratelimiter "github.com/Astrasv/go-gully-backend/middleware/ratelimiter"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	authRoutes := r.Group("/auth")
	authRoutes.Use(RateLimitMiddleware())
	{
		authRoutes.GET("/google/login", auth.Login)
		authRoutes.GET("/google/callback", auth.Callback)
		authRoutes.POST("/logout", auth.Logout)

		// To test without auth
		// authRoutes.POST("/query-agent", handlers.QueryAgent) 
		// authRoutes.GET("/entities", handlers.GetEntities) //-> Comment a block in handlers/sql_agent.go
		

	}

	protected := r.Group("/api")
	protected.Use(middleware.RequireAuth)
	{
		protected.GET("/me", auth.GetCurrentUser)
		protected.GET("/validate", handlers.Validate)
		protected.GET("/protected", handlers.Protected)
		protected.POST("/query-agent", handlers.QueryAgent)
		protected.GET("/query-history", handlers.GetQueryHistory)
		// protected.GET("/entities", handlers.GetEntities)
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return ratelimiter.RequireRateLimiter(ratelimiter.RateLimiter{
		RateLimiterType: ratelimiter.IPRateLimiter,
		Key:             "iplimiter_maximum_requests_for_ip_test",
		Option: ratelimiter.RateLimiterOption{
			Limit: 1,
			Burst: 1,
			Len:   1 * time.Second,
		},
	})
}
