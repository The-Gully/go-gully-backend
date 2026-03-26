package routes

import (
	"time"

	"github.com/Astrasv/go-gully-backend/auth/google"
	"github.com/Astrasv/go-gully-backend/auth/local"
	"github.com/Astrasv/go-gully-backend/auth/verification"
	"github.com/Astrasv/go-gully-backend/handlers"
	"github.com/Astrasv/go-gully-backend/middleware"
	ratelimiter "github.com/Astrasv/go-gully-backend/middleware/ratelimiter"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	authRoutes := r.Group("/auth")
	authRoutes.Use(RateLimitMiddleware())
	{
		authRoutes.GET("/google/login", google.Login)
		authRoutes.GET("/google/callback", google.Callback)
		authRoutes.POST("/logout", google.Logout)

		authRoutes.POST("/register", local.Register)
		authRoutes.POST("/login", local.Login)

		authRoutes.GET("/verify-email", verification.VerifyEmailRedirect)
		authRoutes.POST("/verify-email", verification.VerifyEmailAPI)
		authRoutes.POST("/resend-verification", verification.ResendVerification)

	}

	protected := r.Group("/api")
	protected.Use(middleware.RequireAuth)
	protected.Use(middleware.RequireVerifiedEmail)
	{
		protected.GET("/me", google.GetCurrentUser)
		protected.GET("/validate", handlers.Validate)
		protected.GET("/protected", handlers.Protected)
		protected.POST("/query-agent", handlers.QueryAgent)
		protected.GET("/query-history", handlers.GetQueryHistory)
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
